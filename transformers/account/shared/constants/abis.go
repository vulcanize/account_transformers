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

// These abis have been modifed so that they all use the name "from" for the address which is losing tokens,
// "to" for the address which is getting tokens, and "value" for the number of tokens in every case
// This way we can generically access these values we need for tracking balance
const (
	Transfer2Indexed               string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer2Indexed","type":"event"}]`   // event Transfer(address indexed from, address indexed to, uint256 value);
	Transfer0Indexed               string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer0Indexed","type":"event"}]` // event Transfer(address from, address to, uint256 value);
	Transfer1Indexed               string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer1Indexed","type":"event"}]`  // event Transfer(address indexed from, address to, uint256 value);
	Transfer3Indexed               string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"Transfer3Indexed","type":"event"}]`    // event Transfer(address indexed from, address indexed to, uint256 indexed value);
	Mint1Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Mint1Indexed","type":"event"}]`
	Mint0Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Mint0Indexed","type":"event"}]`
	Mint2Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"Mint2Indexed","type":"event"}]`
	WipedAccount1Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"WipedAccount1Indexed","type":"event"}]`
	WipedAccount0Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"WipedAccount0Indexed","type":"event"}]`
	WipedAccount2Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"WipedAccount2Indexed","type":"event"}]`
	Burn1Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn1Indexed","type":"event"}]`
	Burn0Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn0Indexed","type":"event"}]`
	Burn2Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"Burn2Indexed","type":"event"}]`
	WipeBlacklistedAccount1Indexed string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"WipeBlacklistedAccount1Indexed","type":"event"}]`
	WipeBlacklistedAccount0Indexed string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"WipeBlacklistedAccount0Indexed","type":"event"}]`
	WipeBlacklistedAccount2Indexed string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"WipeBlacklistedAccount2Indexed","type":"event"}]`
	DestroyedBlackFunds0Indexed    string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"DestroyedBlackFunds0Indexed","type":"event"}]`
	DestroyedBlackFunds1Indexed    string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"DestroyedBlackFunds1Indexed","type":"event"}]`
	DestroyedBlackFunds2Indexed    string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"value","type":"uint256"}],"name":"DestroyedBlackFunds2Indexed","type":"event"}]`
	Issue0Indexed                  string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"value","type":"uint256"}],"name":"Issue0Indexed","type":"event"}]`
	Issue1Indexed                  string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"value","type":"uint256"}],"name":"Issue1Indexed","type":"event"}]`
	Redeem0Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"value","type":"uint256"}],"name":"Redeem0Indexed","type":"event"}]`
	Redeem1Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"value","type":"uint256"}],"name":"Redeem1Indexed","type":"event"}]`
	TransferFrom3Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom3Indexed","type":"event"}]`
	TransferFrom2Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom2Indexed","type":"event"}]`
	TransferFrom1Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom1Indexed","type":"event"}]`
	TransferFrom0Indexed           string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"spender","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferFrom0Indexed","type":"event"}]`
	Deposit1Indexed                string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit1Indexed","type":"event"}]`
	Deposit2Indexed                string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"wad","type":"uint256"}],"name":"Deposit2Indexed","type":"event"}]`
	Deposit0Indexed                string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Deposit0Indexed","type":"event"}]`
	Withdrawal1Indexed             string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdrawal1Indexed","type":"event"}]`
	Withdrawal2Indexed             string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"wad","type":"uint256"}],"name":"Withdrawal2Indexed","type":"event"}]`
	Withdrawal0Indexed             string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"wad","type":"uint256"}],"name":"Withdrawal0Indexed","type":"event"}]`
	Sent3Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Sent3Indexed","type":"event"}]`
	Sent2Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Sent2Indexed","type":"event"}]`
	Sent1Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Sent1Indexed","type":"event"}]`
	Sent0Indexed                   string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"operator","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Sent0Indexed","type":"event"}]`
	Minted3Indxed                  string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Minted3Indxed","type":"event"}]`
	Minted2Indxed                  string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Minted2Indxed","type":"event"}]`
	Minted1Indxed                  string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Minted1Indxed","type":"event"}]`
	Minted0Indxed                  string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"operator","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Minted0Indxed","type":"event"}]`
	Burned3Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Burned3Indexed","type":"event"}]`
	Burned2Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Burned2Indexed","type":"event"}]`
	Burned1Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Burned1Indexed","type":"event"}]`
	Burned0Indexed                 string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"operator","type":"address"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"},{"indexed":false,"name":"operatorData","type":"bytes"}],"name":"Burned0Indexed","type":"event"}]`
	TransferWithData3Indexed       string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":true,"name":"data","type":"bytes"}],"name":"TransferWithData3Indexed","type":"event"}]`
	TransferWithData2Indexed       string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"}],"name":"TransferWithData2Indexed","type":"event"}]`
	TransferWithData1Indexed       string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"}],"name":"TransferWithData1Indexed","type":"event"}]`
	TransferWithData0Indexed       string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"data","type":"bytes"}],"name":"TransferWithData0Indexed","type":"event"}]`
	NewTokenGrant3Indexed          string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"value","type":"uint256"},{"indexed":false,"name":"start","type":"uint256"},{"indexed":false,"name":"cliff","type":"uint256"},{"indexed":false,"name":"vesting","type":"uint256"}],"name":"NewTokenGrant3Indexed","type":"event"}]`
	NewTokenGrant2Indexed          string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"start","type":"uint256"},{"indexed":false,"name":"cliff","type":"uint256"},{"indexed":false,"name":"vesting","type":"uint256"}],"name":"NewTokenGrant2Indexed","type":"event"}]`
	NewTokenGrant1Indexed          string = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"start","type":"uint256"},{"indexed":false,"name":"cliff","type":"uint256"},{"indexed":false,"name":"vesting","type":"uint256"}],"name":"NewTokenGrant1Indexed","type":"event"}]`
	NewTokenGrant0Indexed          string = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"},{"indexed":false,"name":"start","type":"uint256"},{"indexed":false,"name":"cliff","type":"uint256"},{"indexed":false,"name":"vesting","type":"uint256"}],"name":"NewTokenGrant0Indexed","type":"event"}]`
)

