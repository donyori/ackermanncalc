package main

import (
	"math/big"
	"testing"
	"time"
)

func TestA42(t *testing.T) {
	a, cost := A(Four, Two, 0, 1, nil)
	t.Logf("Cost: %+v", *cost)
	a42 := big.NewInt(65536)
	a42.Exp(Two, a42, nil)
	a42.Sub(a42, Three)
	if a.Cmp(a42) != 0 {
		t.Error("Result of A(4, 2) is wrong.")
	}
}

func TestATimeout(t *testing.T) {
	var a *big.Int
	var cost *CostInfo
	six := big.NewInt(6)
	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		a, cost = A(six, six, time.Second, 1, nil) // A(6, 6) is big enough.
	}()
	time.Sleep(time.Second + time.Millisecond)
	select {
	case <-doneC:
		t.Log("a:", a)
		t.Logf("Cost: %+v", *cost)
	default:
		t.Fatal("Timeout not work.")
	}
}
