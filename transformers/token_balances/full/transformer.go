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

package full

import (
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/filters"

	"github.com/vulcanize/account_transformers/transformers/token_balances/full/repositories"
	"github.com/vulcanize/account_transformers/transformers/token_balances/shared/constants"
)

type tokenBalanceTransformer struct {
	Converter                    Converter
	FilterRepository             repositories.Topic0FilterRepository
	FilteredLogsRepository       repositories.Topic0FilteredLogsRepository
	ValueTransferEventRepository repositories.ValueTransferEventRepository
	Filters                      []filters.LogFilter
}

func NewTokenBalanceTransformer(db *postgres.DB) *tokenBalanceTransformer {
	return &tokenBalanceTransformer{
		Converter:                    NewConverter(constants.ABIs),
		FilterRepository:             repositories.NewTopic0FilterRepository(db),
		FilteredLogsRepository:       repositories.NewTopic0FilteredLogsRepository(db),
		ValueTransferEventRepository: repositories.NewTokenBalanceRepository(db),
		Filters:                      constants.Filters,
	}
}

func (tbt *tokenBalanceTransformer) Init() error {
	return tbt.FilterRepository.CreateFilters(tbt.Filters)
}

func (tbt *tokenBalanceTransformer) Execute() error {
	for _, filter := range tbt.Filters {
		filteredLogs, err := tbt.FilteredLogsRepository.GetFilteredLogs(filter.Name, filter.FromBlock, filter.ToBlock)
		if err != nil {
			return err
		}
		models, err := tbt.Converter.Convert(filteredLogs)
		if err != nil {
			return err
		}
		err = tbt.ValueTransferEventRepository.CreateBalanceRecords(models)
		if err != nil {
			return err
		}
	}
	return nil
}
