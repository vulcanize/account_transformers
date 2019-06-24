package constants

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var initialized = false

func initConfig() {
	if initialized {
		return
	}

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Sprintf("Could not find environment file: %v", err))
	}
	initialized = true
}

func getStringMapStringSlice(key string) map[string][]string {
	initConfig()
	value := viper.GetStringMapStringSlice(key)
	if value == nil {
		panic(fmt.Sprintf("No environment configuration variable set for key: \"%v\"", key))
	}
	return value
}

func getStringSlice(key string) []string {
	initConfig()
	value := viper.GetStringSlice(key)
	if value == nil {
		panic(fmt.Sprintf("No environment configuration variable set for key: \"%v\"", key))
	}
	return value
}

func getInt(key string) int64 {
	initConfig()
	return viper.GetInt64(key)
}

// Gets a map of top level Token addresses to a list of other addresses that emit events for this Token
func equivalentTokenAddressesMapping() map[string][]string {
	return getStringMapStringSlice("token.equivalents")
}

func EquivalentTokenAddressesMapping() map[common.Address][]common.Address {
	m := equivalentTokenAddressesMapping()
	addrMap := make(map[common.Address][]common.Address)
	for topAddrStr, equivalentsArray := range m {
		topAddr := common.HexToAddress(topAddrStr)
		addrMap[topAddr] = make([]common.Address, 0, len(equivalentsArray))
		for _, equivalentAddr := range equivalentsArray {
			addrMap[topAddr] = append(addrMap[topAddr], common.HexToAddress(equivalentAddr))
		}
	}
	return addrMap
}

func tokenAddresses() []string {
	return getStringSlice("token.addresses")
}

func TokenAddresses() []common.Address {
	strAddrs := tokenAddresses()
	addrs := make([]common.Address, 0, len(strAddrs))
	for _, strAddr := range strAddrs {
		addrs = append(addrs, common.HexToAddress(strAddr))
	}
	return addrs
}

func accountAddresses() []string {
	return getStringSlice("account.addresses")
}

func AccountAddresses() []common.Address {
	strAddrs := accountAddresses()
	addrs := make([]common.Address, 0, len(strAddrs))
	for _, strAddr := range strAddrs {
		addrs = append(addrs, common.HexToAddress(strAddr))
	}
	return addrs
}

func StartingBlock() int64 {
	return getInt("account.start")
}
