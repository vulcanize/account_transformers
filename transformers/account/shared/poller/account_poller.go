// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package poller

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
)

type AccountPoller interface {
	PollAccounts(accounts []common.Address, blockNumber int64) ([]models.CoinBalanceRecord, error)
}

type accountPoller struct {
	blockChain core.BlockChain
}

func NewAccountPoller(bc core.BlockChain) *accountPoller {
	return &accountPoller{
		blockChain: bc,
	}
}

func (ap *accountPoller) PollAccounts(accounts []common.Address, blockNumber int64) ([]models.CoinBalanceRecord, error) {
	balanceRecords := make([]models.CoinBalanceRecord, 0)
	for _, addr := range accounts {
		record := models.CoinBalanceRecord{
			BlockNumber: blockNumber,
			Address:     addr.Bytes(),
		}
		balance, err := ap.blockChain.GetAccountBalance(addr, big.NewInt(blockNumber))
		if err != nil {
			return nil, err
		}
		record.Value = balance.String()
		balanceRecords = append(balanceRecords, record)
	}
	return balanceRecords, nil
}
