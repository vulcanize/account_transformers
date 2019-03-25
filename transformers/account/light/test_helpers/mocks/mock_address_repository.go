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

package mocks

import (
	"github.com/ethereum/go-ethereum/common"
)

type MockAddressRepository struct {
	Addresses       []common.Address
	GetAddressesErr error
	AddAddressErr   error
	PassedAddresses []common.Address
}

func (ar *MockAddressRepository) GetAddresses() ([]common.Address, error) {
	return ar.Addresses, ar.GetAddressesErr
}

func (ar *MockAddressRepository) AddAddress(addr common.Address) error {
	ar.PassedAddresses = append(ar.PassedAddresses, addr)
	return ar.AddAddressErr
}
