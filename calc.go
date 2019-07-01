package main

import (
	"fmt"
	"math"
	"math/big"
	"time"
)

type MNPair struct {
	M int
	N string
}

type CostInfo struct {
	Time         time.Duration
	MaxStackSize int
	NumCalcA     int
}

const StackSizeWarnMinLog float64 = 5.

var One *big.Int = big.NewInt(1)

func A(m int, n *big.Int, verbose int) (a *big.Int, cost *CostInfo) {
	if m < 0 {
		panic(fmt.Errorf("m (%d) is negative", m))
	}
	if n.Sign() < 0 {
		panic(fmt.Errorf("n (%v) is negative", n))
	}
	startTime := time.Now()
	cost = new(CostInfo)
	stack := []MNPair{MNPair{M: m, N: string(n.Bytes())}}
	knownA := make(map[MNPair]*big.Int)
	var top, tmpPair MNPair
	var topN, b *big.Int
	for len(stack) > 0 {
		top = stack[len(stack)-1]
		topN = new(big.Int).SetBytes([]byte(top.N))
		if top.M > 0 {
			if topN.Sign() > 0 {
				tmpPair = MNPair{
					M: top.M,
					N: string(new(big.Int).Sub(topN, One).Bytes()),
				}
				b = knownA[tmpPair]
				if b == nil {
					stack = append(stack, tmpPair)
					if len(stack) > cost.MaxStackSize {
						cost.MaxStackSize = len(stack)
						if verbose > 0 {
							log := math.Log10(float64(cost.MaxStackSize))
							if log == math.Round(log) &&
								log >= StackSizeWarnMinLog {
								fmt.Println("WARNING: stack size up to",
									cost.MaxStackSize)
							}
						}
					}
				} else {
					tmpPair = MNPair{M: top.M - 1, N: string(b.Bytes())}
					b = knownA[tmpPair]
					if b == nil {
						stack = append(stack, tmpPair)
						if len(stack) > cost.MaxStackSize {
							cost.MaxStackSize = len(stack)
							if verbose > 0 {
								log := math.Log10(float64(cost.MaxStackSize))
								if log == math.Round(log) &&
									log >= StackSizeWarnMinLog {
									fmt.Println("WARNING: stack size up to",
										cost.MaxStackSize)
								}
							}
						}
					} else {
						knownA[top] = b
						stack = stack[:len(stack)-1]
					}
				}
			} else {
				tmpPair = MNPair{M: top.M - 1, N: string(One.Bytes())}
				b = knownA[tmpPair]
				if b == nil {
					stack = append(stack, tmpPair)
					if len(stack) > cost.MaxStackSize {
						cost.MaxStackSize = len(stack)
						if verbose > 0 {
							log := math.Log10(float64(cost.MaxStackSize))
							if log == math.Round(log) &&
								log >= StackSizeWarnMinLog {
								fmt.Println("WARNING: stack size up to",
									cost.MaxStackSize)
							}
						}
					}
				} else {
					knownA[top] = b
					stack = stack[:len(stack)-1]
				}
			}
		} else {
			knownA[top] = new(big.Int).Add(topN, One)
			stack = stack[:len(stack)-1]
		}
	}
	cost.Time = time.Since(startTime)
	cost.NumCalcA = len(knownA)
	a = knownA[MNPair{M: m, N: string(n.Bytes())}]
	return
}
