// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package converters

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/libraries/shared/utilities"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
	"github.com/vulcanize/account_transformers/transformers/account/shared"
)

type TokenBalanceConverter interface {
	Convert(mappedTransferRecords map[common.Address][]models.ValueTransferModel, blockNumber int64) []shared.TokenBalanceRecord
}

type tokenBalanceConverter struct{}

func NewTokenBalanceConverter() *tokenBalanceConverter {
	return &tokenBalanceConverter{}
}

func (c *tokenBalanceConverter) Convert(mappedTransferRecords map[common.Address][]models.ValueTransferModel, blockNumber int64) []shared.TokenBalanceRecord {
	balanceRecords := make([]shared.TokenBalanceRecord, 0)
	contractSortedTransferRecords := sortByContract(mappedTransferRecords)
	for addr, mappedRecords := range contractSortedTransferRecords {
		for contract, records := range mappedRecords {
			tokenBalanceRecord := shared.TokenBalanceRecord{
				Address:         addr.Bytes(),
				ContractAddress: contract.Bytes(),
				BlockNumber:     blockNumber,
			}
			value := big.NewInt(0)
			for _, record := range records {
				amount := new(big.Int)
				if record.Dst == addr.Hex() {
					amount.SetString(utilities.NullToZero(record.Amount), 10)
					value = value.Add(value, amount)
				}
				if record.Src == addr.Hex() {
					amount.SetString(utilities.NullToZero(record.Amount), 10)
					value = value.Sub(value, amount)
				}
			}
			tokenBalanceRecord.Value = value.String()
			balanceRecords = append(balanceRecords, tokenBalanceRecord)
		}
	}
	return balanceRecords
}

func sortByContract(mapping map[common.Address][]models.ValueTransferModel) map[common.Address]map[common.Address][]models.ValueTransferModel {
	returnRecords := make(map[common.Address]map[common.Address][]models.ValueTransferModel)
	for addr, records := range mapping {
		returnRecords[addr] = make(map[common.Address][]models.ValueTransferModel)
		for _, record := range records {
			returnRecords[addr][common.HexToAddress(record.Contract)] = append(returnRecords[addr][common.HexToAddress(record.Contract)], record)
		}
	}
	return returnRecords
}
