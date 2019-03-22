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
package config

import (
	"github.com/vulcanize/vulcanizedb/pkg/config"
)

var MainnetAccountConfig = config.ContractConfig{
	Name:    "Accout-mainnet",
	Network: "",
}

var RopstenAccountConfig = config.ContractConfig{
	Name:    "Account-ropsten",
	Network: "ropsten",
}

var RinkebyAccountConfig = config.ContractConfig{
	Name:    "Account-rinkeby",
	Network: "ropsten",
}

var KovanccountConfig = config.ContractConfig{
	Name:    "Account-kovan",
	Network: "kovan",
}
