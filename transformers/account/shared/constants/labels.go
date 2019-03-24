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

package constants

type Label string

const (
	Transfer               Label = "Transfer"
	Mint                   Label = "Mint"
	WipedAccount           Label = "WipedAccount"
	Burn                   Label = "Burn"
	WipeBlacklistedAccount Label = "WipeBlacklistedAccount"
	DestroyedBlackFunds    Label = "DestroyedBlackFunds"
	Issue                  Label = "Issue"
	Redeem                 Label = "Redeem"
	TransferFrom           Label = "TransferFrom"
	Deposit                Label = "Deposit"
	Withdrawal             Label = "Withdrawal"
)

var labels = map[string]Label{
	"Transfer":               Transfer,
	"Mint":                   Mint,
	"WipedAccount":           WipedAccount,
	"Burn":                   Burn,
	"WipeBlacklistedAccount": WipeBlacklistedAccount,
	"DestroyedBlackFunds":    DestroyedBlackFunds,
	"Issue":                  Issue,
	"Redeem":                 Redeem,
	"TransferFrom":           TransferFrom,
	"Deposit":                Deposit,
	"Withdrawal":             Withdrawal,
}

func NewLabel(name string) Label {
	return labels[name]
}

func (label Label) Name() string {
	return string(label)
}
