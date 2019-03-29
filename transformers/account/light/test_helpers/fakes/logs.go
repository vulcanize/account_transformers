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

package fakes

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var FakeLogs = []types.Log{
	{ // This is a mock Mint to 0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4
		Address:     common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // The old TrueUSD
		BlockNumber: 6791666,
		BlockHash:   common.HexToHash("0xMockBlockHash00"),
		TxHash:      common.HexToHash("0xb3e3d6d098c5f8fc4f2f45d650943f18f8fff5367d6991a67dea2aafe8811e93"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
			common.HexToHash("0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{ // This is a mock Mint to 0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4
		Address:     common.HexToAddress("0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"), // The old TrueUSD
		BlockNumber: 6791666,
		BlockHash:   common.HexToHash("0xMockBlockHash00"),
		TxHash:      common.HexToHash("0xb3e3d6d098c5f8fc4f2f45d650943f18f8fff5367d6991a67dea2aafe8811e93"),
		TxIndex:     1,
		Index:       2,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
			common.HexToHash("0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{ // This is a mock Mint to 0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4
		Address:     common.HexToAddress("0x8dd5fbCe2F6a956C3022bA3663759011Dd51e73E"), // The old TrueUSD
		BlockNumber: 6791667,
		BlockHash:   common.HexToHash("0xMockBlockHash01"),
		TxHash:      common.HexToHash("0xb3e3d6d098c5f8fc4f2f45d650943f18f8fff5367d6991a67dea2aafe8811e93"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
			common.HexToHash("0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{ // This a mock Transfer from 0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4 to 0x000000000000000000000000009c1e8674038605c5ae33c74f13bc528e1222b5
		Address:     common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376"), // The new one
		BlockNumber: 6791668,
		BlockHash:   common.HexToHash("0xMockBlockHash02"),
		TxHash:      common.HexToHash("0x98634183774924055c36f18babbb895834dce31a9cb4d2397b639078ffd3605f"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			common.HexToHash("0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4"),
			common.HexToHash("0x000000000000000000000000009c1e8674038605c5ae33c74f13bc528e1222b5"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{ // This is a mock DestroyedBlackFunds of 0x000000000000000000000000009c1e8674038605c5ae33c74f13bc528e1222b5
		Address:     common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376"),
		BlockNumber: 6791669,
		BlockHash:   common.HexToHash("0xMockBlockHash03"),
		TxHash:      common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x61e6e66b0d6339b2980aecc6ccc0039736791f0ccde9ed512e789a7fbdd698c6"),
			common.HexToHash("0x000000000000000000000000009c1e8674038605c5ae33c74f13bc528e1222b5"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{ // This is a mock Mint to 0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4
		Address:     common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // The old TrueUSD
		BlockNumber: 6791669,
		BlockHash:   common.HexToHash("0xMockBlockHash03"),
		TxHash:      common.HexToHash("0xb3e3d6d098c5f8fc4f2f45d650943f18f8fff5367d6991a67dea2aafe8811e93"),
		TxIndex:     1,
		Index:       2,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"),
			common.HexToHash("0x00000000000000000000000048e78948c80e9f8f53190dbdf2990f9a69491ef4"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
}
