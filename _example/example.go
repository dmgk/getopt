package main

import (
	"fmt"

	"github.com/dmgk/getopt"
)

// go run example.go -ba42 -v -z1 -x arg1 arg2
func main() {
	scanner, err := getopt.New("a:bvz:")
	if err != nil {
		fmt.Printf("error creating scanner: %s\n", err)
		return
	}

	for scanner.Scan() {
		opt, err := scanner.Option()
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}

		if opt.HasArg() {
			fmt.Printf("got option %q with arg %q\n", opt.Opt, opt)
		} else {
			fmt.Printf("got option %q\n", opt.Opt)
		}
	}

	fmt.Printf("remaining arguments: %v\n", scanner.Args())
}
