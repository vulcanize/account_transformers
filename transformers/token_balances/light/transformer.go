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

package light

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/account_transformers/transformers/token_balances/shared/constants"
)

type tokenBalanceTransformer struct {
	Converter Converter
	Fetcher   fetcher.Fetcher
	HeaderRepository repository.HeaderRepository
	ValueTransferEventRepository ValueTransferEventRepository
}

func NewTokenBalanceTransformer(db *postgres.DB) *tokenBalanceTransformer {
	return &tokenBalanceTransformer{
		Converter: NewConverter(constants.ABIs),
		HeaderRepository: repository.NewHeaderRepository(db),
		ValueTransferEventRepository: NewTokenBalanceRepository(db),
	}
}

func (tbt *tokenBalanceTransformer) Init() error {
	return nil
}

func (tbt *tokenBalanceTransformer) Execute() error {
	missingHeaders, err := tbt.HeaderRepository.MissingHeaders(0, -1, "token_value_transfer")
	if err != nil {
		return err
	}
	for _, header := range missingHeaders {
		allLogs, err := tbt.Fetcher.FetchLogs(nil, constants.Topic0s, header)
		if err != nil {
			return err
		}
		models, err := tbt.Converter.Convert(allLogs, header.Id)
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