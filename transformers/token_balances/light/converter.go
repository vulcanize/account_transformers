// VulcanizeDB
// Copyright © 2018 Vulcanize

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

package light

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/account_transformers/transformers/token_balances/shared"
	"github.com/vulcanize/account_transformers/transformers/token_balances/shared/constants"
)

type Converter interface {
	Convert(ethLogs []types.Log, headerID int64) ([]ValueTransferModel, error)
}

type converter struct {
	boundContract *bind.BoundContract
	boundEvent    constants.Label
}

func NewConverter(abis map[constants.Label]string) *converter {
	return &converter{}
}

func (c *converter) Convert(ethLogs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	groupedLogs := c.group(ethLogs)
	models := make([]ValueTransferModel, 0)
	for label, logs := range groupedLogs {
		parsedABI, err := geth.ParseAbi(label.ABI())
		if err != nil {
			return nil, err
		}
		c.boundContract = bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		c.boundEvent = label
		models, err := c.unpack(logs)
		if err != nil {
			return nil, err
		}
		models = append(models, models...)
	}
	return models, nil
}

func (c *converter) group(ethLogs []types.Log) map[constants.Label][]types.Log {
	groupedLogs := make(map[constants.Label][]types.Log, 0)
	for _, log := range ethLogs {
		if len(log.Topics) < 1 {
			continue
		}
		label := constants.NewLabelFromSignature(log.Topics[0].Hex())
		if label.String() == "" {
			continue
		}
		groupedLogs[label] = append(groupedLogs[label], log)
	}
	return groupedLogs
}

func (c *converter) unpack(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	var err error
	switch c.boundEvent {
	case constants.Transfer, constants.TransferDev:
		models, err = c.convertTransfers(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.Mint, constants.MintDev:
		models, err = c.convertMints(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.WipedAccount, constants.WipedAccountDev:
		models, err = c.convertWipedAccounts(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.Burn, constants.BurnDev:
		models, err = c.convertBurns(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.WipeBlacklistedAccount, constants.WipeBlacklistedAccountDev:
		models, err = c.convertWipeBlacklistedAccounts(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.DestroyedBlackFunds, constants.DestroyedBlackFundsDev:
		models, err = c.convertDestroyedBlackFunds(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.Issue, constants.IssueDev:
		models, err = c.convertIssues(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.Redeem, constants.RedeemDev:
		models, err = c.convertRedeems(logs, headerID)
		if err != nil {
			return nil, err
		}
	case constants.TransferFrom, constants.TransferFromDev:
		models, err = c.convertTransferFroms(logs, headerID)
		if err != nil {
			return nil, err
		}
	}
	return models, nil
}

func (c *converter) convertTransfers(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedTransfer shared.TransferEntity
		err := c.boundContract.UnpackLog(&unpackedTransfer, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedTransfer.From.Hex(),
			Dst:              unpackedTransfer.To.Hex(),
			Amount:           unpackedTransfer.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertMints(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedMint shared.MintEntity
		err := c.boundContract.UnpackLog(&unpackedMint, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              "0x0",
			Dst:              unpackedMint.To.Hex(),
			Amount:           unpackedMint.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertWipedAccounts(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedWipedAccount shared.WipedAccountEntity
		err := c.boundContract.UnpackLog(&unpackedWipedAccount, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedWipedAccount.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWipedAccount.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertBurns(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedBurn shared.BurnEntity
		err := c.boundContract.UnpackLog(&unpackedBurn, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedBurn.Burner.Hex(),
			Dst:              "0x0",
			Amount:           unpackedBurn.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertWipeBlacklistedAccounts(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedWBAE shared.WipeBlacklistedAccountEntity
		err := c.boundContract.UnpackLog(&unpackedWBAE, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedWBAE.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWBAE.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertDestroyedBlackFunds(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedDBF shared.DestroyedBlackFundsEntity
		err := c.boundContract.UnpackLog(&unpackedDBF, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedDBF.BlackListedUser.Hex(),
			Dst:              "0x0",
			Amount:           unpackedDBF.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertIssues(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedIssue shared.IssueEntity
		err := c.boundContract.UnpackLog(&unpackedIssue, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              "0x0",
			Dst:              log.Address.Hex(),
			Amount:           unpackedIssue.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertRedeems(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedRedeem shared.RedeemEntity
		err := c.boundContract.UnpackLog(&unpackedRedeem, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              log.Address.Hex(),
			Dst:              "0x0",
			Amount:           unpackedRedeem.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}

func (c *converter) convertTransferFroms(logs []types.Log, headerID int64) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, log := range logs {
		var unpackedTransferFrom shared.TransferFromEntity
		err := c.boundContract.UnpackLog(&unpackedTransferFrom, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			Src:              unpackedTransferFrom.From.Hex(),
			Dst:              unpackedTransferFrom.To.Hex(),
			Amount:           unpackedTransferFrom.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return models, nil
}
