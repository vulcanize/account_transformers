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

package converters

import (
	"encoding/json"
	"github.com/vulcanize/account_transformers/transformers/account/light/models"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/account_transformers/transformers/account/shared"
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
)

type ValueTransferConverter interface {
	Convert(ethLogs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error)
}

type valueTransferConverter struct {
	boundContract *bind.BoundContract
	boundEvent    constants.Label
}

func NewValueTransferConverter() *valueTransferConverter {
	return &valueTransferConverter{}
}

func (c *valueTransferConverter) Convert(ethLogs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	groupedLogs := c.group(ethLogs)
	models := make([]models.ValueTransferModel, 0)
	for label, logs := range groupedLogs {
		parsedABI, err := geth.ParseAbi(label.ABI())
		if err != nil {
			return nil, err
		}
		c.boundContract = bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		c.boundEvent = label
		models, err := c.unpack(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
		models = append(models, models...)
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
		if label.String() == "" {
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
	case constants.Transfer, constants.TransferDev:
		records, err = c.convertTransfers(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Mint, constants.MintDev:
		records, err = c.convertMints(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.WipedAccount, constants.WipedAccountDev:
		records, err = c.convertWipedAccounts(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Burn, constants.BurnDev:
		records, err = c.convertBurns(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.WipeBlacklistedAccount, constants.WipeBlacklistedAccountDev:
		records, err = c.convertWipeBlacklistedAccounts(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.DestroyedBlackFunds, constants.DestroyedBlackFundsDev:
		records, err = c.convertDestroyedBlackFunds(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Issue, constants.IssueDev:
		records, err = c.convertIssues(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.Redeem, constants.RedeemDev:
		records, err = c.convertRedeems(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	case constants.TransferFrom, constants.TransferFromDev:
		records, err = c.convertTransferFroms(logs, headerID, blockNumber)
		if err != nil {
			return nil, err
		}
	}
	return records, nil
}

func (c *valueTransferConverter) convertTransfers(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedTransfer.From.Hex(),
			Dst:              unpackedTransfer.To.Hex(),
			Amount:           unpackedTransfer.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertMints(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              "0x0",
			Dst:              unpackedMint.To.Hex(),
			Amount:           unpackedMint.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertWipedAccounts(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedWipedAccount.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWipedAccount.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertBurns(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedBurn.Burner.Hex(),
			Dst:              "0x0",
			Amount:           unpackedBurn.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertWipeBlacklistedAccounts(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedWBAE.Account.Hex(),
			Dst:              "0x0",
			Amount:           unpackedWBAE.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertDestroyedBlackFunds(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedDBF.BlackListedUser.Hex(),
			Dst:              "0x0",
			Amount:           unpackedDBF.Balance.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertIssues(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              "0x0",
			Dst:              log.Address.Hex(),
			Amount:           unpackedIssue.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertRedeems(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              log.Address.Hex(),
			Dst:              "0x0",
			Amount:           unpackedRedeem.Amount.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}

func (c *valueTransferConverter) convertTransferFroms(logs []types.Log, headerID, blockNumber int64) ([]models.ValueTransferModel, error) {
	var records []models.ValueTransferModel
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
		records = append(records, models.ValueTransferModel{
			HeaderID:         headerID,
			Name:             c.boundEvent.String(),
			BlockNumber:      blockNumber,
			Src:              unpackedTransferFrom.From.Hex(),
			Dst:              unpackedTransferFrom.To.Hex(),
			Amount:           unpackedTransferFrom.Value.String(),
			Contract:         log.Address.Hex(),
			LogIndex:         log.Index,
			TransactionIndex: log.TxIndex,
			Raw:              raw,
		})
	}
	return records, nil
}
