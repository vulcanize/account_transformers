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

import "strings"

type Label string

const (
	Transfer Label = "Transfer"
	TransferDev Label = "TransferDeviation"
	Mint Label = "Mint"
	MintDev Label = "MintDeviation"
	WipedAccount Label = "WipedAccount"
	WipedAccountDev Label = "WipedAccountDeviation"
	Burn Label = "Burn"
	BurnDev Label = "BurnDeviation"
	WipeBlacklistedAccount Label = "WipeBlacklistedAccount"
	WipeBlacklistedAccountDev Label = "WipeBlacklistedAccountDeviation"
	DestroyedBlackFunds Label = "DestroyedBlackFunds"
	DestroyedBlackFundsDev Label = "DestroyedBlackFundsDeviation"
	Issue Label = "Issue"
	IssueDev Label = "IssueDeviation"
	Redeem Label = "Redeem"
	RedeemDev Label = "RedeemDeviation"
	TransferFrom Label = "TransferFrom"
	TransferFromDev Label = "TransferFromDeviation"
)

var labels = map[string]Label{
	"Transfer" : Transfer,
	"TransferDeviation" : TransferDev,
	"Mint" : Mint,
	"MintDeviation" : MintDev,
	"WipedAccount" : WipedAccount,
	"WipedAccountDeviation" : WipedAccountDev,
	"Burn" : Burn,
	"BurnDeviation" : BurnDev,
	"WipeBlacklistedAccount" : WipeBlacklistedAccount,
	"WipeBlacklistedAccountDeviation" : WipeBlacklistedAccountDev,
	"DestroyedBlackFunds" : DestroyedBlackFunds,
	"DestroyedBlackFundsDeviation" : DestroyedBlackFundsDev,
	"Issue" : Issue,
	"IssueDeviation" : IssueDev,
	"Redeem" : Redeem,
	"RedeemDeviation" : RedeemDev,
	"TransferFrom" : TransferFrom,
	"TransferFromDeviation" : TransferFromDev,
}

func NewLabel(name string) Label {
	return labels[name]
}

func (label Label) String() string {
	return string(label)
}

func (label Label) Event() string {
	return strings.Replace(label.String(), "Deviation", "", 1)
}