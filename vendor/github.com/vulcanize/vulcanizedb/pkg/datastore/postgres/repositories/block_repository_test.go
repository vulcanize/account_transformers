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

package repositories_test

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/big"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Saving blocks", func() {
	var db *postgres.DB
	var node core.Node
	var blockRepository datastore.BlockRepository

	BeforeEach(func() {
		node = core.Node{
			GenesisBlock: "GENESIS",
			NetworkID:    1,
			ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
			ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
		}
		db = test_config.NewTestDB(node)
		test_config.CleanTestDB(db)
		blockRepository = repositories.NewBlockRepository(db)

	})

	It("associates blocks to a node", func() {
		block := core.Block{
			Number: 123,
		}
		_, insertErr := blockRepository.CreateOrUpdateBlock(block)
		Expect(insertErr).NotTo(HaveOccurred())
		nodeTwo := core.Node{
			GenesisBlock: "0x456",
			NetworkID:    1,
			ID:           "x123456",
			ClientName:   "Geth",
		}
		dbTwo := test_config.NewTestDB(nodeTwo)
		test_config.CleanTestDB(dbTwo)
		repositoryTwo := repositories.NewBlockRepository(dbTwo)

		_, err := repositoryTwo.GetBlock(123)
		Expect(err).To(HaveOccurred())
	})

	It("saves the attributes of the block", func() {
		blockNumber := int64(123)
		gasLimit := uint64(1000000)
		gasUsed := uint64(10)
		blockHash := "x123"
		blockParentHash := "x456"
		blockNonce := "0x881db2ca900682e9a9"
		miner := "x123"
		extraData := "xextraData"
		blockTime := int64(1508981640)
		uncleHash := "x789"
		blockSize := string("1000")
		difficulty := int64(10)
		blockReward := "5132000000000000000"
		unclesReward := "3580000000000000000"
		block := core.Block{
			Reward:       blockReward,
			Difficulty:   difficulty,
			GasLimit:     gasLimit,
			GasUsed:      gasUsed,
			Hash:         blockHash,
			ExtraData:    extraData,
			Nonce:        blockNonce,
			Miner:        miner,
			Number:       blockNumber,
			ParentHash:   blockParentHash,
			Size:         blockSize,
			Time:         blockTime,
			UncleHash:    uncleHash,
			UnclesReward: unclesReward,
		}

		_, insertErr := blockRepository.CreateOrUpdateBlock(block)

		Expect(insertErr).NotTo(HaveOccurred())
		savedBlock, err := blockRepository.GetBlock(blockNumber)
		Expect(err).NotTo(HaveOccurred())
		Expect(savedBlock.Reward).To(Equal(blockReward))
		Expect(savedBlock.Difficulty).To(Equal(difficulty))
		Expect(savedBlock.GasLimit).To(Equal(gasLimit))
		Expect(savedBlock.GasUsed).To(Equal(gasUsed))
		Expect(savedBlock.Hash).To(Equal(blockHash))
		Expect(savedBlock.Nonce).To(Equal(blockNonce))
		Expect(savedBlock.Miner).To(Equal(miner))
		Expect(savedBlock.ExtraData).To(Equal(extraData))
		Expect(savedBlock.Number).To(Equal(blockNumber))
		Expect(savedBlock.ParentHash).To(Equal(blockParentHash))
		Expect(savedBlock.Size).To(Equal(blockSize))
		Expect(savedBlock.Time).To(Equal(blockTime))
		Expect(savedBlock.UncleHash).To(Equal(uncleHash))
		Expect(savedBlock.UnclesReward).To(Equal(unclesReward))
	})

	It("does not find a block when searching for a number that does not exist", func() {
		_, err := blockRepository.GetBlock(111)

		Expect(err).To(HaveOccurred())
	})

	It("saves one transaction associated to the block", func() {
		block := core.Block{
			Number:       123,
			Transactions: []core.TransactionModel{fakes.FakeTransaction},
		}

		_, insertErr := blockRepository.CreateOrUpdateBlock(block)

		Expect(insertErr).NotTo(HaveOccurred())
		savedBlock, getErr := blockRepository.GetBlock(123)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(len(savedBlock.Transactions)).To(Equal(1))
	})

	It("saves two transactions associated to the block", func() {
		block := core.Block{
			Number:       123,
			Transactions: []core.TransactionModel{fakes.FakeTransaction, fakes.FakeTransaction},
		}

		_, insertErr := blockRepository.CreateOrUpdateBlock(block)

		Expect(insertErr).NotTo(HaveOccurred())
		savedBlock, getErr := blockRepository.GetBlock(123)
		Expect(getErr).NotTo(HaveOccurred())
		Expect(len(savedBlock.Transactions)).To(Equal(2))
	})

	It(`replaces blocks and transactions associated to the block
			when a more new block is in conflict (same block number + nodeid)`, func() {
		blockOne := core.Block{
			Number: 123,
			Hash:   "xabc",
			Transactions: []core.TransactionModel{
				fakes.GetFakeTransaction("x123", core.Receipt{}),
				fakes.GetFakeTransaction("x345", core.Receipt{}),
			},
		}
		blockTwo := core.Block{
			Number: 123,
			Hash:   "xdef",
			Transactions: []core.TransactionModel{
				fakes.GetFakeTransaction("x678", core.Receipt{}),
				fakes.GetFakeTransaction("x9ab", core.Receipt{}),
			},
		}

		_, insertErrOne := blockRepository.CreateOrUpdateBlock(blockOne)
		Expect(insertErrOne).NotTo(HaveOccurred())
		_, insertErrTwo := blockRepository.CreateOrUpdateBlock(blockTwo)
		Expect(insertErrTwo).NotTo(HaveOccurred())

		savedBlock, _ := blockRepository.GetBlock(123)
		Expect(len(savedBlock.Transactions)).To(Equal(2))
		Expect(savedBlock.Transactions[0].Hash).To(Equal("x678"))
		Expect(savedBlock.Transactions[1].Hash).To(Equal("x9ab"))
	})

	It(`does not replace blocks when block number is not unique
			     but block number + node id is`, func() {
		blockOne := core.Block{
			Number: 123,
			Transactions: []core.TransactionModel{
				fakes.GetFakeTransaction("x123", core.Receipt{}),
				fakes.GetFakeTransaction("x345", core.Receipt{}),
			},
		}
		blockTwo := core.Block{
			Number: 123,
			Transactions: []core.TransactionModel{
				fakes.GetFakeTransaction("x678", core.Receipt{}),
				fakes.GetFakeTransaction("x9ab", core.Receipt{}),
			},
		}
		_, insertErrOne := blockRepository.CreateOrUpdateBlock(blockOne)
		Expect(insertErrOne).NotTo(HaveOccurred())
		nodeTwo := core.Node{
			GenesisBlock: "0x456",
			NetworkID:    1,
		}
		dbTwo := test_config.NewTestDB(nodeTwo)
		test_config.CleanTestDB(dbTwo)
		repositoryTwo := repositories.NewBlockRepository(dbTwo)

		_, insertErrTwo := blockRepository.CreateOrUpdateBlock(blockOne)
		Expect(insertErrTwo).NotTo(HaveOccurred())
		_, insertErrThree := repositoryTwo.CreateOrUpdateBlock(blockTwo)
		Expect(insertErrThree).NotTo(HaveOccurred())
		retrievedBlockOne, getErrOne := blockRepository.GetBlock(123)
		Expect(getErrOne).NotTo(HaveOccurred())
		retrievedBlockTwo, getErrTwo := repositoryTwo.GetBlock(123)
		Expect(getErrTwo).NotTo(HaveOccurred())

		Expect(retrievedBlockOne.Transactions[0].Hash).To(Equal("x123"))
		Expect(retrievedBlockTwo.Transactions[0].Hash).To(Equal("x678"))
	})

	It("returns 'block exists' error if attempting to add duplicate block", func() {
		block := core.Block{
			Number: 12345,
			Hash:   "0x12345",
		}

		_, err := blockRepository.CreateOrUpdateBlock(block)

		Expect(err).NotTo(HaveOccurred())

		_, err = blockRepository.CreateOrUpdateBlock(block)

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(repositories.ErrBlockExists))
	})

	It("saves the attributes associated to a transaction", func() {
		gasLimit := uint64(5000)
		gasPrice := int64(3)
		nonce := uint64(10000)
		to := "1234567890"
		from := "0987654321"
		var value = new(big.Int)
		value.SetString("34940183920000000000", 10)
		inputData := "0xf7d8c8830000000000000000000000000000000000000000000000000000000000037788000000000000000000000000000000000000000000000000000000000003bd14"
		gethTransaction := types.NewTransaction(nonce, common.HexToAddress(to), value, gasLimit, big.NewInt(gasPrice), common.FromHex(inputData))
		var raw bytes.Buffer
		rlpErr := gethTransaction.EncodeRLP(&raw)
		Expect(rlpErr).NotTo(HaveOccurred())
		transaction := core.TransactionModel{
			Data:     common.Hex2Bytes(inputData),
			From:     from,
			GasLimit: gasLimit,
			GasPrice: gasPrice,
			Hash:     "x1234",
			Nonce:    nonce,
			Raw:      raw.Bytes(),
			Receipt:  core.Receipt{},
			To:       to,
			TxIndex:  2,
			Value:    value.String(),
		}
		block := core.Block{
			Number:       123,
			Transactions: []core.TransactionModel{transaction},
		}

		_, insertErr := blockRepository.CreateOrUpdateBlock(block)
		Expect(insertErr).NotTo(HaveOccurred())

		savedBlock, err := blockRepository.GetBlock(123)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(savedBlock.Transactions)).To(Equal(1))
		savedTransaction := savedBlock.Transactions[0]
		Expect(savedTransaction).To(Equal(transaction))
	})

	Describe("The missing block numbers", func() {
		It("is empty the starting block number is the highest known block number", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(core.Block{Number: 1})

			Expect(insertErr).NotTo(HaveOccurred())
			Expect(len(blockRepository.MissingBlockNumbers(1, 1, node.ID))).To(Equal(0))
		})

		It("is empty if copies of block exist from both current node and another", func() {
			_, insertErrOne := blockRepository.CreateOrUpdateBlock(core.Block{Number: 0})
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(core.Block{Number: 1})
			Expect(insertErrTwo).NotTo(HaveOccurred())
			nodeTwo := core.Node{
				GenesisBlock: "0x456",
				NetworkID:    1,
			}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())
			repositoryTwo := repositories.NewBlockRepository(dbTwo)
			_, insertErrThree := repositoryTwo.CreateOrUpdateBlock(core.Block{Number: 0})
			Expect(insertErrThree).NotTo(HaveOccurred())

			missing := blockRepository.MissingBlockNumbers(0, 1, node.ID)

			Expect(len(missing)).To(BeZero())
		})

		It("is the only missing block number", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(core.Block{Number: 2})
			Expect(insertErr).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(1, 2, node.ID)).To(Equal([]int64{1}))
		})

		It("is both missing block numbers", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(core.Block{Number: 3})
			Expect(insertErr).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(1, 3, node.ID)).To(Equal([]int64{1, 2}))
		})

		It("goes back to the starting block number", func() {
			_, insertErr := blockRepository.CreateOrUpdateBlock(core.Block{Number: 6})
			Expect(insertErr).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(4, 6, node.ID)).To(Equal([]int64{4, 5}))
		})

		It("only includes missing block numbers", func() {
			_, insertErrOne := blockRepository.CreateOrUpdateBlock(core.Block{Number: 4})
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(core.Block{Number: 6})
			Expect(insertErrTwo).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(4, 6, node.ID)).To(Equal([]int64{5}))
		})

		It("includes blocks created by a different node", func() {
			_, insertErrOne := blockRepository.CreateOrUpdateBlock(core.Block{Number: 4})
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(core.Block{Number: 6})
			Expect(insertErrTwo).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(4, 6, "Different node id")).To(Equal([]int64{4, 5, 6}))
		})

		It("is a list with multiple gaps", func() {
			_, insertErrOne := blockRepository.CreateOrUpdateBlock(core.Block{Number: 4})
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(core.Block{Number: 5})
			Expect(insertErrTwo).NotTo(HaveOccurred())
			_, insertErrThree := blockRepository.CreateOrUpdateBlock(core.Block{Number: 8})
			Expect(insertErrThree).NotTo(HaveOccurred())
			_, insertErrFour := blockRepository.CreateOrUpdateBlock(core.Block{Number: 10})
			Expect(insertErrFour).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(3, 10, node.ID)).To(Equal([]int64{3, 6, 7, 9}))
		})

		It("returns empty array when lower bound exceeds upper bound", func() {
			Expect(blockRepository.MissingBlockNumbers(10000, 1, node.ID)).To(Equal([]int64{}))
		})

		It("only returns requested range even when other gaps exist", func() {
			_, insertErrOne := blockRepository.CreateOrUpdateBlock(core.Block{Number: 3})
			Expect(insertErrOne).NotTo(HaveOccurred())
			_, insertErrTwo := blockRepository.CreateOrUpdateBlock(core.Block{Number: 8})
			Expect(insertErrTwo).NotTo(HaveOccurred())

			Expect(blockRepository.MissingBlockNumbers(1, 5, node.ID)).To(Equal([]int64{1, 2, 4, 5}))
		})
	})

	Describe("The block status", func() {
		It("sets the status of blocks within n-20 of chain HEAD as final", func() {
			blockNumberOfChainHead := 25
			for i := 0; i < blockNumberOfChainHead; i++ {
				_, err := blockRepository.CreateOrUpdateBlock(core.Block{Number: int64(i), Hash: strconv.Itoa(i)})
				Expect(err).NotTo(HaveOccurred())
			}

			setErr := blockRepository.SetBlocksStatus(int64(blockNumberOfChainHead))

			Expect(setErr).NotTo(HaveOccurred())
			blockOne, err := blockRepository.GetBlock(1)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockOne.IsFinal).To(Equal(true))
			blockTwo, err := blockRepository.GetBlock(24)
			Expect(err).ToNot(HaveOccurred())
			Expect(blockTwo.IsFinal).To(BeFalse())
		})
	})
})
