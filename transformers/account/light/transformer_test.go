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

package light_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/constants"
	mocks2 "github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers/test_helpers/mocks"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"

	"github.com/vulcanize/account_transformers/transformers/account/config"
	transformer "github.com/vulcanize/account_transformers/transformers/account/light"
	"github.com/vulcanize/account_transformers/transformers/account/light/converters"
	r2 "github.com/vulcanize/account_transformers/transformers/account/light/repositories"
	c2 "github.com/vulcanize/account_transformers/transformers/account/shared/constants"
	"github.com/vulcanize/account_transformers/transformers/account/shared/poller"
	"github.com/vulcanize/account_transformers/transformers/account/shared/test_data/mocks"
	"github.com/vulcanize/account_transformers/transformers/account/shared/test_helpers"
)

var mockLogs = []types.Log{
	{
		Address:     common.HexToAddress(constants.EnsContractAddress),
		BlockNumber: 6885692,
		BlockHash:   common.HexToHash("0xMockBlockHash01"),
		TxHash:      common.HexToHash("0xMockTxHash01"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			common.HexToHash("0xd1115c02622703bb9236a0e6609cb250a874e903494bd9071c25078f4033dac1"),
			common.HexToHash("0xadc756803e4eb4ccfb136b73d5f72e3dc0d452d30ae1f4bc82af394c73ce7115"),
		},
		Data: common.HexToHash("0x00000000000000000000000042032c22c510ad0698f16be9b99640efdeb02832").Bytes(),
	},
	{
		Address:     common.HexToAddress(constants.EnsContractAddress),
		BlockNumber: 6885693,
		BlockHash:   common.HexToHash("0xMockBlockHash02"),
		TxHash:      common.HexToHash("0xMockTxHash02"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0x930a61a57a70a73c2a503615b87e2e54fe5b9cdeacda518270b852296ab1a377"),
			common.HexToHash("0x5954c882606735d75f2775ff380873d6d6b546f63cdf79424f12209b9e15bb91"),
			common.HexToHash("0xadc756803e4eb4ccfb136b73d5f72e3dc0d452d30ae1f4bc82af394c73ce7115"),
		},
		Data: common.HexToHash("0x000000000000000000000000d3ddccdd3b25a8a7423b5bee360a42146eb4baf3").Bytes(),
	},
	{
		Address:     common.HexToAddress("0x000000000000000000000000d3ddccdd3b25a8a7423b5bee360a42146eb4baf3"),
		BlockNumber: 6885694,
		BlockHash:   common.HexToHash("0xMockBlockHash03"),
		TxHash:      common.HexToHash("0xMockTxHash03"),
		TxIndex:     1,
		Index:       1,
		Removed:     false,
		Topics: []common.Hash{
			common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			common.HexToHash("0x5954c882606735d75f2775ff380873d6d6b546f63cdf79424f12209b9e15bb91"),
			common.HexToHash("0xd1115c02622703bb9236a0e6609cb250a874e903494bd9071c25078f4033dac1"),
		},
		Data: common.HexToHash("0x000000000000000000000000a54aef7fa503e75a03b262a4cd73037c1774735d").Bytes(),
	},
}

var _ = Describe("Transformer", func() {
	var db *postgres.DB
	var blockChain core.BlockChain
	var headerRepository repositories.HeaderRepository

	BeforeEach(func() {
		db, blockChain = test_helpers.SetupDBandBC()
		headerRepository = repositories.NewHeaderRepository(db)
		test_helpers.TearDown(db)
	})

	AfterEach(func() {
		test_helpers.TearDown(db)
	})

	Describe("Init", func() {
		It("Doesn't do anything; fills in interface", func() {
			ti := transformer.TokenBalanceTransformer{
				Config: config.MainnetAccountConfig,
			}
			t := ti.NewTransformer(db, blockChain)

			err := t.Init()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Execute", func() {
		BeforeEach(func() { // 6885692 to 6885694
			for i := 6791668; i <= 6791670; i++ {
				header := mocks2.MockHeader1
				header, err := blockChain.GetHeaderByNumber(int64(i))
				Expect(err).ToNot(HaveOccurred())
				_, err = headerRepository.CreateOrUpdateHeader(header)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("With Mock Fetcher: transforms value transfer events into account records", func() {
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       converters.NewValueTransferConverter(c2.EquivalentTokenAddressesMapping()),
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				AddressRepository:            r2.NewAccountHeaderRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			f := mocks.NewMockFetcher(blockChain)
			f.Logs = mockLogs
			t.Fetcher = f
			t.NextStart = 6791668
			t.AddressRepository.AddAddress(common.HexToAddress("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef").String())
			err := t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("If `next start` isn't contiguous with the headers we have available, we can't do anything", func() {
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       converters.NewValueTransferConverter(c2.EquivalentTokenAddressesMapping()),
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				Fetcher:                      fetcher.NewFetcher(blockChain),
				AddressRepository:            r2.NewAccountHeaderRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			t.AddressRepository.AddAddress(common.HexToAddress("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef").String())
			err := t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("Throws an error if there are no account addresses to build records for", func() {
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       converters.NewValueTransferConverter(c2.EquivalentTokenAddressesMapping()),
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				Fetcher:                      fetcher.NewFetcher(blockChain),
				AddressRepository:            r2.NewAccountHeaderRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			t.NextStart = 6791668
			err := t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).To(HaveOccurred())
		})

		It("With real fetcher: transforms value transfer events into account records", func() {
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       converters.NewValueTransferConverter(c2.EquivalentTokenAddressesMapping()),
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				Fetcher:                      fetcher.NewFetcher(blockChain),
				AddressRepository:            r2.NewAccountHeaderRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			t.NextStart = 6791668
			t.AddressRepository.AddAddress(common.HexToAddress("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef").String())
			err := t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
