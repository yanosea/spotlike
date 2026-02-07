package proxy

import (
	"github.com/thanhpk/randstr"
)

// Randstr is an interface that provides a proxy of the methods of randstr.
type Randstr interface {
	Hex(length int) string
}

type randstrProxy struct{}

// NewRandstr returns a new instance of the Randstr interface.
func NewRandstr() Randstr {
	return &randstrProxy{}
}

// Hex is a proxy method that calls the Hex method of the randstr.
func (*randstrProxy) Hex(length int) string {
	return randstr.Hex(length)
}
