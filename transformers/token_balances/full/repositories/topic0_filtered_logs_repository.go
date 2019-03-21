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
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Topic0FilteredLogsRepository interface {
	GetFilteredLogs(name string, start, end int64) ([]*core.WatchedEvent, error)
}

type topic0FilteredLogsRepository struct {
	DB *postgres.DB
}

func NewTopic0FilteredLogsRepository(db *postgres.DB) *topic0FilteredLogsRepository {
	return &topic0FilteredLogsRepository{
		DB: db,
	}
}

func (flr *topic0FilteredLogsRepository) GetFilteredLogs(name string, start, end int64) ([]*core.WatchedEvent, error) {
	query := `SELECT id, name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM accounts.topic0_filtered_logs
			WHERE name = $1 AND block_number BETWEEN $2 AND $3`
	rows, err := flr.DB.Queryx(query, name, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lgs := make([]*core.WatchedEvent, 0)
	for rows.Next() {
		lg := new(core.WatchedEvent)
		err = rows.StructScan(lg)
		if err != nil {
			return nil, err
		}
		lgs = append(lgs, lg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lgs, nil
}
