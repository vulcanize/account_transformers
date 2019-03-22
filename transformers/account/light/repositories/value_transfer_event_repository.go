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

package repositories

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/account_transformers/transformers/account/shared"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
)

type ValueTransferEventRepository interface {
	CreateTokenValueTransferRecords(models []models.ValueTransferModel) error
	GetTokenValueTransferRecordsForAccounts(addresses []common.Address, lastBlock int64) (map[common.Address][]models.ValueTransferModel, error)
}

type valueTransferEventRepository struct {
	DB *postgres.DB
}

func NewValueTransferEventRepository(db *postgres.DB) *valueTransferEventRepository {
	return &valueTransferEventRepository{
		DB: db,
	}
}

func (br *valueTransferEventRepository) CreateTokenValueTransferRecords(records []models.ValueTransferModel) error {
	tx, err := br.DB.Beginx()
	if err != nil {
		return err
	}
	pgStr := `INSERT INTO accounts.token_value_transfers  
		(header_id,
		name,
		dst," 
		src," 
		amount,
		contract,
		log_idx
		tx_idx,
		raw_log) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET 
		(name,
		dst," 
		src," 
		amount,
		contract,
		raw_log) = ($2, $3, $4, $5, $6, $9)`
	for _, model := range records {
		_, err := tx.Exec(pgStr, model.HeaderID, model.Name, model.Dst, model.Src, shared.NullToZero(model.Amount), model.Contract, model.LogIndex, model.TransactionIndex, model.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = repository.MarkHeaderCheckedInTransaction(model.HeaderID, tx, "token_value_transfers")
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (br *valueTransferEventRepository) GetTokenValueTransferRecordsForAccounts(addresses []common.Address, lastBlock int64) (map[common.Address][]models.ValueTransferModel, error) {
	mappedRecords := make(map[common.Address][]models.ValueTransferModel)
	pgStr := `SELECT header_id, name, block_number, dst, src, amount, contract, log_idx, tx_idx, raw_log FROM accounts.token_value_transfers
			WHERE (dst = $1 OR src = $1) AND block_number <= $2`
	for _, addr := range addresses {
		mappedRecords[addr] = make([]models.ValueTransferModel, 0)
		rows, err := br.DB.Queryx(pgStr, addr, lastBlock)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			record := new(models.ValueTransferModel)
			err = rows.StructScan(record)
			if err != nil {
				rows.Close()
				return nil, err
			}
			mappedRecords[addr] = append(mappedRecords[addr], *record)
		}
		if err = rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
	}
	return mappedRecords, nil
}
