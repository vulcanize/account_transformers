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

package repositories

import (
	"database/sql"
	"errors"
	"math/big"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/libraries/shared/utilities"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

const (
	blocksFromHeadBeforeFinal = 20
)

var ErrBlockExists = errors.New("Won't add block that already exists.")

type BlockRepository struct {
	database *postgres.DB
}

func NewBlockRepository(database *postgres.DB) *BlockRepository {
	return &BlockRepository{database: database}
}

func (blockRepository BlockRepository) SetBlocksStatus(chainHead int64) error {
	cutoff := chainHead - blocksFromHeadBeforeFinal
	_, err := blockRepository.database.Exec(`
                  UPDATE blocks SET is_final = TRUE
                  WHERE is_final = FALSE AND number < $1`,
		cutoff)

	return err
}

func (blockRepository BlockRepository) CreateOrUpdateBlock(block core.Block) (int64, error) {
	var err error
	var blockId int64
	retrievedBlockHash, ok := blockRepository.getBlockHash(block)
	if !ok {
		return blockRepository.insertBlock(block)
	}
	if ok && retrievedBlockHash != block.Hash {
		err = blockRepository.removeBlock(block.Number)
		if err != nil {
			return 0, err
		}
		return blockRepository.insertBlock(block)
	}
	return blockId, ErrBlockExists
}

func (blockRepository BlockRepository) MissingBlockNumbers(startingBlockNumber int64, highestBlockNumber int64, nodeId string) []int64 {
	numbers := make([]int64, 0)
	err := blockRepository.database.Select(&numbers,
		`SELECT all_block_numbers
          FROM (
              SELECT generate_series($1::INT, $2::INT) AS all_block_numbers) series
          WHERE all_block_numbers NOT IN (
		  	  SELECT number FROM blocks WHERE eth_node_fingerprint = $3
		  ) `,
		startingBlockNumber,
		highestBlockNumber, nodeId)
	if err != nil {
		log.Error("MissingBlockNumbers: error getting blocks: ", err)
	}
	return numbers
}

func (blockRepository BlockRepository) GetBlock(blockNumber int64) (core.Block, error) {
	blockRows := blockRepository.database.QueryRowx(
		`SELECT id,
                       number,
                       gaslimit,
                       gasused,
                       time,
                       difficulty,
                       hash,
                       nonce,
                       parenthash,
                       size,
                       uncle_hash,
                       is_final,
                       miner,
                       extra_data,
                       reward,
                       uncles_reward
               FROM blocks
               WHERE eth_node_id = $1 AND number = $2`, blockRepository.database.NodeID, blockNumber)
	savedBlock, err := blockRepository.loadBlock(blockRows)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return core.Block{}, datastore.ErrBlockDoesNotExist(blockNumber)
		default:
			log.Error("GetBlock: error loading blocks: ", err)
			return savedBlock, err
		}
	}
	return savedBlock, nil
}

