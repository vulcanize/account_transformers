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

package repositories

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type AddressRepository interface {
	GetAddresses() ([]common.Address, error)
	AddAddress(addr string) error
}

type addressRepository struct {
	DB *postgres.DB
}

func NewAccountHeaderRepository(db *postgres.DB) *addressRepository {
	return &addressRepository{
		DB: db,
	}
}

func (ar *addressRepository) GetAddresses() ([]common.Address, error) {
	dest := make([][]byte, 0)
	err := ar.DB.Select(dest, `SELECT * FROM accounts.addresses`)
	if err != nil {
		return nil, err
	}
	addresses := make([]common.Address, 0, len(dest))
	for _, addrBytes := range dest {
		addr := common.BytesToAddress(addrBytes)
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (ar *addressRepository) AddAddress(addr string) error {
	_, err := ar.DB.Exec(`INSTER INTO accounts.addresses (address) VALUES ($1)`, addr)
	return err
}
