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

package converters

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
)

type ValueTransferConverter interface {
	Convert(ethLogs []types.Log, headerID int64) ([]models.ValueTransferModel, error)
}

type valueTransferConverter struct {
	boundEvent            constants.Event
	mappedEquivalentAddrs map[common.Address][]common.Address
	boundABI              abi.ABI
	boundContract         *bind.BoundContract
}

func NewValueTransferConverter(abiStr string, mappedEquivalentAddrs map[common.Address][]common.Address) (*valueTransferConverter, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedAbi, nil, nil, nil)
	return &valueTransferConverter{
		mappedEquivalentAddrs: mappedEquivalentAddrs,
		boundABI:              parsedAbi,
		boundContract:         boundContract,
	}, nil
}

func (c *valueTransferConverter) Convert(ethLogs []types.Log, headerID int64) ([]models.ValueTransferModel, error) {
	transferModels := make([]models.ValueTransferModel, 0, len(ethLogs))
	for _, log := range ethLogs {
		topicCount := len(log.Topics)
		if topicCount < 1 { // If we don't event have a topic0 then something is very wrong
			continue
		}
		label := constants.NewLabelFromSignature(log.Topics[0].Hex())
		if label.Name() == "" { // This isn't an event we want (shouldn't happen sense we filter by topic to get these events to begin with)
			continue
		}
		c.boundEvent = constants.GetEventFromLabel(label)
		transferModel, err := c.unpack(log, headerID, topicCount)
		if err != nil {
			return nil, err
		}
		transferModels = append(transferModels, transferModel)
	}

	return transferModels, nil
}

func (c *valueTransferConverter) unpack(log types.Log, headerID int64, topicCount int) (models.ValueTransferModel, error) {
	if len(c.boundEvent.Names) < topicCount {
		return models.ValueTransferModel{}, errors.New(fmt.Sprintf("Event type %s cannot have %d number of topics", c.boundEvent.Label.Name(), topicCount))
	}

	unpackMap := make(map[string]interface{})
	abiEventName := c.boundEvent.Names[topicCount-1]
	err := c.boundContract.UnpackLogIntoMap(unpackMap, abiEventName, log)
	if err != nil {
		return models.ValueTransferModel{}, fmt.Errorf("unable to unpack event %s\r\nlog: %v\r\nerror: %v", abiEventName, log, err)
	}
	raw, err := json.Marshal(log)
	if err != nil {
		return models.ValueTransferModel{}, err
	}
	var src, dst common.Address
	amount := big.NewInt(0)
	var ok bool
	if unpackMap["from"] != nil {
		src, ok = unpackMap["from"].(common.Address)
		if !ok {
			return models.ValueTransferModel{}, fmt.Errorf("`From` field in unpacked map should be of type %T but is of type %T", log.Address, unpackMap["from"])
		}
	}
	if unpackMap["to"] != nil {
		dst, ok = unpackMap["to"].(common.Address)
		if !ok {
			return models.ValueTransferModel{}, fmt.Errorf("`To` field in unpacked map should be of type %T but is of type %T", log.Address, unpackMap["to"])
		}
	}
	if unpackMap["value"] != nil {
		amount, ok = unpackMap["value"].(*big.Int)
		if !ok {
			return models.ValueTransferModel{}, fmt.Errorf("`Amount` field in unpacked map should be of type %T but is of type %T", big.NewInt(0), unpackMap["value"])
		}
	}
	return models.ValueTransferModel{
		HeaderID:         headerID,
		Name:             c.boundEvent.Label.Name(),
		BlockNumber:      log.BlockNumber,
		Src:              src.Bytes(),
		Dst:              dst.Bytes(),
		Amount:           amount.String(),
		Contract:         c.getEquivalent(log.Address),
		LogIndex:         log.Index,
		TransactionIndex: log.TxIndex,
		Raw:              raw,
	}, nil
}

func (c *valueTransferConverter) getEquivalent(addr common.Address) []byte {
	for topAddr, equivalents := range c.mappedEquivalentAddrs {
		for _, equivalent := range equivalents {
			if equivalent == addr {
				return topAddr.Bytes()
			}
		}
	}
	return addr.Bytes() // If we find no top level equivalency to map this token address to, return it
}
