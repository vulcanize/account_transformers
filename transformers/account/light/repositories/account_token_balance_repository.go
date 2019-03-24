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

package repositories

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/libraries/shared/utilities"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/account_transformers/transformers/account/shared"
)

type AccountTokenBalanceRepository interface {
	CreateTokenBalanceRecords(records []shared.TokenBalanceRecord, headerID int64) error
}

type accountTokenBalanceRepository struct {
	DB *postgres.DB
}

func NewAccountTokenBalanceRepository(db *postgres.DB) *accountTokenBalanceRepository {
	return &accountTokenBalanceRepository{
		DB: db,
	}
}

func (atbr *accountTokenBalanceRepository) CreateTokenBalanceRecords(records []shared.TokenBalanceRecord, headerID int64) error {
	tx, err := atbr.DB.Beginx()
	if err != nil {
		return err
	}
	pgStr := `INSERT INTO accounts.address_token_balances 
			(address_hash,
			block_number,
			token_contract_address_hash,
			value,
			value_fetched_at,
			inserted_at,
			updated_at) VALUES
			($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (address_hash, block_number, token_contract_address_hash)
			DO UPDATE SET
			(value,
			value_fetched_at,
			updated_at) = ($4, $5, $7)`
	for _, record := range records {
		now := time.Now()
		_, err := tx.Exec(pgStr, record.Address, record.BlockNumber, record.ContractAddress, utilities.NullToZero(record.Value), now, now, now)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = repository.MarkHeaderCheckedInTransaction(headerID, tx, common.BytesToAddress(record.Address).Hex())
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
