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

package light

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/account_transformers/transformers/account/light/converters"
	"github.com/vulcanize/account_transformers/transformers/account/light/repositories"
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
	"github.com/vulcanize/account_transformers/transformers/account/shared/poller"
)

type TokenBalanceTransformer struct {
	Config                       config.ContractConfig
	ValueTransferConverter       converters.ValueTransferConverter
	Fetcher                      fetcher.Fetcher
	HeaderRepository             repository.HeaderRepository
	AddressRepository            repositories.AddressRepository
	WatchedContractRepository    repositories.WatchedContractRepository
	ValueTransferEventRepository repositories.ValueTransferEventRepository
	CoinBalanceRepository        repositories.AccountCoinBalanceRepository
	AccountPoller                poller.AccountPoller
	NextStart                    int64
}

func (tbt TokenBalanceTransformer) NewTransformer(db *postgres.DB, blockChain core.BlockChain) transformer.ContractTransformer {
	vtc, err := converters.NewValueTransferConverter(constants.CombinedABI, constants.EquivalentTokenAddressesMapping())
	if err != nil {
		logrus.Fatal(fmt.Sprintf("invalid abi\r\n%s\r\n", constants.CombinedABI), err)
	}
	return &TokenBalanceTransformer{
		ValueTransferConverter:       vtc,
		Fetcher:                      fetcher.NewFetcher(blockChain),
		HeaderRepository:             repository.NewHeaderRepository(db),
		AddressRepository:            repositories.NewAddressRepository(db),
		WatchedContractRepository:    repositories.NewWatchedContractRepository(db),
		ValueTransferEventRepository: repositories.NewValueTransferEventRepository(db),
		CoinBalanceRepository:        repositories.NewAccountCoinBalanceRepository(db),
		AccountPoller:                poller.NewAccountPoller(db, blockChain),
	}
}

func (tbt *TokenBalanceTransformer) Init() error {
	configuredAccountAddresses := constants.AccountAddresses()
	for _, addr := range configuredAccountAddresses {
		err := tbt.AddressRepository.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	configuredTokenAddresses := constants.TokenAddresses()
	for _, addr := range configuredTokenAddresses {
		err := tbt.WatchedContractRepository.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	tbt.NextStart = constants.StartingBlock()
	return tbt.HeaderRepository.AddCheckColumn("token_value_transfers")
}

func (tbt *TokenBalanceTransformer) Execute() error {
	missingHeaders, err := tbt.HeaderRepository.MissingHeaders(tbt.NextStart, -1, "token_value_transfers")
	if err != nil {
		return err
	}
	// First we need to transform all token value transfer type events into uniform value transfer records
	// Token balance records are a restricted view on this set of records
	for _, header := range missingHeaders {
		tbt.NextStart = header.BlockNumber
		allLogs, err := tbt.Fetcher.FetchLogs(nil, constants.Topic0s, header)
		if err != nil {
			return err
		}
		if len(allLogs) < 1 {
			err = tbt.HeaderRepository.MarkHeaderChecked(header.Id, "token_value_transfers")
			if err != nil {
				return err
			}
			tbt.NextStart++
			continue
		}
		models, err := tbt.ValueTransferConverter.Convert(allLogs, header.Id)
		if err != nil {
			return err
		}
		err = tbt.ValueTransferEventRepository.CreateTokenValueTransferRecords(models)
		if err != nil {
			return err
		}
		err = tbt.HeaderRepository.MarkHeaderChecked(header.Id, "token_value_transfers")
		if err != nil {
			return err
		}
		tbt.NextStart++
	}

	// Get the addresses we want to create eth and token balance records for
	addresses, err := tbt.AddressRepository.GetAddresses()
	if err != nil {
		return err
	}
	if len(addresses) < 1 {
		return errors.New("no addresses to create records for")
	}

	// Now we need to go through and collect eth balances and persist them into db records
	for _, addr := range addresses {
		columnID := "account_" + addr.Hex()
		err = tbt.HeaderRepository.AddCheckColumn(columnID)
		if err != nil {
			return err
		}
		// Retrieve headers which need records for this user
		checkedButNotRecordedHeaders, err := tbt.HeaderRepository.MissingMethodsCheckedEventsIntersection(0, -1, []string{columnID}, []string{"token_value_transfers"})
		if err != nil {
			return err
		}
		if len(checkedButNotRecordedHeaders) < 1 {
			continue
		}
		//mostRecentRecordsBlock := checkedButNotRecordedHeaders[0].BlockNumber - 1
		for _, header := range checkedButNotRecordedHeaders {
			// Let's also poll the eth balance at this header's blockNumber and create eth balance records
			coinBalanceRecords, err := tbt.AccountPoller.PollAccounts(addresses, header.BlockNumber, header.Id)
			if err != nil {
				return err
			}
			// And commit these records to Postgres
			err = tbt.CoinBalanceRepository.CreateCoinBalanceRecord(coinBalanceRecords, header.Id)
			if err != nil {
				return err
			}
			// Mark this header checked for this account
			err = tbt.HeaderRepository.MarkHeaderChecked(header.Id, columnID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (tbt *TokenBalanceTransformer) GetConfig() config.ContractConfig {
	return tbt.Config
}
