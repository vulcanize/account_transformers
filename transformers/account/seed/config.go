package seed

import (
	"github.com/vulcanize/account_transformers/transformers/account/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/config"
)

// Default config subscribes to all canonical headers, all transactions, receipts with value transfer type events
var DefaultConfig = config.Subscription{
	BackFill: true,
	TrxFilter: config.TrxFilter{
		Src: constants.StrAccountAddresses(),
		Dst: constants.StrAccountAddresses(),
	},
	HeaderFilter: config.HeaderFilter{
		FinalOnly: true,
	},
	ReceiptFilter: config.ReceiptFilter{
		Topic0s:   constants.StrTopic0s,
		Contracts: []string{},
	},
	StateFilter: config.StateFilter{
		Addresses: constants.StrAccountAddresses(),
	},
	StorageFilter: config.StorageFilter{
		Off: true,
	},
}
