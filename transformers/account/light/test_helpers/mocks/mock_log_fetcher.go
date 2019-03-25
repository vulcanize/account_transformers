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

package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockFetcher struct {
	Logs                    []types.Log
	PassedHeaders           []core.Header
	PassedContractAddresses []string
}

// Checks all topic0s, on all addresses, fetching matching logs for the given header
func (fetcher *MockFetcher) FetchLogs(contractAddresses []string, topic0s []common.Hash, header core.Header) ([]types.Log, error) {
	fetcher.PassedContractAddresses = contractAddresses
	fetcher.PassedHeaders = append(fetcher.PassedHeaders, header)
	returnLogs := make([]types.Log, 0, len(fetcher.Logs))
	for _, log := range fetcher.Logs {
		if log.BlockNumber == uint64(header.BlockNumber) && checkAllSigs(log, topic0s) {
			returnLogs = append(returnLogs, log)
		}
	}

	return returnLogs, nil
}

func checkAllSigs(log types.Log, topic0s []common.Hash) bool {
	for _, topic0 := range topic0s {
		if log.Topics[0] == topic0 {
			return true
		}
	}

	return false
}
