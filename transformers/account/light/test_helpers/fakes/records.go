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
	"github.com/vulcanize/account_transformers/transformers/account/light/models"
	"github.com/vulcanize/account_transformers/transformers/account/shared"
)

var FakeCoinRecord1 = shared.CoinBalanceRecord{
	Address:     string("FakeAddress"),
	BlockNumber: 6791667,
	Value:       "12345",
}

var FakeCoinRecord2 = shared.CoinBalanceRecord{
	Address:     string("FakeAddress"),
	BlockNumber: 6791668,
	Value:       "10000",
}

var FakeCoinRecord3 = shared.CoinBalanceRecord{
	Address:     string("FakeAddress"),
	BlockNumber: 6791669,
	Value:       "10000",
}

var FakeTokenRecord1 = shared.TokenBalanceRecord{
	Address:         string("FakeAddress"),
	BlockNumber:     6791667,
	ContractAddress: string("FakeContract"),
	Value:           "400",
}

var FakeTokenRecord2 = shared.TokenBalanceRecord{
	Address:         string("FakeAddress"),
	BlockNumber:     6791668,
	ContractAddress: string("FakeContract"),
	Value:           "400",
}

var FakeTokenRecord3 = shared.TokenBalanceRecord{
	Address:         string("FakeAddress"),
	BlockNumber:     6791669,
	ContractAddress: string("FakeContract"),
	Value:           "1400",
}

var FakeValueTransferRecord1 = models.ValueTransferModel{}

var FakeValueTransferRecord2 = models.ValueTransferModel{}

var FakeValueTransferRecord3 = models.ValueTransferModel{}