var CombinedABI = Transfer3Indexed[:len(Transfer3Indexed)-1] + "," +
	Transfer2Indexed[1:len(Transfer2Indexed)-1] + "," +
	Transfer1Indexed[1:len(Transfer1Indexed)-1] + "," +
	Transfer0Indexed[1:len(Transfer0Indexed)-1] + "," +
	Mint2Indexed[1:len(Mint2Indexed)-1] + "," +
	Mint1Indexed[1:len(Mint1Indexed)-1] + "," +
	Mint0Indexed[1:len(Mint0Indexed)-1] + "," +
	WipedAccount2Indexed[1:len(WipedAccount2Indexed)-1] + "," +
	WipedAccount1Indexed[1:len(WipedAccount1Indexed)-1] + "," +
	WipedAccount0Indexed[1:len(WipedAccount0Indexed)-1] + "," +
	Burn2Indexed[1:len(Burn2Indexed)-1] + "," +
	Burn1Indexed[1:len(Burn1Indexed)-1] + "," +
	Burn0Indexed[1:len(Burn0Indexed)-1] + "," +
	WipeBlacklistedAccount2Indexed[1:len(WipeBlacklistedAccount2Indexed)-1] + "," +
	WipeBlacklistedAccount1Indexed[1:len(WipeBlacklistedAccount1Indexed)-1] + "," +
	WipeBlacklistedAccount0Indexed[1:len(WipeBlacklistedAccount0Indexed)-1] + "," +
	DestroyedBlackFunds2Indexed[1:len(DestroyedBlackFunds2Indexed)-1] + "," +
	DestroyedBlackFunds1Indexed[1:len(DestroyedBlackFunds1Indexed)-1] + "," +
	DestroyedBlackFunds0Indexed[1:len(DestroyedBlackFunds0Indexed)-1] + "," +
	DestroyedBlackFunds2Indexed[1:len(DestroyedBlackFunds2Indexed)-1] + "," +
	DestroyedBlackFunds1Indexed[1:len(DestroyedBlackFunds1Indexed)-1] + "," +
	DestroyedBlackFunds0Indexed[1:len(DestroyedBlackFunds0Indexed)-1] + "," +
	TransferFrom3Indexed[1:len(TransferFrom3Indexed)-1] + "," +
	TransferFrom2Indexed[1:len(TransferFrom2Indexed)-1] + "," +
	TransferFrom1Indexed[1:len(TransferFrom1Indexed)-1] + "," +
	TransferFrom0Indexed[1:len(TransferFrom0Indexed)-1] + "," +
	Deposit2Indexed[1:len(Deposit2Indexed)-1] + "," +
	Deposit1Indexed[1:len(Deposit1Indexed)-1] + "," +
	Deposit0Indexed[1:len(Deposit0Indexed)-1] + "," +
	Withdrawal2Indexed[1:len(Withdrawal2Indexed)-1] + "," +
	Withdrawal1Indexed[1:len(Withdrawal1Indexed)-1] + "," +
	Withdrawal0Indexed[1:len(Withdrawal0Indexed)-1] + "," +
	Sent3Indexed[1:len(Sent3Indexed)-1] + "," +
	Sent2Indexed[1:len(Sent2Indexed)-1] + "," +
	Sent1Indexed[1:len(Sent1Indexed)-1] + "," +
	Sent0Indexed[1:len(Sent0Indexed)-1] + "," +
	Minted3Indxed[1:len(Minted3Indxed)-1] + "," +
	Minted2Indxed[1:len(Minted2Indxed)-1] + "," +
	Minted1Indxed[1:len(Minted1Indxed)-1] + "," +
	Minted0Indxed[1:len(Minted0Indxed)-1] + "," +
	Burned3Indexed[1:len(Burned3Indexed)-1] + "," +
	Burned2Indexed[1:len(Burned2Indexed)-1] + "," +
	Burned1Indexed[1:len(Burned1Indexed)-1] + "," +
	Burned0Indexed[1:len(Burned0Indexed)-1] + "," +
	TransferWithData3Indexed[1:len(TransferWithData3Indexed)-1] + "," +
	TransferWithData2Indexed[1:len(TransferWithData2Indexed)-1] + "," +
	TransferWithData1Indexed[1:len(TransferWithData1Indexed)-1] + "," +
	TransferWithData0Indexed[1:len(TransferWithData0Indexed)-1] + "," +
	NewTokenGrant3Indexed[1:len(NewTokenGrant3Indexed)-1] + "," +
	NewTokenGrant2Indexed[1:len(NewTokenGrant2Indexed)-1] + "," +
	NewTokenGrant1Indexed[1:len(NewTokenGrant1Indexed)-1] + "," +
	NewTokenGrant0Indexed[1:]
