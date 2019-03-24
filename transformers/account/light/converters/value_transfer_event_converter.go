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
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/account_transformers/transformers/account/light/models"
	"github.com/vulcanize/account_transformers/transformers/account/shared"
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
)

type ValueTransferConverter interface {
	Convert(ethLogs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error)
}

type valueTransferConverter struct {
	boundEvent            constants.Label
	mappedEquivalentAddrs map[common.Address][]common.Address
}

func NewValueTransferConverter(mappedEquivalentAddrs map[common.Address][]common.Address) *valueTransferConverter {
	return &valueTransferConverter{
		mappedEquivalentAddrs: mappedEquivalentAddrs,
	}
}

func (c *valueTransferConverter) Convert(ethLogs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	groupedLogs := c.group(ethLogs)
	models := make([]models.ValueTransferModel, 0)
	for label, logs := range groupedLogs {
		c.boundEvent = label
		unpackedModels, err := c.unpack(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
		models = append(models, unpackedModels...)
	}
	return models, nil
}

func (c *valueTransferConverter) group(ethLogs []types.Log) map[constants.Label][]types.Log {
	groupedLogs := make(map[constants.Label][]types.Log, 0)
	for _, log := range ethLogs {
		if len(log.Topics) < 1 {
			continue
		}
		label := constants.NewLabelFromSignature(log.Topics[0].Hex())
		if label.Name() == "" {
			continue
		}
		groupedLogs[label] = append(groupedLogs[label], log)
	}
	return groupedLogs
}

func (c *valueTransferConverter) unpack(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var err error
	switch c.boundEvent {
	case constants.Transfer:
		records, err = c.convertTransfers(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Mint:
		records, err = c.convertMints(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.WipedAccount:
		records, err = c.convertWipedAccounts(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Burn:
		records, err = c.convertBurns(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.WipeBlacklistedAccount:
		records, err = c.convertWipeBlacklistedAccounts(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.DestroyedBlackFunds:
		records, err = c.convertDestroyedBlackFunds(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.TransferFrom:
		records, err = c.convertTransferFroms(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Deposit:
		records, err = c.convertDeposits(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Withdrawal:
		records, err = c.convertWithdrawals(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	}
	return records, nil
}

func (c *valueTransferConverter) convertTransfers(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Transfer0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Transfer1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Transfer2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Transfer3Indexed.String()))
			if err != nil {
				return nil, err
			}
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedTransfer shared.TransferEntity
		err := boundContract.UnpackLog(&unpackedTransfer, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedTransfer.From.Hex(),
			Dst:              unpackedTransfer.To.Hex(),
			Amount:           unpackedTransfer.Value.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertMints(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Mint0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Mint1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Mint2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `Mint` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedMint shared.MintEntity
		err := boundContract.UnpackLog(&unpackedMint, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              "0x0",
			Dst:              unpackedMint.To.Hex(),
			Amount:           unpackedMint.Amount.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertWipedAccounts(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipedAccount0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipedAccount1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipedAccount2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `WipedAccount` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedWipedAccount shared.WipedAccountEntity
		err := boundContract.UnpackLog(&unpackedWipedAccount, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedWipedAccount.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWipedAccount.Balance.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertBurns(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Burn0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Burn1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Burn2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `Burn` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedBurn shared.BurnEntity
		err := boundContract.UnpackLog(&unpackedBurn, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedBurn.Burner.Hex(),
			Dst:              "0x0",
			Amount:           unpackedBurn.Value.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertWipeBlacklistedAccounts(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipeBlacklistedAccount0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipeBlacklistedAccount1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.WipeBlacklistedAccount2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `WipeBlacklistedAccount` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedWBAE shared.WipeBlacklistedAccountEntity
		err := boundContract.UnpackLog(&unpackedWBAE, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedWBAE.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWBAE.Balance.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertDestroyedBlackFunds(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.DestroyedBlackFunds0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.DestroyedBlackFunds1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.DestroyedBlackFunds2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `DestroyedBlackFunds` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedDBF shared.DestroyedBlackFundsEntity
		err := boundContract.UnpackLog(&unpackedDBF, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedDBF.BlackListedUser.Hex(),
			Dst:              "0x0",
			Amount:           unpackedDBF.Balance.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertTransferFroms(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.TransferFrom0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.TransferFrom1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.TransferFrom2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			parsedABI, err = abi.JSON(strings.NewReader(constants.TransferFrom3Indexed.String()))
			if err != nil {
				return nil, err
			}
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedTransferFrom shared.TransferFromEntity
		err := boundContract.UnpackLog(&unpackedTransferFrom, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedTransferFrom.From.Hex(),
			Dst:              unpackedTransferFrom.To.Hex(),
			Amount:           unpackedTransferFrom.Value.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertDeposits(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Deposit0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Deposit1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Deposit2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `Deposit` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedDeposit shared.DepositEntity
		err := boundContract.UnpackLog(&unpackedDeposit, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              "0x0",
			Dst:              unpackedDeposit.Dst.Hex(),
			Amount:           unpackedDeposit.Wad.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertWithdrawals(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
	var parsedABI abi.ABI
	var err error
	for _, log := range logs {
		switch len(log.Topics) {
		case 1:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Withdrawal0Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 2:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Withdrawal1Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 3:
			parsedABI, err = abi.JSON(strings.NewReader(constants.Withdrawal2Indexed.String()))
			if err != nil {
				return nil, err
			}
		case 4:
			return nil, errors.New("converter: `Withdrawal` event cannot have 4 topics")
		}
		boundContract := bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		var unpackedDeposit shared.WithdrawalEntity
		err := boundContract.UnpackLog(&unpackedDeposit, c.boundEvent.Name(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.Name(),
			BlockNumber:      blockNumber,
			Src:              unpackedDeposit.Src.Hex(),
			Dst:              "0x0",
			Amount:           unpackedDeposit.Wad.String(),
			Contract:         c.getEquivalent(log.Address),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) getEquivalent(addr common.Address) string {
	for topAddr, equivalents := range c.mappedEquivalentAddrs {
		for _, equivalent := range equivalents {
			if equivalent == addr {
				return topAddr.Hex()
			}
		}
	}
	return addr.Hex() // If we find no top level equivalency to map this token address to, return it
}
