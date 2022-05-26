// Package getopt provides a minimal, getopt(3)-like argument parsing implementation
// with POSIX compatible semantics.
package getopt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// InvalidOptionError is returned when scanner encounters an option not listed in optstring.
type InvalidOptionError rune

func (e InvalidOptionError) Error() string {
	return fmt.Sprintf("invalid option: %q", rune(e))
}

// MissingArgumentError is returned when option is missing a required argument.
type MissingArgumentError rune

func (e MissingArgumentError) Error() string {
	return fmt.Sprintf("option requires an argument: %q", rune(e))
}

// Option contains option name and optional argument value.
type Option struct {
	// Option name
	Opt rune
	// Option argument, if any
	Arg *string
}

func (o *Option) HasArg() bool {
	return o.Arg != nil
}

func (o *Option) String() string {
	if o.Arg != nil {
		return *o.Arg
	}
	return ""
}

func (o *Option) Int() (int, error) {
	v, err := strconv.ParseInt(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (o *Option) Int32() (int32, error) {
	v, err := strconv.ParseInt(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (o *Option) Int64() (int64, error) {
	v, err := strconv.ParseInt(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (o *Option) Uint() (uint, error) {
	v, err := strconv.ParseUint(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (o *Option) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint32(v), nil
}

func (o *Option) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(o.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func (o *Option) Float32() (float32, error) {
	v, err := strconv.ParseFloat(o.String(), 64)
	if err != nil {
		return 0, err
	}
	return float32(v), nil
}

func (o *Option) Float64() (float64, error) {
	v, err := strconv.ParseFloat(o.String(), 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Scanner contains option scanner data.
type Scanner struct {
	// Command line arguments
	argv []string
	// Accepted option characters
	optstring string
	// Current argv index
	optind int
	// Current argv element
	arg string
	// Current index in the arg
	optpos int
	// Current option
	optopt rune
	// Last error, if any
	err error
}

// New returns a new options scanner using os.Args as the command line arguments source.
// The option string optstring may contain the following elements:
// individual characters, and characters followed by a colon to indicate an
// option argument is to follow.
// If optstring starts with ':' then all option argument are treated as optional.
func New(optstring string) (*Scanner, error) {
	return NewArgv(optstring, os.Args)
}

// New returns a new options scanner using passed argv as the command line argument source.
// The option string optstring may contain the following elements:
// individual characters, and characters followed by a colon to indicate an
// option argument is to follow.
// If optstring starts with ':' then all option argument are treated as optional.
func NewArgv(optstring string, argv []string) (*Scanner, error) {
	for _, c := range optstring {
		if !isOptionChar(byte(c)) && c != ':' {
			return nil, fmt.Errorf("invalid optstring character: %q", c)
		}
	}
	return &Scanner{
		argv:      argv,
		optstring: optstring,
		optind:    1,
		optpos:    1,
	}, nil
}

// Scan advances options scanner to the next option.
// It returns false when there are no more options or parsing is cancelled by "--".
func (s *Scanner) Scan() bool {
	if s.optind == len(s.argv) || s.err != nil {
		return false
	}

	s.arg = s.argv[s.optind]
	if s.arg == "--" {
		s.optind += 1
		return false
	}
	if len(s.arg) < 2 || s.arg[0] != '-' || !isOptionChar(s.arg[1]) {
		return false
	}

	return true
}

// Option returns the next option or an error when it encounters an unknown option or
// an option that is missing a required argument.
// If optstring starts with ':' then all arguments are treated as optional and missing
// arguments do not cause errors.
func (s *Scanner) Option() (*Option, error) {
	s.optopt = rune(s.arg[s.optpos])

	idx := strings.IndexRune(s.optstring, s.optopt)
	if idx < 0 {
		s.err = InvalidOptionError(s.optopt)
		return nil, s.err
	}

	if idx < len(s.optstring)-1 && s.optstring[idx+1] == ':' {
		// option with an argument
		if len(s.arg) > s.optpos+1 {
			// option and argument are in the same argv element
			res := &Option{
				Opt: s.optopt,
				Arg: optArg(s.arg[s.optpos+1:]),
			}
			s.optind += 1
			s.optpos = 1
			return res, nil
		} else if s.optind+1 < len(s.argv) {
			// option argument is in the next argv element
			res := &Option{
				Opt: s.optopt,
				Arg: optArg(s.argv[s.optind+1]),
			}
			s.optind += 2
			s.optpos = 1
			return res, nil
		} else {
			if s.optstring != "" && s.optstring[0] == ':' {
				// optstring starts with ':', option argument is optional
				s.optind += 1
				s.optpos = 1
				return &Option{
					Opt: s.optopt,
				}, nil
			} else {
				s.err = MissingArgumentError(s.optopt)
				return nil, s.err
			}
		}
	} else {
		// no-argument option
		s.optpos += 1
		if len(s.arg) == s.optpos {
			// arg has the only one option
			s.optind += 1
			s.optpos = 1
		}
		return &Option{
			Opt: s.optopt,
		}, nil
	}
}

// Args returns remaining command line arguments.
func (s *Scanner) Args() []string {
	if s.optind < len(s.argv) {
		return s.argv[s.optind:]
	}
	return nil
}

func isOptionChar(c byte) bool {
	return ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') || ('0' <= c && c <= '9')
}

func optArg(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}
