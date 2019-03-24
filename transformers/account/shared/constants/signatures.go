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

import "github.com/ethereum/go-ethereum/common"

var EventSignatures = map[string]Label{
	"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": Transfer,
	"0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377": Transfer,
	"0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885": Mint,
	"0xb4683f1a6bfdb19078899ae88702d383183e50a367af29a37f4e7357a7fc13f9": Mint,
	"0xdf58d2368c06216a398f05a7a88c8edc64a25c33f33fd2bd8b56fbc8822c02d8": WipedAccount,
	"0x6dd8c32d75aebbe065ce6cbbce463910d4d71fcfa91bb7fb276358b35223b8c1": WipedAccount,
	"0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5": Burn,
	"0x654fd2845115fd0c81b8dd4d62aa4344b1990142afe190faf6dc30a9fa0f5e31": Burn,
	"0xfa8f14973a436f651cdc72fcb50527f364a3b92681dc7aacb0ebeed1e7fb7070": WipeBlacklistedAccount,
	"0x5ad932ee61740a6d8228a885eb9b654a6495a14372f49fed16686c4c80cc031a": WipeBlacklistedAccount,
	"0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6": DestroyedBlackFunds,
	"0x2ca1f69d129207da56149d0f73a0e0dce0d262588ea7da9e1523fc6fe04ad407": DestroyedBlackFunds,
	"0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a": Issue,
	"0xbc5846b4d860cbc33f6f6b6b2f7648fcbaca7425ffa5a8e26a9a70e5fd092f49": Issue,
	"0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44": Redeem,
	"0x13ffabc1e5c4de958be06f1541c4a24fdf8a4deffcbf851b2e805e41613e5e22": Redeem,
	"0x5f7542858008eeb041631f30e6109ae94b83a58e9a58261dd2c42c508850f939": TransferFrom,
	"0xbf480be9a10591c488f4bdf0305c23b7114779b0c1a42fb58401ef0a8574bc79": TransferFrom,
	"0x9b40ecb08ede8b5f14f7401641540771bbdaaac8638c21f5813ae7ba76d75155": Deposit,
	"0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c": Deposit,
	"0x1be94c6778a9e1751832385a51994fbb7b20c9c08ebfa22735a951d4b84ebb1e": Withdrawal,
	"0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65": Withdrawal,
}

var Topic0s = []common.Hash{
	common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), // Transfer
	common.HexToHash("0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377"), // Transfer
	common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"), // Mint
	common.HexToHash("0xb4683f1a6bfdb19078899ae88702d383183e50a367af29a37f4e7357a7fc13f9"), // Mint
	common.HexToHash("0xdf58d2368c06216a398f05a7a88c8edc64a25c33f33fd2bd8b56fbc8822c02d8"), // WipedAccount
	common.HexToHash("0x6dd8c32d75aebbe065ce6cbbce463910d4d71fcfa91bb7fb276358b35223b8c1"), // WipedAccount
	common.HexToHash("0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5"), // Burn
	common.HexToHash("0x654fd2845115fd0c81b8dd4d62aa4344b1990142afe190faf6dc30a9fa0f5e31"), // Burn
	common.HexToHash("0xfa8f14973a436f651cdc72fcb50527f364a3b92681dc7aacb0ebeed1e7fb7070"), // WipeBlacklistedAccount
	common.HexToHash("0x5ad932ee61740a6d8228a885eb9b654a6495a14372f49fed16686c4c80cc031a"), // WipeBlacklistedAccount
	common.HexToHash("0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6"), // DestroyedBlackFunds
	common.HexToHash("0x2ca1f69d129207da56149d0f73a0e0dce0d262588ea7da9e1523fc6fe04ad407"), // DestroyedBlackFunds
	//common.HexToHash("0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a"), // Issue
	//common.HexToHash("0xbc5846b4d860cbc33f6f6b6b2f7648fcbaca7425ffa5a8e26a9a70e5fd092f49"), // Issue
	//common.HexToHash("0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44"), // Redeem
	//common.HexToHash("0x13ffabc1e5c4de958be06f1541c4a24fdf8a4deffcbf851b2e805e41613e5e22"), // Redeem
	common.HexToHash("0x5f7542858008eeb041631f30e6109ae94b83a58e9a58261dd2c42c508850f939"), // TransferFrom
	common.HexToHash("0xbf480be9a10591c488f4bdf0305c23b7114779b0c1a42fb58401ef0a8574bc79"), // TransferFrom
	common.HexToHash("0x9b40ecb08ede8b5f14f7401641540771bbdaaac8638c21f5813ae7ba76d75155"), // Deposit
	common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c"), // Deposit
	common.HexToHash("0x1be94c6778a9e1751832385a51994fbb7b20c9c08ebfa22735a951d4b84ebb1e"), // Withdrawal
	common.HexToHash("0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65"), // Withdrawal
}

func NewLabelFromSignature(sig string) Label {
	for signature, label := range EventSignatures {
		if sig == signature {
			return label
		}
	}
	return ""
}
