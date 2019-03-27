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
	"database/sql"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/test_config"
)

var _ = Describe("Block header repository", func() {
	var (
		rawHeader []byte
		err       error
		timestamp string
		db        *postgres.DB
		repo      repositories.HeaderRepository
		header    core.Header
	)

	BeforeEach(func() {
		rawHeader, err = json.Marshal(types.Header{})
		Expect(err).NotTo(HaveOccurred())
		timestamp = big.NewInt(123456789).String()

		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repo = repositories.NewHeaderRepository(db)
		header = core.Header{
			BlockNumber: 100,
			Hash:        common.BytesToHash([]byte{1, 2, 3, 4, 5}).Hex(),
			Raw:         rawHeader,
			Timestamp:   timestamp,
		}
	})

	Describe("creating or updating a header", func() {
		It("adds a header", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			var dbHeader core.Header
			err = db.Get(&dbHeader, `SELECT block_number, hash, raw, block_timestamp FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("adds node data to header", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			var ethNodeId int64
			err = db.Get(&ethNodeId, `SELECT eth_node_id FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(ethNodeId).To(Equal(db.NodeID))
			var ethNodeFingerprint string
			err = db.Get(&ethNodeFingerprint, `SELECT eth_node_fingerprint FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(ethNodeFingerprint).To(Equal(db.Node.ID))
		})

		It("returns valid header exists error if attempting duplicate headers", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(repositories.ErrValidHeaderExists))

			var dbHeaders []core.Header
			err = db.Select(&dbHeaders, `SELECT block_number, hash, raw FROM public.headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(1))
		})

		It("replaces header if hash is different", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repo.CreateOrUpdateHeader(headerTwo)

			Expect(err).NotTo(HaveOccurred())
			var dbHeader core.Header
			err = db.Get(&dbHeader, `SELECT block_number, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.Hash).To(Equal(headerTwo.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(headerTwo.Raw))
		})

		It("does not replace header if node fingerprint is different", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())

			repoTwo := repositories.NewHeaderRepository(dbTwo)
			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repoTwo.CreateOrUpdateHeader(headerTwo)

			Expect(err).NotTo(HaveOccurred())
			var dbHeaders []core.Header
			err = dbTwo.Select(&dbHeaders, `SELECT block_number, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(2))
		})

		It("only replaces header with matching node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())

			repoTwo := repositories.NewHeaderRepository(dbTwo)
			headerTwo := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{5, 4, 3, 2, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}
			_, err = repoTwo.CreateOrUpdateHeader(headerTwo)
			headerThree := core.Header{
				BlockNumber: header.BlockNumber,
				Hash:        common.BytesToHash([]byte{1, 1, 1, 1, 1}).Hex(),
				Raw:         rawHeader,
				Timestamp:   timestamp,
			}

			_, err = repoTwo.CreateOrUpdateHeader(headerThree)

			Expect(err).NotTo(HaveOccurred())
			var dbHeaders []core.Header
			err = dbTwo.Select(&dbHeaders, `SELECT block_number, hash, raw FROM headers WHERE block_number = $1`, header.BlockNumber)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbHeaders)).To(Equal(2))
			Expect(dbHeaders[0].Hash).To(Or(Equal(header.Hash), Equal(headerThree.Hash)))
			Expect(dbHeaders[1].Hash).To(Or(Equal(header.Hash), Equal(headerThree.Hash)))
			Expect(dbHeaders[0].Raw).To(Or(MatchJSON(header.Raw), MatchJSON(headerThree.Raw)))
			Expect(dbHeaders[1].Raw).To(Or(MatchJSON(header.Raw), MatchJSON(headerThree.Raw)))
		})
	})

	Describe("creating a transaction", func() {
		It("adds a transaction", func() {
			headerID, err := repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txIndex := big.NewInt(123)
			transaction := core.TransactionModel{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
			}

			insertErr := repo.CreateTransaction(headerID, transaction)

			Expect(insertErr).NotTo(HaveOccurred())
			var dbTransaction core.TransactionModel
			err = db.Get(&dbTransaction,
				`SELECT hash, gaslimit, gasprice, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.light_sync_transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbTransaction).To(Equal(transaction))
		})

		It("silently ignores duplicate inserts", func() {
			headerID, err := repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())
			fromAddress := common.HexToAddress("0x1234")
			toAddress := common.HexToAddress("0x5678")
			txHash := common.HexToHash("0x9876")
			txIndex := big.NewInt(123)
			transaction := core.TransactionModel{
				Data:     []byte{},
				From:     fromAddress.Hex(),
				GasLimit: 0,
				GasPrice: 0,
				Hash:     txHash.Hex(),
				Nonce:    0,
				Raw:      []byte{},
				Receipt:  core.Receipt{},
				To:       toAddress.Hex(),
				TxIndex:  txIndex.Int64(),
				Value:    "0",
			}

			insertErr := repo.CreateTransaction(headerID, transaction)
			Expect(insertErr).NotTo(HaveOccurred())

			insertTwoErr := repo.CreateTransaction(headerID, transaction)
			Expect(insertTwoErr).NotTo(HaveOccurred())

			var dbTransactions []core.TransactionModel
			err = db.Select(&dbTransactions,
				`SELECT hash, gaslimit, gasprice, input_data, nonce, raw, tx_from, tx_index, tx_to, "value"
				FROM public.light_sync_transactions WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(dbTransactions)).To(Equal(1))
		})
	})

	Describe("Getting a header", func() {
		It("returns header if it exists", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			dbHeader, err := repo.GetHeader(header.BlockNumber)

			Expect(err).NotTo(HaveOccurred())
			Expect(dbHeader.Id).NotTo(BeZero())
			Expect(dbHeader.BlockNumber).To(Equal(header.BlockNumber))
			Expect(dbHeader.Hash).To(Equal(header.Hash))
			Expect(dbHeader.Raw).To(MatchJSON(header.Raw))
			Expect(dbHeader.Timestamp).To(Equal(header.Timestamp))
		})

		It("does not return header for a different node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(header)
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			_, err = repoTwo.GetHeader(header.BlockNumber)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})
	})

	Describe("Getting missing headers", func() {
		It("returns block numbers for headers not in the database", func() {
			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 1,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 3,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 5,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			missingBlockNumbers, err := repo.MissingBlockNumbers(1, 5, db.Node.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{2, 4}))
		})

		It("does not count headers created by a different node fingerprint", func() {
			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 1,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 3,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.CreateOrUpdateHeader(core.Header{
				BlockNumber: 5,
				Raw:         rawHeader,
				Timestamp:   timestamp,
			})
			Expect(err).NotTo(HaveOccurred())

			nodeTwo := core.Node{ID: "FingerprintTwo"}
			dbTwo, err := postgres.NewDB(test_config.DBConfig, nodeTwo)
			Expect(err).NotTo(HaveOccurred())
			repoTwo := repositories.NewHeaderRepository(dbTwo)

			missingBlockNumbers, err := repoTwo.MissingBlockNumbers(1, 5, nodeTwo.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(missingBlockNumbers).To(ConsistOf([]int64{1, 2, 3, 4, 5}))
		})
	})
})
