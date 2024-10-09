package randstrproxy

import (
	"github.com/thanhpk/randstr"
)

// Randstr is an interface for randstr.
type Randstr interface {
	Hex(length int) string
}

// RandstrProxy is a struct that implements Randstr.
type RandstrProxy struct{}

// New is a constructor for RandstrProxy.
func New() Randstr {
	return &RandstrProxy{}
}

// Hex is a proxy for randstr.Hex.
func (*RandstrProxy) Hex(length int) string {
	return randstr.Hex(length)
}
