package bigmath

import "math/big"

// ReIm represents a complex number represented by real and imaginary parts.
type ReIm struct {
	Re *big.Float
	Im *big.Float
}

// NewReIm creates a new complex number with real and imaginary parts, the pointers of re and im are stored.
func NewReIm(re, im *big.Float) *ReIm {
	return &ReIm{Re: re, Im: im}
}

// ModArg represents a complex number represented by modulus and argument.
// There is no guarantee that Mod be non-negative, but the real and imaginary part for the underlying complex number always equal to Mod.cos(arg) and Mod.sin(arg).
type ModArg struct {
	Mod *big.Float
	Arg *big.Float
}

// NewModArg creates a new complex number with modulus and argument, the pointers of mod and arg are stored.
func NewModArg(mod, arg *big.Float) *ModArg {
	return &ModArg{Mod: mod, Arg: arg}
}
