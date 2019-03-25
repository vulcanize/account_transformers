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
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockHeaderRepository struct {
	CheckColumnIds            []string
	CheckedHeaders            map[string][]int64
	ReturnErr                 error
	ReturnMissingHeaders      map[string][]core.Header
	ReturnIntersectionHeaders []core.Header
	PassedStartingBlock       int64
	PassedEndingBlock         int64
}

// Adds a checked_header column for the provided column id
func (r *MockHeaderRepository) AddCheckColumn(id string) error {
	r.CheckColumnIds = append(r.CheckColumnIds, id)
	return r.ReturnErr
}

// Adds a checked_header column for all of the provided column ids
func (r *MockHeaderRepository) AddCheckColumns(ids []string) error {
	panic("implement me")
}

// Marks the header checked for the provided column id
func (r *MockHeaderRepository) MarkHeaderChecked(headerID int64, id string) error {
	r.CheckedHeaders[id] = append(r.CheckedHeaders[id], headerID)
	return r.ReturnErr
}

// Marks the header checked for all of the provided column ids
func (r *MockHeaderRepository) MarkHeaderCheckedForAll(headerID int64, ids []string) error {
	panic("implement me")
}

// Marks all of the provided headers checked for each of the provided column ids
func (r *MockHeaderRepository) MarkHeadersCheckedForAll(headers []core.Header, ids []string) error {
	panic("implement me")
}

// Returns missing headers for the provided checked_headers column id
func (r *MockHeaderRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64, id string) ([]core.Header, error) {
	r.PassedEndingBlock = endingBlockNumber
	r.PassedStartingBlock = startingBlockNumber
	return r.ReturnMissingHeaders[id], r.ReturnErr
}

// Returns missing headers for all of the provided checked_headers column ids
func (r *MockHeaderRepository) MissingHeadersForAll(startingBlockNumber, endingBlockNumber int64, ids []string) ([]core.Header, error) {
	panic("implement me")
}

// Takes in an ordered sequence of headers and returns only the first contiguous segment
// Enforce continuity with previous segment with the appropriate startingBlockNumber
func contiguousHeaders(headers []core.Header, startingBlockNumber int64) []core.Header {
	if len(headers) < 1 {
		return headers
	}
	previousHeader := headers[0].BlockNumber
	if previousHeader != startingBlockNumber {
		return []core.Header{}
	}
	for i := 1; i < len(headers); i++ {
		previousHeader++
		if headers[i].BlockNumber != previousHeader {
			return headers[:i]
		}
	}

	return headers
}

// Returns headers that have been checked for all of the provided event ids but not for the provided method ids
func (r *MockHeaderRepository) MissingMethodsCheckedEventsIntersection(startingBlockNumber, endingBlockNumber int64, methodIds, eventIds []string) ([]core.Header, error) {
	r.PassedEndingBlock = endingBlockNumber
	r.PassedStartingBlock = startingBlockNumber
	return r.ReturnIntersectionHeaders, r.ReturnErr
}

// Check the repositories column id cache for a value
func (r *MockHeaderRepository) CheckCache(key string) (interface{}, bool) {
	panic("implement me")
}
