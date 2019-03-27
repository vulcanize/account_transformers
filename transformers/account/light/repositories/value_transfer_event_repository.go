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
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
	"github.com/vulcanize/account_transformers/transformers/account/shared"
)

type ValueTransferEventRepository interface {
	CreateTokenValueTransferRecords(models []models.ValueTransferModel) error
	GetTokenValueTransferRecordsForAccount(address common.Address, firstBlock, lastBlock int64) ([]models.ValueTransferModel, error)
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
		block_number,
		name,
		dst, 
		src,
		amount,
		contract,
		log_idx,
		tx_idx,
		raw_log) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET 
		(block_number,
		name,
		dst, 
		src, 
		amount,
		contract,
		raw_log) = ($2, $3, $4, $5, $6, $7, $10)`
	for _, record := range records {
		_, err := tx.Exec(pgStr, record.HeaderID, record.BlockNumber, record.Name, record.Dst, record.Src, shared.NullToZero(record.Amount), record.Contract, record.LogIndex, record.TransactionIndex, record.Raw)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (br *valueTransferEventRepository) GetTokenValueTransferRecordsForAccount(address common.Address, firstBlock, lastBlock int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	pgStr := `SELECT header_id, name, block_number, dst, src, amount, contract, log_idx, tx_idx, raw_log FROM accounts.token_value_transfers
			WHERE (dst = $1 OR src = $1) AND block_number BETWEEN $2 AND $3`
	rows, err := br.DB.Queryx(pgStr, address.Hex(), firstBlock, lastBlock)
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
		records = append(records, *record)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()
	return records, nil
}

func (br *valueTransferEventRepository) GetTokenValueTransferRecordsForAccounts(addresses []common.Address, lastBlock int64) (map[common.Address][]models.ValueTransferModel, error) {
	mappedRecords := make(map[common.Address][]models.ValueTransferModel)
	pgStr := `SELECT header_id, name, block_number, dst, src, amount, contract, log_idx, tx_idx, raw_log FROM accounts.token_value_transfers
			WHERE (dst = $1 OR src = $1) AND block_number <= $2`
	for _, addr := range addresses {
		rows, err := br.DB.Queryx(pgStr, addr.Hex(), lastBlock)
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
