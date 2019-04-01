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

package poller

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"

	"github.com/vulcanize/account_transformers/transformers/account/shared"
)

type AccountPoller interface {
	PollAccount(addr common.Address, headers []core.Header) ([]shared.CoinBalanceRecord, error)
}

type accountPoller struct {
	db               *postgres.DB
	headerRepository repositories.HeaderRepository
	blockChain       core.BlockChain
	balanceCache     map[common.Address]*big.Int
}

func NewAccountPoller(db *postgres.DB, bc core.BlockChain) *accountPoller {
	return &accountPoller{
		db:               db,
		headerRepository: repositories.NewHeaderRepository(db),
		blockChain:       bc,
		balanceCache:     make(map[common.Address]*big.Int),
	}
}

func (ap *accountPoller) PollAccount(addr common.Address, headers []core.Header) ([]shared.CoinBalanceRecord, error) {
	balanceRecords := make([]shared.CoinBalanceRecord, 0, len(headers))
	for _, header := range headers {
		if ap.balanceCache[addr] == nil {
			ap.balanceCache[addr] = big.NewInt(0)
		}
		record := shared.CoinBalanceRecord{
			BlockNumber: header.BlockNumber,
			Address:     addr.Bytes(),
			HeaderID:    header.Id,
		}
		balance, err := ap.blockChain.GetAccountBalance(addr, big.NewInt(header.BlockNumber))
		if err != nil {
			return nil, err
		}
		if ap.balanceCache[addr].String() != balance.String() {
			err = ap.pollTx(addr, header.BlockNumber, header.Id)
			if err != nil {
				return nil, err
			}
		}
		ap.balanceCache[addr] = balance
		record.Value = balance.String()
		balanceRecords = append(balanceRecords, record)
	}

	return balanceRecords, nil
}

func (ap *accountPoller) pollTx(addr common.Address, blockNumber, headerID int64) error {
	blk, err := ap.blockChain.GetBlockByNumber(blockNumber)
	if err != nil {
		return err
	}
	tx, err := ap.db.Beginx()
	if err != nil {
		return err
	}
	for _, trx := range blk.Transactions {
		if strings.ToLower(trx.From) == strings.ToLower(addr.String()) || strings.ToLower(trx.To) == strings.ToLower(addr.String()) {
			txId, err := ap.headerRepository.CreateTransactionInTx(tx, headerID, trx)
			if err != nil {
				tx.Rollback()
				return err
			}
			_, err = ap.headerRepository.CreateReceiptInTx(tx, headerID, txId, trx.Receipt)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit()
}
