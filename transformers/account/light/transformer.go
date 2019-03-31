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
	"fmt"
	log "github.com/sirupsen/logrus"

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

type AccountTransformer struct {
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
	initialized                  bool
	QuitChannel                  chan bool
}

func (tbt AccountTransformer) NewTransformer(db *postgres.DB, blockChain core.BlockChain) transformer.ContractTransformer {
	return &AccountTransformer{
		Fetcher:                      fetcher.NewFetcher(blockChain),
		HeaderRepository:             repository.NewHeaderRepository(db),
		AddressRepository:            repositories.NewAddressRepository(db),
		WatchedContractRepository:    repositories.NewWatchedContractRepository(db),
		ValueTransferEventRepository: repositories.NewValueTransferEventRepository(db),
		CoinBalanceRepository:        repositories.NewAccountCoinBalanceRepository(db),
		AccountPoller:                poller.NewAccountPoller(db, blockChain),
		QuitChannel:                  make(chan bool),
	}
}

func (tbt *AccountTransformer) Init() error {
	var err error
	// Stop createBalanceRecords goroutine if it is already running, so that we can restart it with new init
	if tbt.initialized {
		tbt.QuitChannel <- true
	}
	// Get the list of account addresses we want to create records for from the config and add them to Postgres
	configuredAccountAddresses := constants.AccountAddresses()
	for _, addr := range configuredAccountAddresses {
		err = tbt.AddressRepository.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	// Get the list of token addresses we want to watch for from the config and add them to Postgres
	configuredTokenAddresses := constants.TokenAddresses()
	for _, addr := range configuredTokenAddresses {
		err = tbt.WatchedContractRepository.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	// Get the mapping of equivalent contracts from the config and use this mapping to initialize a value transfer event converter
	tbt.ValueTransferConverter, err = converters.NewValueTransferConverter(constants.CombinedABI, constants.EquivalentTokenAddressesMapping())
	if err != nil {
		return fmt.Errorf("invalid abi\r\n%s\r\n%v", constants.CombinedABI, err)
	}
	// Get the starting block from the config
	tbt.NextStart = constants.StartingBlock()
	// Add a column ID to the checked_headers table to track value transfer event log processing progress
	return tbt.HeaderRepository.AddCheckColumn("token_value_transfers")
}

func (tbt *AccountTransformer) Execute() error {
	missingHeaders, err := tbt.HeaderRepository.MissingHeaders(tbt.NextStart, -1, "token_value_transfers")
	if err != nil {
		return err
	}
	// First we need to transform all token value transfer type events into uniform value transfer records
	// Token balance records are a restricted view on this set of records
	for _, header := range missingHeaders {
		tbt.NextStart = header.BlockNumber // Set to the current header so that we restart on it if something goes wrong
		allLogs, err := tbt.Fetcher.FetchLogs(nil, constants.Topic0s, header)
		if err != nil {
			return err
		}
		if len(allLogs) < 1 {
			// No logs to process at this header, mark it checked and continue to the next
			err = tbt.HeaderRepository.MarkHeaderChecked(header.Id, "token_value_transfers")
			if err != nil {
				return err
			}
			tbt.NextStart++ // Set next header we need to start processing at
			continue
		}
		// Convert all of the fetched logs into a generic/uniform value transfer event model
		models, err := tbt.ValueTransferConverter.Convert(allLogs, header.Id)
		if err != nil {
			return err
		}
		// Write these models to Postgres in our token_value_transfer table
		// Views across these models for a particular user and contract create out token balance records
		err = tbt.ValueTransferEventRepository.CreateTokenValueTransferRecords(models)
		if err != nil {
			return err
		}
		// Mark this header checked for value transfer events
		err = tbt.HeaderRepository.MarkHeaderChecked(header.Id, "token_value_transfers")
		if err != nil {
			return err
		}
		tbt.NextStart++ // Set next header we need to start processing at
	}

	// If we haven't spun up a goroutine to create eth balance records (or if it was brought down by an error), do so now
	if !tbt.initialized {
		go tbt.createBalanceRecords()
		tbt.initialized = true
	}

	return nil
}

func (tbt *AccountTransformer) GetConfig() config.ContractConfig {
	return tbt.Config
}

func (tbt *AccountTransformer) createBalanceRecords() {
	for {
		select {
		case <-tbt.QuitChannel:
			tbt.initialized = false
			return
		default:
			// Get the addresses we want to create eth balance records for
			// We check each time in case a new one has been added to the Postgres table from an external source
			addresses, err := tbt.AddressRepository.GetAddresses()
			if err != nil {
				tbt.throwErr(err)
				return
			}
			// Cycle through the account addresses
			for _, addr := range addresses {
				// And create a checked_header id for them (IF NOT EXISTS)
				columnID := "account_" + addr.Hex()
				err = tbt.HeaderRepository.AddCheckColumn(columnID)
				if err != nil {
					tbt.throwErr(err)
					return
				}
				// Retrieve the headers which still need eth balance records for this user
				checkedButNotRecordedHeaders, err := tbt.HeaderRepository.MissingMethodsCheckedEventsIntersection(0, -1, []string{columnID}, []string{"token_value_transfers"})
				if err != nil {
					tbt.throwErr(err)
					return
				}
				// Create coin balance records for this account at each header
				coinBalanceRecords, err := tbt.AccountPoller.PollAccount(addr, checkedButNotRecordedHeaders)
				if err != nil {
					tbt.throwErr(err)
					return
				}
				// And commit these records to Postgres
				err = tbt.CoinBalanceRepository.CreateCoinBalanceRecords(coinBalanceRecords)
				if err != nil {
					tbt.throwErr(err)
					return
				}
				// Mark these headers checked for this account
				err = tbt.HeaderRepository.MarkHeadersCheckedForAll(checkedButNotRecordedHeaders, []string{columnID})
				if err != nil {
					tbt.throwErr(err)
					return
				}
			}
		}
	}
}

func (tbt *AccountTransformer) throwErr(err error) {
	log.Error("createBalanceRecords: error", err)
	tbt.initialized = false
}
