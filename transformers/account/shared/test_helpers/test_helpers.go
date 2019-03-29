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

package test_helpers

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	rpc2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"

	"github.com/vulcanize/account_transformers/test_config"
)

func SetupDBandBC() (*postgres.DB, core.BlockChain) {
	con := test_config.InfuraClient
	Expect(con).ToNot(Equal(nil))
	infuraIPC := con.IPCPath
	rawRpcClient, err := rpc.Dial(infuraIPC)
	Expect(err).NotTo(HaveOccurred())
	rpcClient := client.NewRpcClient(rawRpcClient, infuraIPC)
	ethClient := ethclient.NewClient(rawRpcClient)
	blockChainClient := client.NewEthClient(ethClient)
	node := node.MakeNode(rpcClient)
	transactionConverter := rpc2.NewRpcTransactionConverter(ethClient)
	blockChain := geth.NewBlockChain(blockChainClient, rpcClient, node, transactionConverter)

	db, err := postgres.NewDB(config.Database{
		Hostname: "localhost",
		Name:     "vulcanize_private",
		Port:     5432,
	}, blockChain.Node())
	Expect(err).NotTo(HaveOccurred())

	return db, blockChain
}

func TearDown(db *postgres.DB) {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DROP TABLE checked_headers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`CREATE TABLE checked_headers (id SERIAL PRIMARY KEY, header_id INTEGER UNIQUE NOT NULL REFERENCES headers (id) ON DELETE CASCADE);`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM accounts.addresses`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM accounts.token_value_transfers`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM accounts.address_coin_balances`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM accounts.addresses`)
	Expect(err).NotTo(HaveOccurred())

	_, err = tx.Exec(`DELETE FROM accounts.watched_contracts`)
	Expect(err).NotTo(HaveOccurred())

	err = tx.Commit()
	Expect(err).NotTo(HaveOccurred())
}
