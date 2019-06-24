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

package fakes

import (
	"encoding/json"

	"github.com/vulcanize/vulcanizedb/pkg/core"
)

var rawFakeHeader, _ = json.Marshal(core.Header{})

var FakeHeader1 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad123ert",
	BlockNumber: 6791667,
	Raw:         rawFakeHeader,
	Timestamp:   "50000000",
	Id:          1,
}

var FakeHeader2 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad456yui",
	BlockNumber: 6791668,
	Raw:         rawFakeHeader,
	Timestamp:   "50000015",
	Id:          2,
}

var FakeHeader3 = core.Header{
	Hash:        "0x135391a0962a63944e5908e6fedfff90fb4be3e3290a21017861099bad234hfs",
	BlockNumber: 6791669,
	Raw:         rawFakeHeader,
	Timestamp:   "50000030",
	Id:          3,
}
