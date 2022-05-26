# Minimal, POSIX compatible argument parsing in Go

Package getopt provides a minimal, getopt(3)-like argument parsing implementation with POSIX compatible semantics.

## Example

```go
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
```

```
$ go run example.go -ba42 -v -z1 -x arg1 arg2
got option 'b'
got option 'a' with arg "42"
got option 'v'
got option 'z' with arg "1"
error: invalid option: 'x'
remaining arguments: [-x arg1 arg2]
```
