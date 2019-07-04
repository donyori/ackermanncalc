package main

import (
	"math/big"
	"strings"
)

type BigIntIndicator string

func MakeBigIntIndicator(x *big.Int) BigIntIndicator {
	if x == nil {
		return ""
	}
	var builder strings.Builder
	b := x.Bytes()
	builder.Grow(1 + len(b))
	if x.Sign() >= 0 {
		builder.WriteRune('+')
	} else {
		builder.WriteRune('-')
	}
	builder.Write(b)
	return BigIntIndicator(builder.String())
}

func (bii BigIntIndicator) ToBigInt() *big.Int {
	if bii == "" {
		return nil
	}
	x := new(big.Int).SetBytes([]byte(bii)[1:])
	if bii[0] == '-' {
		x.Neg(x)
	}
	return x
}
