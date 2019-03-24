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

package shared

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransferEntity struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type MintEntity struct {
	To     common.Address
	Amount *big.Int
}

type WipedAccountEntity struct {
	Account common.Address
	Balance *big.Int
}

type BurnEntity struct {
	Burner common.Address
	Value  *big.Int
}

type WipeBlacklistedAccountEntity struct {
	Account common.Address
	Balance *big.Int
}

type DestroyedBlackFundsEntity struct {
	BlackListedUser common.Address
	Balance         *big.Int
}

type IssueEntity struct {
	Amount *big.Int
}

type RedeemEntity struct {
	Amount *big.Int
}

type TransferFromEntity struct {
	Spender common.Address
	From    common.Address
	To      common.Address
	Value   *big.Int
}

type DepositEntity struct {
	Dst common.Address
	Wad *big.Int
}

type WithdrawalEntity struct {
	Src common.Address
	Wad *big.Int
}
