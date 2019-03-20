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

import "github.com/ethereum/go-ethereum/common"

var EventSignatures = map[Label]string{
	Transfer : "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
	TransferDev : "0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377",
	Mint : "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885",
	MintDev : "0xb4683f1a6bfdb19078899ae88702d383183e50a367af29a37f4e7357a7fc13f9",
	WipedAccount : "0xdf58d2368c06216a398f05a7a88c8edc64a25c33f33fd2bd8b56fbc8822c02d8",
	WipedAccountDev : "0x6dd8c32d75aebbe065ce6cbbce463910d4d71fcfa91bb7fb276358b35223b8c1",
	Burn : "0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5",
	BurnDev : "0x654fd2845115fd0c81b8dd4d62aa4344b1990142afe190faf6dc30a9fa0f5e31",
	WipeBlacklistedAccount : "0xfa8f14973a436f651cdc72fcb50527f364a3b92681dc7aacb0ebeed1e7fb7070",
	WipeBlacklistedAccountDev : "0x5ad932ee61740a6d8228a885eb9b654a6495a14372f49fed16686c4c80cc031a",
	DestroyedBlackFunds : "0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6",
	DestroyedBlackFundsDev : "0x2ca1f69d129207da56149d0f73a0e0dce0d262588ea7da9e1523fc6fe04ad407",
	Issue : "0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a",
	IssueDev : "0xbc5846b4d860cbc33f6f6b6b2f7648fcbaca7425ffa5a8e26a9a70e5fd092f49",
	Redeem : "0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44",
	RedeemDev : "0x13ffabc1e5c4de958be06f1541c4a24fdf8a4deffcbf851b2e805e41613e5e22",
	TransferFrom : "0x5f7542858008eeb041631f30e6109ae94b83a58e9a58261dd2c42c508850f939",
	TransferFromDev : "0xbf480be9a10591c488f4bdf0305c23b7114779b0c1a42fb58401ef0a8574bc79",
}

var Topic0s = []common.Hash{
	common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
	common.HexToHash("0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377"),
	common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
	common.HexToHash("0xb4683f1a6bfdb19078899ae88702d383183e50a367af29a37f4e7357a7fc13f9"),
	common.HexToHash("0xdf58d2368c06216a398f05a7a88c8edc64a25c33f33fd2bd8b56fbc8822c02d8"),
	common.HexToHash("0x6dd8c32d75aebbe065ce6cbbce463910d4d71fcfa91bb7fb276358b35223b8c1"),
	common.HexToHash("0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5"),
	common.HexToHash("0x654fd2845115fd0c81b8dd4d62aa4344b1990142afe190faf6dc30a9fa0f5e31"),
	common.HexToHash("0xfa8f14973a436f651cdc72fcb50527f364a3b92681dc7aacb0ebeed1e7fb7070"),
	common.HexToHash("0x5ad932ee61740a6d8228a885eb9b654a6495a14372f49fed16686c4c80cc031a"),
	common.HexToHash("0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6"),
	common.HexToHash("0x2ca1f69d129207da56149d0f73a0e0dce0d262588ea7da9e1523fc6fe04ad407"),
	common.HexToHash("0xcb8241adb0c3fdb35b70c24ce35c5eb0c17af7431c99f827d44a445ca624176a"),
	common.HexToHash("0xbc5846b4d860cbc33f6f6b6b2f7648fcbaca7425ffa5a8e26a9a70e5fd092f49"),
	common.HexToHash("0x702d5967f45f6513a38ffc42d6ba9bf230bd40e8f53b16363c7eb4fd2deb9a44"),
	common.HexToHash("0x13ffabc1e5c4de958be06f1541c4a24fdf8a4deffcbf851b2e805e41613e5e22"),
	common.HexToHash("0x5f7542858008eeb041631f30e6109ae94b83a58e9a58261dd2c42c508850f939"),
	common.HexToHash("0xbf480be9a10591c488f4bdf0305c23b7114779b0c1a42fb58401ef0a8574bc79"),
}

func (label Label) Sig() string {
	return EventSignatures[label]
}

func NewLabelFromSignature(sig string) Label {
	for label, signature := range EventSignatures {
		if sig == signature {
			return label
		}
	}
	return ""
}