package bigmath

import (
	"math/big"
	"sync"
)

var (
	// Default precision of float64.
	floatPrec uint = 53
	once      sync.Once
)

// SetFloatPrec sets the app-wide precision for big.Float.
func SetFloatPrec(prec uint) {
	once.Do(func() {
		floatPrec = prec
	})
}

// BlankFloat returns a new big.Float with value 0 and app-wide precision set via SetFloatPrec.
func BlankFloat() *big.Float { return big.NewFloat(0).SetPrec(floatPrec) }

// NewFloat returns a new big.Float with value x and app-wide precision set via SetFloatPrec.
func NewFloat(x float64) *big.Float { return big.NewFloat(x).SetPrec(floatPrec) }

// BlankRat returns a new big.Rat with value 0.
func BlankRat() *big.Rat { return big.NewRat(0, 1) }

// NewRat returns a new big.Rat with value n/d.
func NewRat(n, d int) *big.Rat { return big.NewRat(int64(n), int64(d)) }

// BlankInt returns a new big.Int with value 0.
func BlankInt() *big.Int { return big.NewInt(0) }

// NewInt returns a new big.Int with value x.
func NewInt(x int) *big.Int { return big.NewInt(int64(x)) }
