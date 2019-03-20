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

package constants

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/filters"
)

var Filters = []filters.LogFilter{
	{
		Name: Transfer.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Transfer.Sig()},
	},
	{
		Name: TransferDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{TransferDev.Sig()},
	},
	{
		Name: Mint.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Mint.Sig()},
	},
	{
		Name: MintDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{MintDev.Sig()},
	},
	{
		Name: WipedAccount.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{WipedAccount.Sig()},
	},
	{
		Name: WipedAccountDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{WipedAccountDev.Sig()},
	},
	{
		Name: Burn.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Burn.Sig()},
	},
	{
		Name: BurnDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{BurnDev.Sig()},
	},
	{
		Name: Burn.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Burn.Sig()},
	},
	{
		Name: BurnDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{BurnDev.Sig()},
	},
	{
		Name: WipeBlacklistedAccount.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{WipeBlacklistedAccount.Sig()},
	},
	{
		Name: WipeBlacklistedAccountDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{WipeBlacklistedAccountDev.Sig()},
	},
	{
		Name: DestroyedBlackFunds.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{DestroyedBlackFunds.Sig()},
	},
	{
		Name: DestroyedBlackFundsDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{DestroyedBlackFundsDev.Sig()},
	},
	{
		Name: Issue.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Issue.Sig()},
	},
	{
		Name: IssueDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{IssueDev.Sig()},
	},
	{
		Name: Redeem.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{Redeem.Sig()},
	},
	{
		Name: RedeemDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{RedeemDev.Sig()},
	},
	{
		Name: TransferFrom.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{TransferFrom.Sig()},
	},
	{
		Name: TransferFromDev.String(),
		FromBlock: 0,
		ToBlock: -1,
		Topics: core.Topics{TransferFromDev.Sig()},
	},
}
