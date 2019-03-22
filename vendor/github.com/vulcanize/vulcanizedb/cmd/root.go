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

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/pkg/config"
	"github.com/vulcanize/vulcanizedb/pkg/geth"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	vRpc "github.com/vulcanize/vulcanizedb/pkg/geth/converters/rpc"
	"github.com/vulcanize/vulcanizedb/pkg/geth/node"
)

var (
	cfgFile             string
	databaseConfig      config.Database
	genConfig           config.Plugin
	ipc                 string
	levelDbPath         string
	startingBlockNumber int64
	storageDiffsPath    string
	syncAll             bool
	endingBlockNumber   int64
	recheckHeadersArg   bool
)

const (
	pollingInterval  = 7 * time.Second
	validationWindow = 15
)

var rootCmd = &cobra.Command{
	Use:              "vulcanizedb",
	PersistentPreRun: database,
}

func Execute() {
	log.Info("----- Starting vDB -----")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func database(cmd *cobra.Command, args []string) {
	ipc = viper.GetString("client.ipcpath")
	levelDbPath = viper.GetString("client.leveldbpath")
	storageDiffsPath = viper.GetString("filesystem.storageDiffsPath")
	databaseConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}
	viper.Set("database.config", databaseConfig)
}

func init() {
	cobra.OnInitialize(initConfig)
	// When searching for env variables, replace dots in config keys with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file location")
	rootCmd.PersistentFlags().String("database-name", "vulcanize_public", "database name")
	rootCmd.PersistentFlags().Int("database-port", 5432, "database port")
	rootCmd.PersistentFlags().String("database-hostname", "localhost", "database hostname")
	rootCmd.PersistentFlags().String("database-user", "", "database user")
	rootCmd.PersistentFlags().String("database-password", "", "database password")
	rootCmd.PersistentFlags().String("client-ipcPath", "", "location of geth.ipc file")
	rootCmd.PersistentFlags().String("client-levelDbPath", "", "location of levelDb chaindata")
	rootCmd.PersistentFlags().String("filesystem-storageDiffsPath", "", "location of storage diffs csv file")
	rootCmd.PersistentFlags().String("exporter-name", "exporter", "name of exporter plugin")

	viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("database-name"))
	viper.BindPFlag("database.port", rootCmd.PersistentFlags().Lookup("database-port"))
	viper.BindPFlag("database.hostname", rootCmd.PersistentFlags().Lookup("database-hostname"))
	viper.BindPFlag("database.user", rootCmd.PersistentFlags().Lookup("database-user"))
	viper.BindPFlag("database.password", rootCmd.PersistentFlags().Lookup("database-password"))
	viper.BindPFlag("client.ipcPath", rootCmd.PersistentFlags().Lookup("client-ipcPath"))
	viper.BindPFlag("client.levelDbPath", rootCmd.PersistentFlags().Lookup("client-levelDbPath"))
	viper.BindPFlag("filesystem.storageDiffsPath", rootCmd.PersistentFlags().Lookup("filesystem-storageDiffsPath"))
	viper.BindPFlag("exporter.fileName", rootCmd.PersistentFlags().Lookup("exporter-name"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		noConfigError := "No config file passed with --config flag"
		fmt.Println("Error: ", noConfigError)
		log.Fatal(noConfigError)
		os.Exit(1)
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n\n", viper.ConfigFileUsed())
	} else {
		invalidConfigError := "Couldn't read config file"
		fmt.Println("Error: ", invalidConfigError)
		log.Fatal(invalidConfigError)
		os.Exit(1)
	}
}

func getBlockChain() *geth.BlockChain {
	rawRpcClient, err := rpc.Dial(ipc)

	if err != nil {
		log.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := vRpc.NewRpcTransactionConverter(ethClient)
	return geth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
}
