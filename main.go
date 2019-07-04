package main

import (
	"flag"
	"fmt"
	"math/big"
	"time"
)

const (
	DefaultTimeout time.Duration = time.Minute
	DefaultVerbose int           = 0
)

func main() {
	var timeout time.Duration
	var verbose int
	timeoutUsage := "Time limit, only valid when it is positive"
	verboseUsage := "> 0 to show debug message"
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
	m, ok := new(big.Int).SetString(args[0], 0)
	if !ok {
		fmt.Println("Argument m cannot be parsed as an integer.")
		return
	}
	n, ok := new(big.Int).SetString(args[1], 0)
	if !ok {
		fmt.Println("Argument n cannot be parsed as an integer.")
		return
	}
	a, cost := A(m, n, timeout, verbose, func(err error, stackTrace string) {
		fmt.Println("Error:", err)
		if stackTrace != "" {
			fmt.Println("Stack trace:")
			fmt.Println(stackTrace)
		}
	})
	if a != nil {
		fmt.Printf("A(%v, %v) = %v\n", m, n, a)
	} else {
		fmt.Printf("Timeout. (time limit: %v)\n", timeout)
	}
	fmt.Printf("Elapsed %v. Max stack size: %v. Number of saved Ackermann–Péter function: %v.\n",
		cost.Time, cost.MaxStackSize, cost.NumSaved)
}
