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
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

type Topic0FilterRepository interface {
	CreateFilters(filters []filters.LogFilter) error
}

type topic0FilterRepository struct {
	DB 	*postgres.DB
}

func NewTopic0FilterRepository(db *postgres.DB) *topic0FilterRepository {
	return &topic0FilterRepository{
		DB: db,
	}
}

func (fr *topic0FilterRepository) CreateFilters(filters []filters.LogFilter) error {
	tx, err := fr.DB.Beginx()
	if err != nil {
		return err
	}
	pgStr := `INSERT INTO accounts.topic0_filters
				(name,
				from_block,
				to_block,
				topic0,
				topic1,
				topic2,
				topic3)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (name) DO UPDATE SET
				(from_block,
				to_block,
				topic0,
				topic1,
				topic2,
				topic3) = ($2, $3, $4, $5, $6, $7)`
	for _, filter := range filters {
		_, err = tx.Exec(pgStr, filter.Name, filter.FromBlock, filter.ToBlock, filter.Topics[0], filter.Topics[1], filter.Topics[2], filter.Topics[3])
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
