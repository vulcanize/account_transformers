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

package repositories

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type AddressRepository interface {
	GetAddresses() ([]common.Address, error)
	AddAddress(addr common.Address) error
}

type addressRepository struct {
	DB *postgres.DB
}

func NewAddressRepository(db *postgres.DB) *addressRepository {
	return &addressRepository{
		DB: db,
	}
}

func (ar *addressRepository) GetAddresses() ([]common.Address, error) {
	dest := new([]string)
	err := ar.DB.Select(dest, `SELECT * FROM accounts.addresses`)
	if err != nil {
		return nil, err
	}
	addresses := make([]common.Address, 0, len(*dest))
	for _, addrStrings := range *dest {
		addr := common.HexToAddress(addrStrings)
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (ar *addressRepository) AddAddress(addr common.Address) error {
	_, err := ar.DB.Exec(`INSERT INTO accounts.addresses (address) VALUES ($1) ON CONFLICT (address) DO NOTHING`, addr.String())
	return err
}
