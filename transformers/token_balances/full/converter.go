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

package full

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/contract_watcher/shared/helpers"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/account_transformers/transformers/token_balances/shared"
	"github.com/vulcanize/account_transformers/transformers/token_balances/shared/constants"
)

type Converter interface {
	Convert(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error)
}

type converter struct {
	boundContract *bind.BoundContract
	boundEvent     constants.Label
}

func NewConverter(abis map[constants.Label]string) *converter {
	return &converter{}
}

func (c *converter) Convert(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	groupedWatchedEvents := c.group(watchedEvents)
	convertedModels := make([]ValueTransferModel, 0)
	for label, watchedEvents := range groupedWatchedEvents {
		parsedABI, err := geth.ParseAbi(label.ABI())
		if err != nil {
			return nil, err
		}
		c.boundContract = bind.NewBoundContract(common.HexToAddress("0x0"), parsedABI, nil, nil, nil)
		c.boundEvent = label
		models, err := c.unpack(watchedEvents)
		if err != nil {
			return nil, err
		}
		convertedModels = append(convertedModels, models...)
	}
	return convertedModels, nil
}

func (c *converter) group(watchedEvents []*core.WatchedEvent) map[constants.Label][]*core.WatchedEvent {
	groupedWatchedEvents := make(map[constants.Label][]*core.WatchedEvent, 0)
	for _, we := range watchedEvents {
		label := constants.NewLabel(we.Name)
		if label.String() == "" {
			continue
		}
		groupedWatchedEvents[label] = append(groupedWatchedEvents[label], we)
	}
	return groupedWatchedEvents
}

func (c *converter) unpack(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	var err error
	switch c.boundEvent {
	case constants.Transfer, constants.TransferDev:
		models, err = c.convertTransfers(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.Mint, constants.MintDev:
		models, err = c.convertMints(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.WipedAccount, constants.WipedAccountDev:
		models, err = c.convertWipedAccounts(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.Burn, constants.BurnDev:
		models, err = c.convertBurns(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.WipeBlacklistedAccount, constants.WipeBlacklistedAccountDev:
		models, err = c.convertWipeBlacklistedAccounts(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.DestroyedBlackFunds, constants.DestroyedBlackFundsDev:
		models, err = c.convertDestroyedBlackFunds(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.Issue, constants.IssueDev:
		models, err = c.convertIssues(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.Redeem, constants.RedeemDev:
		models, err = c.convertRedeems(watchedEvents)
		if err != nil {
			return nil, err
		}
	case constants.TransferFrom, constants.TransferFromDev:
		models, err = c.convertTransferFroms(watchedEvents)
		if err != nil {
			return nil, err
		}
	}
	return models, nil
}

func (c *converter) convertTransfers(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedTransfer shared.TransferEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedTransfer, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedTransfer.From.Hex(),
			Dst: unpackedTransfer.To.Hex(),
			Amount: unpackedTransfer.Value.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertMints(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedMint shared.MintEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedMint, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: "0x0",
			Dst: unpackedMint.To.Hex(),
			Amount: unpackedMint.Amount.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertWipedAccounts(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedWipedAccount shared.WipedAccountEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedWipedAccount, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedWipedAccount.Account.Hex(),
			Dst: "0x0",
			Amount: unpackedWipedAccount.Balance.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertBurns(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedBurn shared.BurnEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedBurn, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedBurn.Burner.Hex(),
			Dst: "0x0",
			Amount: unpackedBurn.Value.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertWipeBlacklistedAccounts(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedWBAE shared.WipeBlacklistedAccountEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedWBAE, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedWBAE.Account.Hex(),
			Dst: "0x0",
			Amount: unpackedWBAE.Balance.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertDestroyedBlackFunds(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedDBF shared.DestroyedBlackFundsEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedDBF, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedDBF.BlackListedUser.Hex(),
			Dst: "0x0",
			Amount: unpackedDBF.Balance.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertIssues(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedIssue shared.IssueEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedIssue, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: "0x0",
			Dst: we.Address,
			Amount: unpackedIssue.Amount.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertRedeems(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedRedeem shared.RedeemEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedRedeem, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: we.Address,
			Dst: "0x0",
			Amount: unpackedRedeem.Amount.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}

func (c *converter) convertTransferFroms(watchedEvents []*core.WatchedEvent) ([]ValueTransferModel, error) {
	var models []ValueTransferModel
	for _, we := range watchedEvents {
		var unpackedTransferFrom shared.TransferFromEntity
		log := helpers.ConvertToLog(*we)
		err := c.boundContract.UnpackLog(&unpackedTransferFrom, c.boundEvent.Event(), log)
		if err != nil {
			return nil, err
		}
		models = append(models, ValueTransferModel{
			Name: we.Name,
			Src: unpackedTransferFrom.From.Hex(),
			Dst: unpackedTransferFrom.To.Hex(),
			Amount: unpackedTransferFrom.Value.String(),
			Contract: we.Address,
			VulcanizeLogID: we.LogID,
		})
	}
	return models, nil
}