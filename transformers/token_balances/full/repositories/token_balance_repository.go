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
	"github.com/vulcanize/account_transformers/transformers/token_balances/full"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type ValueTransferEventRepository interface {
	CreateBalanceRecords(models []full.ValueTransferModel) error
}

type valueTransferEventRepository struct {
	DB *postgres.DB
}

func NewTokenBalanceRepository(db *postgres.DB) *valueTransferEventRepository {
	return &valueTransferEventRepository{
		DB: db,
	}
}

func (br *valueTransferEventRepository) CreateBalanceRecords(models []full.ValueTransferModel) error {
	tx, err := br.DB.Beginx()
	if err != nil {
		return err
	}
	pgStr := `INSERT INTO accounts.value_transfer_events  
		(name,
		dst," 
		src," 
		amount,
		contract,
		log_id) VALUES
		($1, $2, $3, $4, $5, $6)
		ON CONFLICT (log_id) DO UPDATE SET 
		(name,
		dst," 
		src," 
		amount,
		contract) = ($1, $2, $3, $4, $5)`
	for _, model := range models {
		_, err := tx.Exec(pgStr, model.Name, model.Dst, model.Src, model.Amount, model.Contract, model.VulcanizeLogID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}
