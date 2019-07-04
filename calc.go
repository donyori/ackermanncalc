package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/donyori/gorecover"
)

type CostInfo struct {
	Time         time.Duration
	MaxStackSize *big.Int
	NumSaved     *big.Int
}

type stackElement struct {
	target         *BigIntPair
	numArrowCopies *big.Int
}

func A(m, n *big.Int, timeout time.Duration, verbose int,
	panicHandler func(err error, stackTrace string)) (
	a *big.Int, cost *CostInfo) {
	if m.Sign() < 0 {
		panic(fmt.Errorf("m (%v) is negative", m))
	}
	if n.Sign() < 0 {
		panic(fmt.Errorf("n (%v) is negative", n))
	}
	startTime := time.Now()
	var timeoutC <-chan time.Time
	if timeout > 0 {
		timeoutC = time.After(timeout)
	}
	doneC := make(chan struct{})
	cost = &CostInfo{MaxStackSize: big.NewInt(0)}
	savedMap := NewBigMap(0)
	f := func() {
		defer close(doneC)
		stack := NewBigStack()
		stack.Push(&stackElement{
			target: NewBigIntPair(Sub(m, Two), Add(n, Three)),
		})
		var v interface{}
		var top *stackElement
		var sLen, x, b, r *big.Int
		didPush := false
		for !stack.IsEmpty() {
			v, _ = stack.Top()
			top = v.(*stackElement)
			if didPush {
				didPush = false
				sLen = stack.Len()
				if cost.MaxStackSize.Cmp(sLen) < 0 {
					cost.MaxStackSize.Set(sLen)
				}
				if verbose > 0 {
					fmt.Printf("Pushed A(%v, %v).\n",
						Add(top.target.X(), Two),
						Sub(top.target.Y(), Three))
				}
			}
			if top.numArrowCopies != nil {
				// r = 2 ↑(n_top_target-1) b0
				if top.numArrowCopies.Cmp(One) > 0 {
					// r_top_target = [2 ↑(n_top_target-1)](top.numArrowCopies-1 copies) r
					// So decrease top.numArrowCopies and push 2 ↑(n_top_target-1) r.
					top.numArrowCopies.Sub(top.numArrowCopies, One)
					stack.Push(&stackElement{
						target: NewBigIntPair(Sub(top.target.X(), One), r),
					})
					didPush = true
				} else {
					// r is just the result of the top target.
					// Save the result and pop.
					saveKnuthsUpArrowNotationValue(savedMap, top.target, r)
					stack.Pop()
					if verbose > 0 {
						fmt.Printf("Saved A(%v, %v).\n",
							Add(top.target.X(), Two),
							Sub(top.target.Y(), Three))
						fmt.Println("Popped.")
					}
				}
			} else {
				x = simpleKnuthsUpArrowNotationValue(top.target)
				if x == nil {
					b, x = knownKnuthsUpArrowNotationValue(
						savedMap, top.target.XInd())
					top.numArrowCopies = top.target.YInd().ToBigInt() // b of the target, must > b.
					top.numArrowCopies.Sub(top.numArrowCopies, b)     // Total (bt - b) arrows.
					stack.Push(&stackElement{
						target: NewBigIntPair(Sub(top.target.X(), One), x),
					})
					didPush = true
				} else {
					r = x
					stack.Pop()
					if verbose > 0 {
						fmt.Println("Popped.")
					}
				}
			}
		}
		if r != nil {
			a = Sub(r, Three)
		}
	}
	go func() {
		if panicHandler != nil {
			err, st := gorecover.RecoverWithStackTrace(f)
			if err != nil {
				panicHandler(err, st)
			}
		} else {
			f()
		}
	}()
	select {
	case <-doneC:
	case <-timeoutC:
	}
	cost.Time = time.Since(startTime)
	cost.NumSaved = savedMap.Len()
	return
}

func Add(x, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

func Sub(x, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

func Mul(x, y *big.Int) *big.Int {
	return new(big.Int).Mul(x, y)
}

func Exp(a, x *big.Int) *big.Int {
	return new(big.Int).Exp(a, x, nil)
}

func simpleKnuthsUpArrowNotationValue(nAndB *BigIntPair) (x *big.Int) {
	n := nAndB.X()
	b := nAndB.Y()
	cmp := n.Cmp(One)
	if cmp > 0 {
		switch b {
		case Two:
			x = Four
		case One:
			x = Two
		case Zero:
			x = One
		}
	} else if cmp == 0 {
		x = Exp(Two, b)
	} else if n == Zero {
		x = Mul(Two, b)
	} else if n == NegOne {
		x = Add(Two, b)
	} else { // n == -2
		x = Add(b, One)
	}
	return
}

func knownKnuthsUpArrowNotationValue(m *BigMap, nInd BigIntIndicator) (
	b, x *big.Int) {
	v, ok := m.Load(string(nInd))
	if !ok {
		return Two, Four
	}
	bx := v.(*BigIntPair)
	return bx.X(), bx.Y()
}

func saveKnuthsUpArrowNotationValue(m *BigMap, nAndB *BigIntPair, x *big.Int) {
	m.Store(string(nAndB.XInd()), NewBigIntPair(nAndB.YInd(), x))
}
