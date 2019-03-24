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

type ABI string

const (
	Transfer2Indexed               ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`   // event Transfer(address indexed from, address indexed to, uint256 value);
	Transfer0Indexed               ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]` // event Transfer(address from, address to, uint256 value);
	Transfer1Indexed               ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`  // event Transfer(address indexed from, address to, uint256 value);
	Transfer3Indexed               ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`    // event Transfer(address indexed from, address indexed to, uint256 indexed value);
	Mint1Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"Mint","type":"event"}]`
	Mint0Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"Mint","type":"event"}]`
	Mint2Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"amount","type":"uint256"}],"name":"Mint","type":"event"}]`
	WipedAccount1Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipedAccount","type":"event"}]`
	WipedAccount0Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipedAccount","type":"event"}]`
	WipedAccount2Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":true,"name":"balance","type":"uint256"}],"name":"WipedAccount","type":"event"}]`
	Burn1Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"burner","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`
	Burn0Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"burner","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`
	Burn2Indexed                   ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"burner","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]`
	WipeBlacklistedAccount1Indexed ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipeBlacklistedAccount","type":"event"}]`
	WipeBlacklistedAccount0Indexed ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"account","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"WipeBlacklistedAccount","type":"event"}]`
	WipeBlacklistedAccount2Indexed ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":true,"name":"balance","type":"uint256"}],"name":"WipeBlacklistedAccount","type":"event"}]`
	DestroyedBlackFunds0Indexed    ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"blackListedUser","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"}]`
	DestroyedBlackFunds1Indexed    ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"blackListedUser","type":"address"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"}]`
	DestroyedBlackFunds2Indexed    ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"blackListedUser","type":"address"},{"indexed":true,"name":"balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"}]`
	Issue0Indexed                  ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Issue","type":"event"}]`
	Issue1Indexed                  ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"amount","type":"uint256"}],"name":"Issue","type":"event"}]`
	Redeem0Indexed                 ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Redeem","type":"event"}]`
	RedeemDev1Indexed              ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"amount","type":"uint256"}],"name":"Redeem","type":"event"}]`
	TransferFrom3Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom","type":"event"}]`
	TransferFrom2Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom","type":"event"}]`
	TransferFrom1Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom","type":"event"}]`
	TransferFrom0Indexed           ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"spender","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom","type":"event"}]`
	Deposit1Indexed                ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"dst","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit","type":"event"}]`
	Deposit2Indexed                ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"dst","type":"address"},{"indexed":true,"name":"wad","type":"uint256"}],"name":"Deposit","type":"event"}]`
	Deposit0Indexed                ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"dst","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit","type":"event"}]`
	Withdrawal1Indexed             ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"src","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdrawal","type":"event"}]`
	Withdrawal2Indexed             ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"src","type":"address"},{"indexed":true,"name":"wad","type":"uint256"}],"name":"Withdrawal","type":"event"}]`
	Withdrawal0Indexed             ABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"src","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdrawal","type":"event"}]`
)

func (abi ABI) String() string {
	return string(abi)
}
