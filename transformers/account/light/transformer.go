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
	"github.com/ethereum/go-ethereum/common"
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
	// Get the addresses we want to create eth balance records for
	// We check each execution cycle in case a new one has been added to the Postgres table from an external source
	addresses, err := tbt.AddressRepository.GetAddresses()
	if err != nil {
		return err
	}
	// Spin up a goroutine to process eth balance records in the background
	go tbt.processEthBalanceRecords(addresses)
	// Be sure to bring this goroutine down at the end, this blocks until the
	// Goroutine has finished its current processing cycle (select case), receives the quit signal,
	// Finishes catching up with processing headers
	defer func() {
		tbt.QuitChannel <- true
		<- tbt.QuitChannel
	}()
	// Get all the headers which we need to process token value transfer events for
	missingHeaders, err := tbt.HeaderRepository.MissingHeaders(tbt.NextStart, -1, "token_value_transfers")
	if err != nil {
		return err
	}
	// Transform all token value transfer type events into uniform value transfer records
	// User's token balance records are a view on this set of records
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
	return nil
}

func (tbt *AccountTransformer) GetConfig() config.ContractConfig {
	return tbt.Config
}

func (tbt *AccountTransformer) processEthBalanceRecords(addresses []common.Address) {
	for {
		select {
		// If we get a quit signal, finish one more cycle of record processing to catch up before shutting the goroutine down
		case <-tbt.QuitChannel:
			tbt.createEthBalanceRecords(addresses)
			tbt.QuitChannel <- true
			return
		default:
			tbt.createEthBalanceRecords(addresses)
		}
	}
}

func (tbt *AccountTransformer) createEthBalanceRecords(addresses []common.Address) {
	// Cycle through the account addresses
	for _, addr := range addresses {
		// And create a checked_header id for them (IF NOT EXISTS, also this repository uses a lru cache to avoid excess db connections)
		columnID := "account_" + addr.Hex()
		err := tbt.HeaderRepository.AddCheckColumn(columnID)
		if err != nil {
			log.Errorf("transformer: error adding columnID %s", columnID, err)
		}
		// Retrieve the headers which still need eth balance records for this user
		checkedButNotRecordedHeaders, err := tbt.HeaderRepository.MissingMethodsCheckedEventsIntersection(0, -1, []string{columnID}, []string{"token_value_transfers"})
		if err != nil {
			log.Errorf("transformer: error fetching missing headers for %s", columnID, err)
		}
		if len(checkedButNotRecordedHeaders) < 1 {
			continue
		}
		// Create coin balance records for this account at each header
		coinBalanceRecords, err := tbt.AccountPoller.PollAccount(addr, checkedButNotRecordedHeaders)
		if err != nil {
			log.Errorf("transformer: error creating coin balance records for %s", columnID, err)
		}
		// And commit these records to Postgres
		err = tbt.CoinBalanceRepository.CreateCoinBalanceRecords(coinBalanceRecords)
		if err != nil {
			log.Errorf("transformer: error persisting coin balance records for %s", columnID, err)
		}
		// Mark these headers checked for this account
		err = tbt.HeaderRepository.MarkHeadersCheckedForAll(checkedButNotRecordedHeaders, []string{columnID})
		if err != nil {
			log.Errorf("transformer: error marking headers checked for %s", columnID, err)
		}
	}
}
