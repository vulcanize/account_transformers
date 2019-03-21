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
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Topic0FilteredLogsRepository interface {
	GetFilteredLogs(name string, start, end int64) ([]*core.WatchedEvent, int64, error)
}

type topic0FilteredLogsRepository struct {
	DB *postgres.DB
}

func NewTopic0FilteredLogsRepository(db *postgres.DB) *topic0FilteredLogsRepository {
	return &topic0FilteredLogsRepository{
		DB: db,
	}
}

func (flr *topic0FilteredLogsRepository) GetFilteredLogs(name string, start, end int64) ([]*core.WatchedEvent, int64, error) {
	var query string
	var rows *sqlx.Rows
	var err error
	if end < 0 {
		query = `SELECT id, name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM accounts.topic0_filtered_logs
			WHERE name = $1 AND block_number >= $2
			ORDER BY accounts.topic0_filtered_logs.block_number LIMIT 500`
		rows, err = flr.DB.Queryx(query, name, start)
		if err != nil {
			return nil, start, err
		}
	} else {
		query = `SELECT id, name, block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data FROM accounts.topic0_filtered_logs
			WHERE name = $1 AND block_number BETWEEN $2 AND $3
			ORDER BY accounts.topic0_filtered_logs.block_number LIMIT 500`
		rows, err = flr.DB.Queryx(query, name, start, end)
		if err != nil {
			return nil, start, err
		}
	}
	defer rows.Close()

	logs := make([]*core.WatchedEvent, 0)
	for rows.Next() {
		log := new(core.WatchedEvent)
		err = rows.StructScan(log)
		if err != nil {
			return nil, start, err
		}
		logs = append(logs, log)
	}
	if err = rows.Err(); err != nil {
		return nil, start, err
	}
	contiguousLogs, nextStart := contiguousFilteredLogs(logs, start)
	return contiguousLogs, nextStart, nil
}

func contiguousFilteredLogs(logs []*core.WatchedEvent, start int64) ([]*core.WatchedEvent, int64) {
	if len(logs) < 1 {
		return []*core.WatchedEvent{}, start
	}
	if logs[0].BlockNumber != start {
		return []*core.WatchedEvent{}, start
	}
	contiguousLogs := make([]*core.WatchedEvent, 0, len(logs))
	contiguousLogs = append(contiguousLogs, logs[0])
	index := start + 1
	for _, log := range logs {
		if log.BlockNumber != index {
			return contiguousLogs, index
		}
		contiguousLogs = append(contiguousLogs, log)
		index++
	}
	return contiguousLogs, index
}