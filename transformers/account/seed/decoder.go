package seed

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	common2 "github.com/vulcanize/vulcanizedb/pkg/geth/converters/common"
	"github.com/vulcanize/vulcanizedb/pkg/ipfs"
)

// RlpDecoder implements methods for decoding ETH rlp
type Decoder interface {
	DecodePayload(payload ipfs.ResponsePayload) (DecodedPayload, error)
}

// SeedNodeDecoder implements methods for decoding ETH rlp from the seed node
type SeedNodeDecoder struct {
	HeaderConverter common2.HeaderConverter
}

// NewSeedNodeDecoder returns a pointer to a new SeedNodeDecoder
func NewSeedNodeDecoder() *SeedNodeDecoder {
	return &SeedNodeDecoder{}
}

// DecodePayload decodes all fields in the
func (snd *SeedNodeDecoder) DecodePayload(payload ipfs.ResponsePayload) (DecodedPayload, error) {
	headers, err := snd.DecodeHeaders(payload.HeadersRlp)
	if err != nil {
		return DecodedPayload{}, err
	}
	uncles, err := snd.DecodeUncles(payload.UnclesRlp)
	if err != nil {
		return DecodedPayload{}, err
	}
	transactions, err := snd.DecodeTransactions(payload.TransactionsRlp)
	if err != nil {
		return DecodedPayload{}, err
	}
	receipts, ethLogs, err := snd.DecodeReceipts(payload.ReceiptsRlp)
	if err != nil {
		return DecodedPayload{}, err
	}
	accounts, err := snd.DecodeStateAccounts(payload.StateNodesRlp)
	if err != nil {
		return DecodedPayload{}, err
	}
	return DecodedPayload{
		Headers:      headers,
		Uncles:       uncles,
		Transactions: transactions,
		Receipts:     receipts,
		Logs:         ethLogs,
		Accounts:     accounts,
	}, nil
}

// DecodeHeader decodes header rlp into a types.Header
func (snd *SeedNodeDecoder) DecodeHeaders(headersRlp [][]byte) ([]core.Header, error) {
	headers := make([]core.Header, len(headersRlp))
	for _, headerRlp := range headersRlp {
		header := new(types.Header)
		err := rlp.DecodeBytes(headerRlp, header)
		if err != nil {
			return nil, err
		}
		headers = append(headers, snd.HeaderConverter.Convert(header, header.Hash().Hex()))
	}
	return headers, nil
}

// DecodeUncles decodes header rlp into a types.Header
func (snd *SeedNodeDecoder) DecodeUncles(unclesRlp [][]byte) ([]core.Header, error) {
	uncles := make([]core.Header, len(unclesRlp))
	for _, headerRlp := range unclesRlp {
		uncle := new(types.Header)
		err := rlp.DecodeBytes(headerRlp, uncle)
		if err != nil {
			return nil, err
		}
		uncles = append(uncles, snd.HeaderConverter.Convert(uncle, uncle.Hash().Hex()))
	}
	return uncles, nil
}

// DecodeTransactions decodes transaction rlps into type.Transactions
func (snd *SeedNodeDecoder) DecodeTransactions(trxRlps [][]byte) ([]core.TransactionModel, error) {
	transactions := make([]types.Transaction, 0, len(trxRlps))
	for _, trxRlp := range trxRlps {
		transaction := new(types.Transaction)
		err := rlp.DecodeBytes(trxRlp, transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, *transaction)
	}
	return snd.convertTrxs(transactions)
}

func (snd *SeedNodeDecoder) convertTrxs(trxs []types.Transaction) ([]core.TransactionModel, error) {
	coreTransactions := make([]core.TransactionModel, 0, len(trxs))
	for _, gethTransaction := range trxs {
		signer := getSigner(gethTransaction)
		sender, err := signer.Sender(&gethTransaction)
		if err != nil {
			return nil, err
		}
		coreTransaction := transToCoreTrans(gethTransaction, &sender)
		coreTransactions = append(coreTransactions, coreTransaction)
	}
	return coreTransactions, nil
}

// DecodeReceipts decodes receipt rlps into type.Receipts
func (snd *SeedNodeDecoder) DecodeReceipts(receiptRlps [][]byte) ([]core.Receipt, []types.Log, error) {
	vulcReceipts := make([]core.Receipt, 0, len(receiptRlps))
	logs := make([]types.Log, 0)
	for _, receiptRlp := range receiptRlps {
		receipt := new(types.ReceiptForStorage)
		err := rlp.DecodeBytes(receiptRlp, receipt)
		if err != nil {
			return nil, nil, err
		}
		for _, log := range receipt.Logs {
			logs = append(logs, *log)
		}
		vulcReceipt, err := common2.ToCoreReceipt(&types.Receipt{
			Bloom:             receipt.Bloom,
			BlockHash:         receipt.BlockHash,
			BlockNumber:       receipt.BlockNumber,
			CumulativeGasUsed: receipt.CumulativeGasUsed,
			ContractAddress:   receipt.ContractAddress,
			Logs:              receipt.Logs,
			Status:            receipt.Status,
			GasUsed:           receipt.GasUsed,
			PostState:         receipt.PostState,
			TransactionIndex:  receipt.TransactionIndex,
			TxHash:            receipt.TxHash,
		})
		if err != nil {
			return nil, nil, err
		}
		vulcReceipts = append(vulcReceipts, vulcReceipt)
	}
	return vulcReceipts, logs, nil
}

// DecodeStateAccounts decodes state rlp into state.Accounts
func (snd *SeedNodeDecoder) DecodeStateAccounts(stateRlps map[common.Hash][]byte) (map[common.Hash]state.Account, error) {
	accounts := make(map[common.Hash]state.Account)
	for hashKey, stateRlp := range stateRlps {
		account := new(state.Account)
		err := rlp.DecodeBytes(stateRlp, account)
		if err != nil {
			return nil, err
		}
		accounts[hashKey] = *account
	}
	return accounts, nil
}
