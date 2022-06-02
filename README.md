## getopt

Package getopt provides a minimal, getopt(3)-like argument parsing implementation with POSIX compatible semantics.

[![Go Reference](https://pkg.go.dev/badge/github.com/dmgk/getopt.svg)](https://pkg.go.dev/github.com/dmgk/getopt)
![Tests](https://github.com/dmgk/getopt/actions/workflows/tests.yml/badge.svg)

#### Example

```go
package main

import (
	"fmt"

	"github.com/dmgk/getopt"
)

// go run example.go -ba42 -v -z -- -w arg1 arg2
func main() {
	// -a requires an argument
	// -b and -v have no arguments
	// -z may have an optional argument
	opts, err := getopt.New("a:bz::v")
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
```

```
$ go run example.go -ba42 -v -z -- -w arg1 arg2
example: got option 'b'
example: got option 'a' with arg "42"
example: got option 'v'
example: got option 'z'
example: remaining arguments: [-w arg1 arg2]
```
