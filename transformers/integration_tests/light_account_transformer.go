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

package integration_tests

import (
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/fetcher"
	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/light/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"

	"github.com/vulcanize/account_transformers/transformers/account/config"
	transformer "github.com/vulcanize/account_transformers/transformers/account/light"
	"github.com/vulcanize/account_transformers/transformers/account/light/converters"
	r2 "github.com/vulcanize/account_transformers/transformers/account/light/repositories"
	"github.com/vulcanize/account_transformers/transformers/account/light/test_helpers/fakes"
	"github.com/vulcanize/account_transformers/transformers/account/light/test_helpers/mocks"
	"github.com/vulcanize/account_transformers/transformers/account/shared"
	c2 "github.com/vulcanize/account_transformers/transformers/account/shared/constants"
	"github.com/vulcanize/account_transformers/transformers/account/shared/poller"
	"github.com/vulcanize/account_transformers/transformers/account/shared/test_helpers"
)

var _ = Describe("Transformer", func() {
	var db *postgres.DB
	var blockChain core.BlockChain
	var headerRepository repositories.HeaderRepository
	var headerID, headerID2, headerID3 int64

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
			header, err := blockChain.GetHeaderByNumber(6791667)
			Expect(err).ToNot(HaveOccurred())
			headerID, err = headerRepository.CreateOrUpdateHeader(header)
			Expect(err).ToNot(HaveOccurred())
			header, err = blockChain.GetHeaderByNumber(6791668)
			Expect(err).ToNot(HaveOccurred())
			headerID2, err = headerRepository.CreateOrUpdateHeader(header)
			Expect(err).ToNot(HaveOccurred())
			header, err = blockChain.GetHeaderByNumber(6791669)
			Expect(err).ToNot(HaveOccurred())
			headerID3, err = headerRepository.CreateOrUpdateHeader(header)
			Expect(err).ToNot(HaveOccurred())
		})

		It("With Mock Fetcher: transforms value transfer events into account records", func() {
			vtc, err := converters.NewValueTransferConverter(c2.CombinedABI, c2.EquivalentTokenAddressesMapping())
			Expect(err).ToNot(HaveOccurred())
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       vtc,
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				AddressRepository:            r2.NewAddressRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			f := mocks.MockFetcher{}
			f.Logs = fakes.FakeLogs
			t.Fetcher = &f
			t.NextStart = 6791667
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
			addrs := []common.Address{
				common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4"),
				common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5"),
			}

			// We get token transfer records up to the provided block
			Expect(f.PassedHeaders[0].Id).To(Equal(headerID))
			transferRecords, err := t.ValueTransferEventRepository.GetTokenValueTransferRecordsForAccounts(addrs, 6791667)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(transferRecords[addrs[0]])).To(Equal(1))
			Expect(len(transferRecords[addrs[1]])).To(Equal(0))
			Expect(transferRecords[addrs[0]][0].BlockNumber).To(Equal(uint64(6791667)))
			Expect(transferRecords[addrs[0]][0].Contract).To(Equal("0x0000000000085d4780B73119b644AE5ecd22b376"))
			Expect(transferRecords[addrs[0]][0].HeaderID).To(Equal(headerID))
			Expect(transferRecords[addrs[0]][0].Src).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(transferRecords[addrs[0]][0].Dst).To(Equal("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4"))
			Expect(transferRecords[addrs[0]][0].Amount).To(Equal("376864137882094974530501285544524832305182681138"))

			Expect(f.PassedHeaders[1].Id).To(Equal(headerID2))
			transferRecords, err = t.ValueTransferEventRepository.GetTokenValueTransferRecordsForAccounts(addrs, 6791668)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(transferRecords[addrs[0]])).To(Equal(2))
			Expect(len(transferRecords[addrs[1]])).To(Equal(1))
			Expect(transferRecords[addrs[0]][1].BlockNumber).To(Equal(uint64(6791668)))
			Expect(transferRecords[addrs[0]][1].Contract).To(Equal("0x0000000000085d4780B73119b644AE5ecd22b376"))
			Expect(transferRecords[addrs[0]][1].HeaderID).To(Equal(headerID2))
			Expect(transferRecords[addrs[0]][1].Src).To(Equal("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4"))
			Expect(transferRecords[addrs[0]][1].Dst).To(Equal("0x009C1E8674038605C5AE33C74f13bC528E1222B5"))
			Expect(transferRecords[addrs[0]][1].Amount).To(Equal("376864137882094974530501285544524832305182681138"))
			Expect(transferRecords[addrs[1]][0].BlockNumber).To(Equal(uint64(6791668)))
			Expect(transferRecords[addrs[1]][0].Contract).To(Equal("0x0000000000085d4780B73119b644AE5ecd22b376"))
			Expect(transferRecords[addrs[1]][0].HeaderID).To(Equal(headerID2))
			Expect(transferRecords[addrs[1]][0].Src).To(Equal("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4"))
			Expect(transferRecords[addrs[1]][0].Dst).To(Equal("0x009C1E8674038605C5AE33C74f13bC528E1222B5"))
			Expect(transferRecords[addrs[1]][0].Amount).To(Equal("376864137882094974530501285544524832305182681138"))

			Expect(f.PassedHeaders[2].Id).To(Equal(headerID3))
			transferRecords, err = t.ValueTransferEventRepository.GetTokenValueTransferRecordsForAccounts(addrs, 6791669)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(transferRecords[addrs[0]])).To(Equal(2))
			Expect(len(transferRecords[addrs[1]])).To(Equal(2))
			Expect(transferRecords[addrs[1]][1].BlockNumber).To(Equal(uint64(6791669)))
			Expect(transferRecords[addrs[1]][1].Contract).To(Equal("0x0000000000085d4780B73119b644AE5ecd22b376"))
			Expect(transferRecords[addrs[1]][1].HeaderID).To(Equal(headerID3))
			Expect(transferRecords[addrs[1]][1].Src).To(Equal("0x009C1E8674038605C5AE33C74f13bC528E1222B5"))
			Expect(transferRecords[addrs[1]][1].Dst).To(Equal("0x0000000000000000000000000000000000000000"))
			Expect(transferRecords[addrs[1]][1].Amount).To(Equal("376864137882094974530501285544524832305182681138"))

			var coinRecord shared.CoinBalanceRecord
			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(), 6791667)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791667)))
			Expect(coinRecord.Value).To(Equal("18780247519"))

			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes(), 6791667)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791667)))
			Expect(coinRecord.Value).To(Equal("171056198103568077794"))

			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(), 6791668)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791668)))
			Expect(coinRecord.Value).To(Equal("165525871780247519"))

			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes(), 6791668)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791668)))
			Expect(coinRecord.Value).To(Equal("172845293271568077794"))

			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(), 6791669)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791669)))
			Expect(coinRecord.Value).To(Equal("165525871780247519"))

			err = db.Get(&coinRecord, `SELECT address_hash, block_number, value FROM accounts.address_coin_balances WHERE address_hash = $1 AND block_number = $2`, common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes(), 6791669)
			Expect(err).ToNot(HaveOccurred())
			Expect(coinRecord.Address).To(Equal(common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes()))
			Expect(coinRecord.BlockNumber).To(Equal(int64(6791669)))
			Expect(coinRecord.Value).To(Equal("172845293271568077794"))

			// First header, one event with only our first watched address => one record
			var tokenRecord shared.TokenBalanceRecord
			err = db.Get(&tokenRecord, `SELECT address_hash, block_number, value, token_contract_address_hash 
											FROM accounts.address_token_balances 
											WHERE address_hash = $1 AND token_contract_address_hash = $2 AND block_number = $3`,
				common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(),
				common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes(),
				6791667)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(tokenRecord.BlockNumber).To(Equal(int64(6791667)))
			Expect(tokenRecord.Value).To(Equal("376864137882094974530501285544524832305182681138"))
			Expect(tokenRecord.ContractAddress).To(Equal(common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes()))

			// Second header, another event with both our watched addresses => two more records
			err = db.Get(&tokenRecord, `SELECT address_hash, block_number, value, token_contract_address_hash 
											FROM accounts.address_token_balances 
											WHERE address_hash = $1 AND token_contract_address_hash = $2 AND block_number = $3`,
				common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(),
				common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes(),
				6791668)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(tokenRecord.BlockNumber).To(Equal(int64(6791668)))
			Expect(tokenRecord.Value).To(Equal("0"))
			Expect(tokenRecord.ContractAddress).To(Equal(common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes()))

			err = db.Get(&tokenRecord, `SELECT address_hash, block_number, value, token_contract_address_hash 
											FROM accounts.address_token_balances 
											WHERE address_hash = $1 AND token_contract_address_hash = $2 AND block_number = $3`,
				common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes(),
				common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes(),
				6791668)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenRecord.Address).To(Equal(common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes()))
			Expect(tokenRecord.BlockNumber).To(Equal(int64(6791668)))
			Expect(tokenRecord.Value).To(Equal("376864137882094974530501285544524832305182681138"))
			Expect(tokenRecord.ContractAddress).To(Equal(common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes()))

			// Third header, another event with only the 2nd address => two more records, one changed and one unchanged
			err = db.Get(&tokenRecord, `SELECT address_hash, block_number, value, token_contract_address_hash 
											FROM accounts.address_token_balances 
											WHERE address_hash = $1 AND token_contract_address_hash = $2 AND block_number = $3`,
				common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes(),
				common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes(),
				6791669)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenRecord.Address).To(Equal(common.HexToAddress("0x48E78948C80e9f8F53190DbDF2990f9a69491ef4").Bytes()))
			Expect(tokenRecord.BlockNumber).To(Equal(int64(6791669)))
			Expect(tokenRecord.Value).To(Equal("0"))
			Expect(tokenRecord.ContractAddress).To(Equal(common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes()))

			err = db.Get(&tokenRecord, `SELECT address_hash, block_number, value, token_contract_address_hash 
											FROM accounts.address_token_balances 
											WHERE address_hash = $1 AND token_contract_address_hash = $2 AND block_number = $3`,
				common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes(),
				common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes(),
				6791669)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenRecord.Address).To(Equal(common.HexToAddress("0x009C1E8674038605C5AE33C74f13bC528E1222B5").Bytes()))
			Expect(tokenRecord.BlockNumber).To(Equal(int64(6791669)))
			Expect(tokenRecord.Value).To(Equal("0"))
			Expect(tokenRecord.ContractAddress).To(Equal(common.HexToAddress("0x0000000000085d4780B73119b644AE5ecd22b376").Bytes()))
		})

		It("If `next start` isn't contiguous with the headers we have available, we can't do anything", func() {
			vtc, err := converters.NewValueTransferConverter(c2.CombinedABI, c2.EquivalentTokenAddressesMapping())
			Expect(err).ToNot(HaveOccurred())
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       vtc,
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				Fetcher:                      fetcher.NewFetcher(blockChain),
				AddressRepository:            r2.NewAddressRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

		It("With real fetcher: transforms value transfer events into account records", func() {
			vtc, err := converters.NewValueTransferConverter(c2.CombinedABI, c2.EquivalentTokenAddressesMapping())
			Expect(err).ToNot(HaveOccurred())
			t := transformer.TokenBalanceTransformer{
				ValueTransferConverter:       vtc,
				TokenBalanceConverter:        converters.NewTokenBalanceConverter(),
				HeaderRepository:             repository.NewHeaderRepository(db),
				Fetcher:                      fetcher.NewFetcher(blockChain),
				AddressRepository:            r2.NewAddressRepository(db),
				ValueTransferEventRepository: r2.NewValueTransferEventRepository(db),
				CoinBalanceRepository:        r2.NewAccountCoinBalanceRepository(db),
				TokenBalanceRepository:       r2.NewAccountTokenBalanceRepository(db),
				AccountPoller:                poller.NewAccountPoller(blockChain),
			}
			t.NextStart = 6791668
			err = t.Init()
			Expect(err).ToNot(HaveOccurred())
			err = t.Execute()
			Expect(err).ToNot(HaveOccurred())
		})

	})
})
