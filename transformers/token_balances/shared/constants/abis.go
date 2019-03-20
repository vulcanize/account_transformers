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

var ABIs = map[Label]string{
	Transfer : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`,
	TransferDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint"}],"name":"Transfer","type":"event"}]`,
	Mint : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"Mint","type":"event"}]`,
	MintDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint"}],"name":"Mint","type":"event"}]`,
	WipedAccount : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipedAccount","type":"event"}]`,
	WipedAccountDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint"}],"name":"WipedAccount","type":"event"}]`,
	Burn : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"burner","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`,
	BurnDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"burner","type":"address"},{"indexed":false,"name":"value","type":"uint"}],"name":"Burn","type":"event"}]`,
	WipeBlacklistedAccount : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipeBlacklistedAccount","type":"event"}]`,
	WipeBlacklistedAccountDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint"}],"name":"WipeBlacklistedAccount","type":"event"}]`,
	DestroyedBlackFunds : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"blackListedUser","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"}]`,
	DestroyedBlackFundsDev : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"blackListedUser","type":"address"},{"indexed":false,"name":"balance","type":"uint"}],"name":"DestroyedBlackFunds","type":"event"}]`,
	Issue : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Issue","type":"event"}]`,
	IssueDev : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint"}],"name":"Issue","type":"event"}]`,
	Redeem : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Redeem","type":"event"}]`,
	RedeemDev : `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint"}],"name":"Redeem","type":"event"}]`,
	TransferFrom : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom","type":"event"}]`,
	TransferFromDev : `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint"}],"name":"TransferFrom","type":"event"}]`,
}

func (label Label) ABI() string {
	return ABIs[label]
}