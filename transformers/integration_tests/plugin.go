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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"plugin"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/libraries/shared/watcher"
	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	p2 "github.com/vulcanize/vulcanizedb/pkg/plugin"
	"github.com/vulcanize/vulcanizedb/pkg/plugin/helpers"

	"github.com/vulcanize/account_transformers/transformers/account/shared/test_helpers"
)

var accountConfig = config.Plugin{
	Home: "github.com/vulcanize/account_transformers",
	Transformers: map[string]config.Transformer{
		"account": {
			Path:           "transformers/account/light/initializer",
			Type:           config.EthContract,
			MigrationPath:  "db/migrations",
			MigrationRank:  0,
			RepositoryPath: "github.com/vulcanize/account_transformers",
		},
	},
	FileName: "testAccountTransformer",
	FilePath: "$GOPATH/src/github.com/vulcanize/account_transformers/transformers/integration_tests/plugin",
	Save:     false,
}

var dbConfig = config.Database{
	Hostname: "localhost",
	Port:     5432,
	Name:     "vulcanize_private",
}

type Exporter interface {
	Export() ([]transformer.EventTransformerInitializer, []transformer.StorageTransformerInitializer, []transformer.ContractTransformerInitializer)
}

var _ = Describe("Plugin test", func() {
	var g p2.Generator
	var goPath, soPath string
	var err error
	var bc core.BlockChain
	var db *postgres.DB
	var hr repositories.HeaderRepository
	viper.SetConfigName("integration")
	viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/account_transformers/environments/")

	Describe("Account Transformer Plugin", func() {
		BeforeEach(func() {
			goPath, soPath, err = accountConfig.GetPluginPaths()
			Expect(err).ToNot(HaveOccurred())
			g, err = p2.NewGenerator(accountConfig, dbConfig)
			Expect(err).ToNot(HaveOccurred())
			err = g.GenerateExporterPlugin()
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() {
			err := helpers.ClearFiles(goPath, soPath)
			Expect(err).ToNot(HaveOccurred())
		})

		Describe("GenerateTransformerPlugin", func() {
			It("It bundles the specified  TransformerInitializers into a Exporter object and creates .so", func() {
				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				event, storage, contract := exporter.Export()
				Expect(len(event)).To(Equal(0))
				Expect(len(storage)).To(Equal(0))
				Expect(len(contract)).To(Equal(1))
			})

			It("Loads our generated Exporter and uses it to import an arbitrary set of TransformerInitializers that we can execute over", func() {
				db, bc = test_helpers.SetupDBandBC()
				defer test_helpers.TearDown(db)

				hr = repositories.NewHeaderRepository(db)
				header1, err := bc.GetHeaderByNumber(6791668)
				Expect(err).ToNot(HaveOccurred())
				_, err = hr.CreateOrUpdateHeader(header1)
				Expect(err).ToNot(HaveOccurred())

				plug, err := plugin.Open(soPath)
				Expect(err).ToNot(HaveOccurred())
				symExporter, err := plug.Lookup("Exporter")
				Expect(err).ToNot(HaveOccurred())
				exporter, ok := symExporter.(Exporter)
				Expect(ok).To(Equal(true))
				_, _, contractInitializers := exporter.Export()

				w := watcher.NewContractWatcher(db, bc)
				w.AddTransformers(contractInitializers)
				err = w.Execute()
				Expect(err).ToNot(HaveOccurred())

				type coinBalanceRecord struct {
					Address     []byte `db:"address_hash"`
					BlockNumber int64  `db:"block_number"`
					Value       string
				}

				var coinRecords []coinBalanceRecord
				rows, err := db.Queryx(`SELECT * FROM accounts.address_coin_balances WHERE block_number = $1`, 6791668)
				Expect(err).ToNot(HaveOccurred())
				for rows.Next() {
					record := new(coinBalanceRecord)
					err = rows.StructScan(record)
					Expect(err).ToNot(HaveOccurred())
					coinRecords = append(coinRecords, *record)
				}
				err = rows.Err()
				Expect(err).ToNot(HaveOccurred())
				rows.Close()

				type tokenBalanceRecord struct {
					Address         []byte `db:"address_hash"`
					BlockNumber     int64  `db:"block_number"`
					ContractAddress []byte `db:"token_contract_address_hash"`
					Value           string
				}
				var tokenRecords []tokenBalanceRecord
				rows, err = db.Queryx(`SELECT * FROM accounts.address_token_balances WHERE block_number = $1`, 6791668)
				Expect(err).ToNot(HaveOccurred())
				for rows.Next() {
					record := new(tokenBalanceRecord)
					err = rows.StructScan(record)
					Expect(err).ToNot(HaveOccurred())
					tokenRecords = append(tokenRecords, *record)
				}
				err = rows.Err()
				Expect(err).ToNot(HaveOccurred())
				rows.Close()
			})
		})
	})
})
