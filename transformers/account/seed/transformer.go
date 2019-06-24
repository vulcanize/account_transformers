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

package seed

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/streamer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/header/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"

	"github.com/vulcanize/account_transformers/transformers/account/light/converters"
	"github.com/vulcanize/account_transformers/transformers/account/light/repositories"
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
)

type AccountTransformer struct {
	// Streamer interface is used to subscribe to a ResponsePayloads from the seed node
	Streamer streamer.IStreamer
	// SubscriptionConfig contains the subscription parameters to send to the seed node to direct which data should be relayed
	SubscriptionConfig config.Subscription
	// HeaderRepository provides access to the public.headers table
	HeaderRepository repository.HeaderRepository
	// AddressRepository provides access to the account.addresses table
	AddressRepository repositories.AddressRepository
	// WatchedContractRepository provides access to the accounts.contract_addresses table
	WatchedContractRepository repositories.WatchedContractRepository
	// Decode seed node payloads
	Decoder Decoder
	// Interface for persist decoded payload
	Repository Repository
	// Interface for converting token value transfer events
	TokenConverter converters.ValueTransferConverter
	// Interface for persisting token value transfer records
	TokenRepository repositories.ValueTransferEventRepository

	NextStart     int64
	routine       bool
	StreamChannel chan ipfs.ResponsePayload
	QuitChannel   chan bool
	WaitGroup     *sync.WaitGroup
}

func (tbt AccountTransformer) NewTransformer(db *postgres.DB, subCon config.Subscription, client core.RpcClient) transformer.SeedNodeTransformer {
	return &AccountTransformer{
		SubscriptionConfig:        subCon,
		Streamer:                  streamer.NewSeedStreamer(client),
		HeaderRepository:          repository.NewHeaderRepository(db),
		AddressRepository:         repositories.NewAddressRepository(db),
		WatchedContractRepository: repositories.NewWatchedContractRepository(db),
		TokenRepository:           repositories.NewValueTransferEventRepository(db),
		Decoder:                   NewSeedNodeDecoder(),
		QuitChannel:               make(chan bool, 1),
		StreamChannel:             make(chan ipfs.ResponsePayload, 20000),
	}
}

func (tbt AccountTransformer) Init() error {
	var err error
	// Get the list of token addresses we want to watch for from the config and add them to Postgres
	configuredTokenAddresses := constants.TokenAddresses()
	for _, addr := range configuredTokenAddresses {
		err = tbt.WatchedContractRepository.AddAddress(addr)
		if err != nil {
			return err
		}
	}
	// Configure token converter with abis and token equivalencies
	tbt.TokenConverter, err = converters.NewValueTransferConverter(constants.CombinedABI, constants.EquivalentTokenAddressesMapping())
	return err
}

func (tbt *AccountTransformer) Execute() error {
	// Subscribe to the seed node service with the given config/filter parameters
	sub, err := tbt.Streamer.Stream(tbt.StreamChannel, tbt.SubscriptionConfig)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case response := <-tbt.StreamChannel:
			decodedResponse, err := tbt.Decoder.DecodePayload(response)
			if err != nil {
				log.Error(err)
			}
			headerId, err := tbt.Repository.CommitPayload(decodedResponse)
			if err != nil {
				log.Error(err)
			}
			tokenTransfers, err := tbt.TokenConverter.Convert(decodedResponse.Logs, headerId)
			if err != nil {
				log.Error(err)
			}
			log.Error(tbt.TokenRepository.CreateTokenValueTransferRecords(tokenTransfers))
		case <-tbt.QuitChannel:
			return nil
		case err := <-sub.Err():
			return err
		}
	}
}

func (tbt *AccountTransformer) GetConfig() config.Subscription {
	return tbt.SubscriptionConfig
}
