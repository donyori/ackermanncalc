package main

import "math/big"

var (
	NegOne = big.NewInt(-1)
	Zero   = big.NewInt(0)
	One    = big.NewInt(1)
	Two    = big.NewInt(2)
	Three  = big.NewInt(3)
	Four   = big.NewInt(4)
)

var indicatorToBigIntMap map[BigIntIndicator]*big.Int

func init() {
	indicatorToBigIntMap = make(map[BigIntIndicator]*big.Int)
	indicatorToBigIntMap[MakeBigIntIndicator(NegOne)] = NegOne
	indicatorToBigIntMap[MakeBigIntIndicator(Zero)] = Zero
	indicatorToBigIntMap[MakeBigIntIndicator(One)] = One
	indicatorToBigIntMap[MakeBigIntIndicator(Two)] = Two
	indicatorToBigIntMap[MakeBigIntIndicator(Three)] = Three
	indicatorToBigIntMap[MakeBigIntIndicator(Four)] = Four
}

func IndicatorToBigInt(indicator BigIntIndicator) *big.Int {
	x := indicatorToBigIntMap[indicator]
	if x == nil {
		x = indicator.ToBigInt()
	}
	return x
}
