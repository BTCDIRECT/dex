package wrappers

import (
	"database/sql/driver"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// Hash is a wrapper around common.Hash with additional functionalities
type Hash struct {
	common.Hash
}

// Scan implements Scanner.Scan for common.Hash
func (hash *Hash) Scan(value interface{}) error {
	hex := strings.TrimPrefix(strings.TrimPrefix(value.(string), "0X"), "0x")
	hash.SetBytes(common.Hex2Bytes(hex))
	return nil
}

// Value implements driver.Valuer interface
func (hash *Hash) Value() (driver.Value, error) {
	return hash.Hex(), nil
}

// WrapHash wraps common.Hash
func WrapHash(hash *common.Hash) *Hash {
	return &Hash{*hash}
}
