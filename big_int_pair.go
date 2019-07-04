package main

import (
	"errors"
	"math/big"
)

type BigIntPair struct {
	x, y BigIntIndicator
}

func NewBigIntPair(x, y interface{}) *BigIntPair {
	if x == nil {
		x = Zero
	}
	if y == nil {
		y = Zero
	}
	p := new(BigIntPair)
	switch x.(type) {
	case BigIntIndicator:
		p.x = x.(BigIntIndicator)
	case string:
		p.x = BigIntIndicator(x.(string))
	case *big.Int:
		p.x = MakeBigIntIndicator(x.(*big.Int))
	default:
		panic(errors.New("x is not a BigIntIndicator, string, or *big.Int"))
	}
	switch y.(type) {
	case BigIntIndicator:
		p.y = y.(BigIntIndicator)
	case string:
		p.y = BigIntIndicator(y.(string))
	case *big.Int:
		p.y = MakeBigIntIndicator(y.(*big.Int))
	default:
		panic(errors.New("y is not a BigIntIndicator, string, or *big.Int"))
	}
	return p
}

func (bip *BigIntPair) X() *big.Int {
	if bip == nil {
		return nil
	}
	return IndicatorToBigInt(bip.x)
}

func (bip *BigIntPair) Y() *big.Int {
	if bip == nil {
		return nil
	}
	return IndicatorToBigInt(bip.y)
}

func (bip *BigIntPair) XInd() BigIntIndicator {
	if bip == nil {
		return ""
	}
	return bip.x
}

func (bip *BigIntPair) YInd() BigIntIndicator {
	if bip == nil {
		return ""
	}
	return bip.y
}
