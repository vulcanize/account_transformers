// VulcanizeDB
// Copyright © 2019 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTstringLITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package constants

type Event struct {
	Label Label    // The Label for this event
	Names []string // Name used to find this event in the ABI, depending on how many indexed args it has (e.g. Transfer3Indexed)
}

func (e *Event) GetLabel() Label {
	return e.Label
}

func (e *Event) GetName(i int) string {
	if i > len(e.Names)-1 || i < 0 {
		return ""
	}
	return e.Names[i]
}

func GetEventFromLabel(label Label) Event {
	if event, ok := events[label]; ok {
		return event
	}
	return Event{}
}

func GetEventFromName(name string) Event {
	for label, event := range events {
		if label.Name() == name {
			return event
		}
	}
	return Event{}
}

var events = map[Label]Event{
	Transfer: {
		Label: Transfer,
		Names: []string{
			"Transfer0Indexed",
			"Transfer1Indexed",
			"Transfer2Indexed",
			"Transfer3Indexed",
		},
	},
	Mint: {
		Label: Mint,
		Names: []string{
			"Mint0Indexed",
			"Mint1Indexed",
			"Mint2Indexed",
		},
	},
	WipedAccount: {
		Label: WipedAccount,
		Names: []string{
			"WipedAccount0Indexed",
			"WipedAccount1Indexed",
			"WipedAccount2Indexed",
		},
	},
	Burn: {
		Label: Burn,
		Names: []string{
			"Burn0Indexed",
			"Burn1Indexed",
			"Burn2Indexed",
		},
	},
	WipeBlacklistedAccount: {
		Label: WipeBlacklistedAccount,
		Names: []string{
			"WipeBlacklistedAccount0Indexed",
			"WipeBlacklistedAccount1Indexed",
			"WipeBlacklistedAccount2Indexed",
		},
	},
	DestroyedBlackFunds: {
		Label: DestroyedBlackFunds,
		Names: []string{
			"DestroyedBlackFunds0Indexed",
			"DestroyedBlackFunds1Indexed",
			"DestroyedBlackFunds2Indexed",
		},
	},
	Issue: {
		Label: Issue,
		Names: []string{
			"Issue0Indexed",
			"Issue1Indexed",
		},
	},
	Redeem: {
		Label: Redeem,
		Names: []string{
			"Redeem0Indexed",
			"Redeem1Indexed",
		},
	},
	TransferFrom: {
		Label: TransferFrom,
		Names: []string{
			"TransferFrom0Indexed",
			"TransferFrom1Indexed",
			"TransferFrom2Indexed",
			"TransferFrom3Indexed",
		},
	},
	Deposit: {
		Label: Deposit,
		Names: []string{
			"Deposit0Indexed",
			"Deposit1Indexed",
			"Deposit2Indexed",
		},
	},
	Withdrawal: {
		Label: Withdrawal,
		Names: []string{
			"Withdrawal0Indexed",
			"Withdrawal1Indexed",
			"Withdrawal2Indexed",
		},
	},
	Sent: {
		Label: Sent,
		Names: []string{
			"Sent0Indexed",
			"Sent1Indexed",
			"Sent2Indexed",
			"Sent3Indexed",
		},
	},
	Minted: {
		Label: Minted,
		Names: []string{
			"Minted0Indxed",
			"Minted1Indxed",
			"Minted2Indxed",
			"Minted3Indxed",
		},
	},
	Burned: {
		Label: Burned,
		Names: []string{
			"Burned0Indexed",
			"Burned1Indexed",
			"Burned2Indexed",
			"Burned3Indexed",
		},
	},
	TransferWithData: {
		Label: TransferWithData,
		Names: []string{
			"TransferWithData0Indexed",
			"TransferWithData1Indexed",
			"TransferWithData2Indexed",
			"TransferWithData3Indexed",
		},
	},
	NewTokenGrant: {
		Label: NewTokenGrant,
		Names: []string{
			"NewTokenGrant0Indexed",
			"NewTokenGrant1Indexed",
			"NewTokenGrant2Indexed",
			"NewTokenGrant3Indexed",
		},
	},
}
