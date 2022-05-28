package main

import (
	"fmt"

	"github.com/dmgk/getopt"
)

// go run example.go -ba42 -v -z1 -x arg1 arg2
func main() {
	opts, err := getopt.New("a:bvz:")
	if err != nil {
		fmt.Printf("error creating scanner: %s\n", err)
		return
	}

	for opts.Scan() {
		opt, err := opts.Option()
		if err != nil {
			fmt.Printf("%s: error parsing option: %s\n", opts.ProgramName(), err)
			continue
		}

		if opt.HasArg() {
			fmt.Printf("%s: got option %q with arg %q\n", opts.ProgramName(), opt.Opt, opt)
		} else {
			fmt.Printf("%s: got option %q\n", opts.ProgramName(), opt.Opt)
		}
	}

	fmt.Printf("%s: remaining arguments: %v\n", opts.ProgramName(), opts.Args())
}
