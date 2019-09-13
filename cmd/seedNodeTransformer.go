// Copyright © 2019 Vulcanize, Inc
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vulcanize/account_transformers/transformers/account/seed"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth/client"
	"github.com/vulcanize/vulcanizedb/utils"
)

// seedNodeTransformerCmd represents the seedNodeTransformer command
var seedNodeTransformerCmd = &cobra.Command{
	Use:   "seedNodeTransformer",
	Short: "Account transformer that uses the VulcanizeDB seed node",
	Long: `This command spins up an account transformer which receives data from
a VulcanizeDB seed node, it converts this data into more useful models and persists
them in Postgres`,
	Run: func(cmd *cobra.Command, args []string) {
		seedNodeTransformer()
	},
}

func seedNodeTransformer() {
	// Create rpc client for the seed node
	rpcClient := getRpcClient()
	gethNode := new(core.Node)
	// Fetch the geth node info from the remote seed node for use in configuring the database
	// Not sure if "statediff" is supposed to be the "method" arg, if it follows the same convention as subscriptions
	// it would be and the first arg would actually be the method name ("node"); need to test this still
	err := rpcClient.CallContext(context.Background(), gethNode, "statediff", "node")
	if err != nil {
		log.Fatal(err)
	}
	transIniter := seed.AccountTransformer{}
	db := utils.LoadPostgres(databaseConfig, *gethNode)
	subCon := seed.DefaultConfig
	trans := transIniter.NewTransformer(&db, subCon, rpcClient)
	if err := trans.Init(); err != nil {
		log.Fatal(err)
	}
	if err := trans.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(seedNodeTransformerCmd)
}

func getRpcClient() core.RpcClient {
	vulcPath := viper.GetString("subscription.path")
	if vulcPath == "" {
		vulcPath = "ws://127.0.0.1:8080" // default to and try the default ws url if no path is provided
	}
	rawRpcClient, err := rpc.Dial(vulcPath)
	if err != nil {
		log.Fatal(err)
	}
	return client.NewRpcClient(rawRpcClient, vulcPath)
}
