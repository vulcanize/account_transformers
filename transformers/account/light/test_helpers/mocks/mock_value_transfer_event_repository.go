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

package mocks

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
)

type MockValueTransferEventRepository struct {
	PassedRecords                              []models.ValueTransferModel
	CreateTokenValueTransferRecordsErr         error
	PassedAddresses                            []common.Address
	PassedLastBlock                            int64
	TokenValueTransferRecords                  map[common.Address][]models.ValueTransferModel
	GetTokenValueTransferRecordsForAccountsErr error
}

func (br *MockValueTransferEventRepository) CreateTokenValueTransferRecords(records []models.ValueTransferModel) error {
	br.PassedRecords = records
	return br.CreateTokenValueTransferRecordsErr
}

func (br *MockValueTransferEventRepository) GetTokenValueTransferRecordsForAccounts(addresses []common.Address, lastBlock int64) (map[common.Address][]models.ValueTransferModel, error) {
	br.PassedAddresses = addresses
	br.PassedLastBlock = lastBlock
	return br.TokenValueTransferRecords, br.GetTokenValueTransferRecordsForAccountsErr
}
