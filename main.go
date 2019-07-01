package main

import (
	"flag"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/donyori/gorecover"
)

type Result struct {
	A    *big.Int
	Cost *CostInfo
}

const (
	DefaultTimeout time.Duration = time.Minute
	DefaultVerbose int           = 1
)

func main() {
	var timeout time.Duration
	var verbose int
	timeoutUsage := "Time limit, only valid when it is positive"
	verboseUsage := "> 0 to show warning message"
	flag.DurationVar(&timeout, "timeout", DefaultTimeout, timeoutUsage)
	flag.DurationVar(&timeout, "t", DefaultTimeout, timeoutUsage+" (shorthand)")
	flag.IntVar(&verbose, "verbose", DefaultVerbose, verboseUsage)
	flag.IntVar(&verbose, "v", DefaultVerbose, verboseUsage+" (shorthand)")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("You should give exactly two arguments: m and n.")
		return
	}
	m, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Argument m cannot be parsed as an integer or overflows. Error message:", err)
		return
	}
	n, ok := new(big.Int).SetString(args[1], 0)
	if !ok {
		fmt.Println("Argument n cannot be parsed as an integer.")
		return
	}
	var timeoutC <-chan time.Time
	resultC := make(chan *Result)
	if timeout > 0 {
		timeoutC = time.After(timeout)
	}
	go func() {
		var r *Result
		err := gorecover.Recover(func() {
			a, cost := A(m, n, verbose)
			r = &Result{
				A:    a,
				Cost: cost,
			}
		})
		if err != nil {
			fmt.Println("Error:", err)
			r = nil
		}
		resultC <- r
	}()
	select {
	case r := <-resultC:
		if r != nil {
			fmt.Printf("A(%d, %v) = %v\n", m, n, r.A)
		}
		fmt.Printf("Elapsed %v. Max stack size: %d. Number of calculated Ackermann–Péter function: %d.\n",
			r.Cost.Time, r.Cost.MaxStackSize, r.Cost.NumCalcA)
	case <-timeoutC:
		fmt.Printf("Timeout. (time limit: %v)\n", timeout)
	}
}
