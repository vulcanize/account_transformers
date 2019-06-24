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

package seed

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
)

type Repository interface {
	CommitPayload(payload DecodedPayload) (int64, error)
}

type BlockExplorerRepository struct {
	db               *postgres.DB
	headerRepository repositories.HeaderRepository
}

func NewBlockExplorerRepository(db *postgres.DB) *BlockExplorerRepository {
	return &BlockExplorerRepository{
		db:               db,
		headerRepository: repositories.NewHeaderRepository(db),
	}
}

// Commits a decoded seed node payload to the Postgres repository
// Returns the headerId for use in persisting the token records at this blockheight
func (ber *BlockExplorerRepository) CommitPayload(payload DecodedPayload) (int64, error) {
	tx, err := ber.db.Beginx()
	if err != nil {
		return 0, err
	}
	headerId, err := ber.CommitHeaders(tx, payload.Headers, payload.BlockNumber)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = ber.CommitUncles(tx, payload.Uncles, headerId); err != nil {
		tx.Rollback()
		return 0, err
	}
	mappedTrxIds, err := ber.CommitTransactions(tx, payload.Transactions, headerId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = ber.CommitReceipts(tx, payload.Receipts, headerId, mappedTrxIds); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = ber.CommitAccounts(tx, payload.Accounts, headerId); err != nil {
		tx.Rollback()
		return 0, err
	}
	return headerId, tx.Commit()
}

func (ber *BlockExplorerRepository) CommitHeaders(tx *sqlx.Tx, headers []core.Header, blockNumber int64) (int64, error) {
	if len(headers) < 1 {
		return 0, errors.New("BlockExplorerRepository.CommitHeaders() expects at least one header")
	}
	if len(headers) > 1 {
		// We have a problem if we have two headers that are considered final
		// For now, just use the first and leave a warning that something is going wrong
		log.Errorf("BlockExplorerRepository.CommitHeaders() expects a single `final` header; %d headers found for block %d", len(headers), blockNumber)
	}
	return ber.headerRepository.CreateOrUpdateHeader(headers[0])
}

func (ber *BlockExplorerRepository) CommitUncles(tx *sqlx.Tx, uncles []core.Header, headerId int64) error {
	return nil
}

func (ber *BlockExplorerRepository) CommitTransactions(tx *sqlx.Tx, trxs []core.TransactionModel, headerId int64) (map[string]int64, error) {
	txIdMapping := make(map[string]int64)
	for _, trx := range trxs {
		txId, err := ber.headerRepository.CreateTransactionInTx(tx, headerId, trx)
		if err != nil {
			return nil, err
		}
		txIdMapping[trx.Hash] = txId
	}
	return txIdMapping, nil
}

func (ber *BlockExplorerRepository) CommitReceipts(tx *sqlx.Tx, receipts []core.Receipt, headerId int64, txIdMapping map[string]int64) error {
	for _, rct := range receipts {
		_, err := ber.headerRepository.CreateReceiptInTx(tx, headerId, txIdMapping[rct.TxHash], rct)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ber *BlockExplorerRepository) CommitAccounts(tx *sqlx.Tx, mappedAccounts map[common.Hash]state.Account, headerId int64) error {
	for accountKey, account := range mappedAccounts {
		err := ber.createAccount(tx, accountKey, account, headerId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ber *BlockExplorerRepository) createAccount(tx *sqlx.Tx, accountKey common.Hash, account state.Account, headerId int64) error {
	_, err := tx.Exec(`INSERT INTO accounts.state_accounts
						  	(header_id, account_key, balance, root, nonce, code_hash)
							VALUES ($1, $2, $3, $4, $5, $6)
							ON CONFLICT (header_id, account_key) DO UPDATE
							SET (balance, root, nonce, code_hash) = ($3::NUMERIC, $4, $5, $6)
							RETURNING id`, headerId, accountKey, account.Balance.Int64(), account.Root.Hex(), account.Nonce, account.CodeHash)
	return err
}