func (blockRepository BlockRepository) insertBlock(block core.Block) (int64, error) {
	var blockId int64
	tx, _ := blockRepository.database.Beginx()
	err := tx.QueryRow(
		`INSERT INTO blocks
                (eth_node_id, number, gaslimit, gasused, time, difficulty, hash, nonce, parenthash, size, uncle_hash, is_final, miner, extra_data, reward, uncles_reward, eth_node_fingerprint)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
                RETURNING id `,
		blockRepository.database.NodeID,
		block.Number,
		block.GasLimit,
		block.GasUsed,
		block.Time,
		block.Difficulty,
		block.Hash,
		block.Nonce,
		block.ParentHash,
		block.Size,
		block.UncleHash,
		block.IsFinal,
		block.Miner,
		block.ExtraData,
		utilities.NullToZero(block.Reward),
		utilities.NullToZero(block.UnclesReward),
		blockRepository.database.Node.ID).
		Scan(&blockId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(block.MappedUncleRewards) > 0 {
		err = blockRepository.createUncleRewards(tx, blockId, block.Hash, block.MappedUncleRewards)
		if err != nil {
			tx.Rollback()
			return 0, postgres.ErrDBInsertFailed
		}
	}
	if len(block.Transactions) > 0 {
		err = blockRepository.createTransactions(tx, blockId, block.Transactions)
		if err != nil {
			tx.Rollback()
			return 0, postgres.ErrDBInsertFailed
		}
	}
	tx.Commit()
	return blockId, nil
}

func (blockRepository BlockRepository) createUncleRewards(tx *sqlx.Tx, blockId int64, blockHash string, mappedUncleRewards map[string]map[string]*big.Int) error {
	for miner, uncleRewards := range mappedUncleRewards {
		for uncleHash, reward := range uncleRewards {
			err := blockRepository.createUncleReward(tx, blockId, blockHash, miner, uncleHash, reward.String())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (blockRepository BlockRepository) createUncleReward(tx *sqlx.Tx, blockId int64, blockHash, miner, uncleHash, amount string) error {
	_, err := tx.Exec(
		`INSERT INTO uncle_rewards
       (block_id, block_hash, uncle_hash, uncle_reward, miner_address)
       VALUES ($1, $2, $3, $4, $5)
       RETURNING id`,
		blockId, blockHash, miner, uncleHash, utilities.NullToZero(amount))
	return err
}

func (blockRepository BlockRepository) createTransactions(tx *sqlx.Tx, blockId int64, transactions []core.Transaction) error {
	for _, transaction := range transactions {
		err := blockRepository.createTransaction(tx, blockId, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

//Fields like value lose precision if converted to
//int64 so convert to string instead. But nil
//big.Int -> string = "" so convert to "0"
func nullStringToZero(s string) string {
	if s == "" {
		return "0"
	}
	return s
}

func (blockRepository BlockRepository) createTransaction(tx *sqlx.Tx, blockId int64, transaction core.Transaction) error {
	_, err := tx.Exec(
		`INSERT INTO transactions
       (block_id, hash, nonce, tx_to, tx_from, gaslimit, gasprice, value, input_data)
       VALUES ($1, $2, $3, $4, $5, $6, $7,  $8::NUMERIC, $9)
       RETURNING id`,
		blockId,
		transaction.Hash,
		transaction.Nonce,
		transaction.To,
		transaction.From,
		transaction.GasLimit,
		transaction.GasPrice,
		nullStringToZero(transaction.Value),
		transaction.Data)
	if err != nil {
		return err
	}
	if hasReceipt(transaction) {
		receiptId, err := blockRepository.createReceipt(tx, blockId, transaction.Receipt)
		if err != nil {
			return err
		}
		if hasLogs(transaction) {
			err = blockRepository.createLogs(tx, transaction.Receipt.Logs, receiptId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func hasLogs(transaction core.Transaction) bool {
	return len(transaction.Receipt.Logs) > 0
}

func hasReceipt(transaction core.Transaction) bool {
	return transaction.Receipt.TxHash != ""
}

func (blockRepository BlockRepository) createReceipt(tx *sqlx.Tx, blockId int64, receipt core.Receipt) (int, error) {
	//Not currently persisting log bloom filters
	var receiptId int
	err := tx.QueryRow(
		`INSERT INTO receipts
               (contract_address, tx_hash, cumulative_gas_used, gas_used, state_root, status, block_id)
               VALUES ($1, $2, $3, $4, $5, $6, $7) 
               RETURNING id`,
		receipt.ContractAddress, receipt.TxHash, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.StateRoot, receipt.Status, blockId).Scan(&receiptId)
	if err != nil {
		log.Error("createReceipt: error inserting receipt: ", err)
		return receiptId, err
	}
	return receiptId, nil
}

func (blockRepository BlockRepository) getBlockHash(block core.Block) (string, bool) {
	var retrievedBlockHash string
	blockRepository.database.Get(&retrievedBlockHash,
		`SELECT hash
               FROM blocks
               WHERE number = $1 AND eth_node_id = $2`,
		block.Number, blockRepository.database.NodeID)
	return retrievedBlockHash, blockExists(retrievedBlockHash)
}

func (blockRepository BlockRepository) createLogs(tx *sqlx.Tx, logs []core.Log, receiptId int) error {
	for _, tlog := range logs {
		_, err := tx.Exec(
			`INSERT INTO logs (block_number, address, tx_hash, index, topic0, topic1, topic2, topic3, data, receipt_id)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                `,
			tlog.BlockNumber, tlog.Address, tlog.TxHash, tlog.Index, tlog.Topics[0], tlog.Topics[1], tlog.Topics[2], tlog.Topics[3], tlog.Data, receiptId,
		)
		if err != nil {
			return postgres.ErrDBInsertFailed
		}
	}
	return nil
}

func blockExists(retrievedBlockHash string) bool {
	return retrievedBlockHash != ""
}

func (blockRepository BlockRepository) removeBlock(blockNumber int64) error {
	_, err := blockRepository.database.Exec(
		`DELETE FROM
                blocks
                WHERE number=$1 AND eth_node_id=$2`,
		blockNumber, blockRepository.database.NodeID)
	if err != nil {
		return postgres.ErrDBDeleteFailed
	}
	return nil
}

func (blockRepository BlockRepository) loadBlock(blockRows *sqlx.Row) (core.Block, error) {
	type b struct {
		ID int
		core.Block
	}
	var block b
	err := blockRows.StructScan(&block)
	if err != nil {
		log.Error("loadBlock: error loading block: ", err)
		return core.Block{}, err
	}
	transactionRows, err := blockRepository.database.Queryx(`
            SELECT hash,
				   nonce,
				   tx_to,
				   tx_from,
				   gaslimit,
				   gasprice,
				   value,
				   input_data
            FROM transactions
            WHERE block_id = $1
            ORDER BY hash`, block.ID)
	if err != nil {
		log.Error("loadBlock: error fetting transactions: ", err)
		return core.Block{}, err
	}
	block.Transactions = blockRepository.LoadTransactions(transactionRows)
	return block.Block, nil
}

func (blockRepository BlockRepository) LoadTransactions(transactionRows *sqlx.Rows) []core.Transaction {
	var transactions []core.Transaction
	for transactionRows.Next() {
		var transaction core.Transaction
		err := transactionRows.StructScan(&transaction)
		if err != nil {
			log.Fatal(err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
