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

package utils

type ValueType int

const (
	Uint256 ValueType = iota
	Uint48
	Bytes32
	Address
)

type Key string

type StorageValueMetadata struct {
	Name string
	Keys map[Key]string
	Type ValueType
}

func GetStorageValueMetadata(name string, keys map[Key]string, t ValueType) StorageValueMetadata {
	return StorageValueMetadata{
		Name: name,
		Keys: keys,
		Type: t,
	}
}
