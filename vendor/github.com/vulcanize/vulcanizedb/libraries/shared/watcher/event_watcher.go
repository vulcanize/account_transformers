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

package watcher

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	chunk "github.com/vulcanize/vulcanizedb/libraries/shared/chunker"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"
	fetch "github.com/vulcanize/vulcanizedb/libraries/shared/fetcher"
	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type EventWatcher struct {
	Transformers  []transformer.EventTransformer
	DB            *postgres.DB
	Fetcher       fetch.LogFetcher
	Chunker       chunk.Chunker
	Addresses     []common.Address
	Topics        []common.Hash
	StartingBlock *int64
}

func NewEventWatcher(db *postgres.DB, bc core.BlockChain) EventWatcher {
	chunker := chunk.NewLogChunker()
	fetcher := fetch.NewFetcher(bc)
	return EventWatcher{
		DB:      db,
		Fetcher: fetcher,
		Chunker: chunker,
	}
}

// Adds transformers to the watcher and updates the chunker, so that it will consider the new transformers.
func (watcher *EventWatcher) AddTransformers(initializers []transformer.EventTransformerInitializer) {
	var contractAddresses []common.Address
	var topic0s []common.Hash
	var configs []transformer.EventTransformerConfig

	for _, initializer := range initializers {
		t := initializer(watcher.DB)
		watcher.Transformers = append(watcher.Transformers, t)

		config := t.GetConfig()
		configs = append(configs, config)

		if watcher.StartingBlock == nil {
			watcher.StartingBlock = &config.StartingBlockNumber
		} else if earlierStartingBlockNumber(config.StartingBlockNumber, *watcher.StartingBlock) {
			watcher.StartingBlock = &config.StartingBlockNumber
		}

		addresses := transformer.HexStringsToAddresses(config.ContractAddresses)
		contractAddresses = append(contractAddresses, addresses...)
		topic0s = append(topic0s, common.HexToHash(config.Topic))
	}

	watcher.Addresses = append(watcher.Addresses, contractAddresses...)
	watcher.Topics = append(watcher.Topics, topic0s...)
	watcher.Chunker.AddConfigs(configs)
}

func (watcher *EventWatcher) Execute(recheckHeaders constants.TransformerExecution) error {
	if watcher.Transformers == nil {
		return fmt.Errorf("No transformers added to watcher")
	}

	checkedColumnNames, err := repo.GetCheckedColumnNames(watcher.DB)
	if err != nil {
		return err
	}
	notCheckedSQL := repo.CreateNotCheckedSQL(checkedColumnNames, recheckHeaders)

	missingHeaders, err := repo.MissingHeaders(*watcher.StartingBlock, -1, watcher.DB, notCheckedSQL)
	if err != nil {
		log.Error("Fetching of missing headers failed in watcher!")
		return err
	}

	for _, header := range missingHeaders {
		// TODO Extend FetchLogs for doing several blocks at a time
		logs, err := watcher.Fetcher.FetchLogs(watcher.Addresses, watcher.Topics, header)
		if err != nil {
			// TODO Handle fetch error in watcher
			log.Errorf("Error while fetching logs for header %v in watcher", header.Id)
			return err
		}

		chunkedLogs := watcher.Chunker.ChunkLogs(logs)

		// Can't quit early and mark as checked if there are no logs. If we are running continuousLogSync,
		// not all logs we're interested in might have been fetched.
		for _, t := range watcher.Transformers {
			transformerName := t.GetConfig().TransformerName
			logChunk := chunkedLogs[transformerName]
			err = t.Execute(logChunk, header, constants.HeaderMissing)
			if err != nil {
				log.Errorf("%v transformer failed to execute in watcher: %v", transformerName, err)
				return err
			}
		}
	}
	return err
}

func earlierStartingBlockNumber(transformerBlock, watcherBlock int64) bool {
	return transformerBlock < watcherBlock
}
