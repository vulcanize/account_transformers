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

package client

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
)

type RpcClient struct {
	client  *rpc.Client
	ipcPath string
}

type BatchElem struct {
	Method string
	Args   []interface{}
	Result interface{}
	Error  error
}

func NewRpcClient(client *rpc.Client, ipcPath string) RpcClient {
	return RpcClient{
		client:  client,
		ipcPath: ipcPath,
	}
}

func (client RpcClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	//If an empty interface (or other nil object) is passed to CallContext, when the JSONRPC message is created the params will
	//be interpreted as [null]. This seems to work fine for most of the ethereum clients (which presumably ignore a null parameter.
	//Ganache however does not ignore it, and throws an 'Incorrect number of arguments' error.
	if args == nil {
		return client.client.CallContext(ctx, result, method)
	} else {
		return client.client.CallContext(ctx, result, method, args...)
	}
}

func (client RpcClient) IpcPath() string {
	return client.ipcPath
}

func (client RpcClient) SupportedModules() (map[string]string, error) {
	return client.client.SupportedModules()
}

func (client RpcClient) BatchCall(batch []BatchElem) error {
	var rpcBatch []rpc.BatchElem
	for _, batchElem := range batch {
		var newBatchElem = rpc.BatchElem{
			Result: batchElem.Result,
			Method: batchElem.Method,
			Args:   batchElem.Args,
			Error:  batchElem.Error,
		}

		rpcBatch = append(rpcBatch, newBatchElem)
	}
	return client.client.BatchCall(rpcBatch)
}
