package getopt

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func ExampleNewArgv() {
	scanner, err := NewArgv("a:bz::v", []string{"getopt", "-ba42", "-v", "-z", "--", "-w", "arg1", "arg2"})
	if err != nil {
		panic("error creating scanner: " + err.Error())
	}

	for scanner.Scan() {
		opt, err := scanner.Option()
		if err != nil {
			panic("error: " + err.Error())
		}

		if opt.HasArg() {
			fmt.Printf("%s: got option %q with arg %q\n", scanner.ProgramName(), opt.Opt, opt)
		} else {
			fmt.Printf("%s: got option %q\n", scanner.ProgramName(), opt.Opt)
		}
	}
	fmt.Printf("%s: remaining arguments: %q\n", scanner.ProgramName(), scanner.Args())
	// Output:
	// getopt: got option 'b'
	// getopt: got option 'a' with arg "42"
	// getopt: got option 'v'
	// getopt: got option 'z'
	// getopt: remaining arguments: ["-w" "arg1" "arg2"]
}

func TestOptionsNoArgs(t *testing.T) {
	examples := []struct {
		optstring string
		argv      []string
		expected  []*Option
		errors    []error
	}{
		{
			"",
			[]string{"getopt", "-a", "-b"},
			nil,
			[]error{InvalidOptionError('a')},
		},
		{
			"c",
			[]string{"getopt", "-a", "-b"},
			nil,
			[]error{InvalidOptionError('a')},
		},
		{
			"ab",
			[]string{"getopt", "-a", "-b"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}},
			nil,
		},
		{
			"Ac5",
			[]string{"getopt", "-A", "-c", "-5"},
			[]*Option{{Opt: 'A'}, {Opt: 'c'}, {Opt: '5'}},
			nil,
		},
		{
			"bc",
			[]string{"getopt", "-bc", "-z"},
			[]*Option{{Opt: 'b'}, {Opt: 'c'}},
			[]error{InvalidOptionError('z')},
		},
		{
			"abc",
			[]string{"getopt", "-abc"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c'}},
			nil,
		},
		{
			"abc",
			[]string{"getopt", "-ab", "-c"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c'}},
			nil,
		},
		{
			"abc",
			[]string{"getopt", "-a", "-bc"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c'}},
			nil,
		},
		{
			"abc",
			[]string{"getopt", "-a", "-z"},
			[]*Option{{Opt: 'a'}},
			[]error{InvalidOptionError('z')},
		},
		{
			"abc",
			[]string{"getopt", "-az"},
			[]*Option{{Opt: 'a'}},
			[]error{InvalidOptionError('z')},
		},
	}

	for i, ex := range examples {
		actual, errors, _ := parseOptions(t, ex.optstring, ex.argv)
		if len(errors) > 0 || len(ex.errors) > 0 {
			if len(errors) > 0 && len(ex.errors) == 0 {
				t.Errorf("example %d: expected no errors, got\n%s", i+1, dumpErrors(errors))
			} else if len(errors) == 0 && len(ex.errors) > 0 {
				t.Errorf("example %d: expected errors\n%s\ngot none", i+1, dumpErrors(ex.errors))
			} else {
				expectedErrors := dumpErrors(ex.errors)
				actualErrors := dumpErrors(errors)
				if expectedErrors != actualErrors {
					t.Errorf("example %d: expected errors\n%s\ngot\n%s", i+1, expectedErrors, actualErrors)
				}
			}
		} else {
			if !reflect.DeepEqual(ex.expected, actual) {
				t.Errorf("example %d: expected options\n%s\ngot\n%s", i+1, dumpOptions(ex.expected), dumpOptions(actual))
			}
		}
	}
}

func TestOptionsWithArgs(t *testing.T) {
	examples := []struct {
		optstring string
		argv      []string
		expected  []*Option
		errors    []error
	}{
		{
			"a:b",
			[]string{"getopt", "-a1", "-b"},
			[]*Option{{Opt: 'a', Arg: optArg("1")}, {Opt: 'b'}},
			nil,
		},
		{
			"ab:",
			[]string{"getopt", "-a", "-bfoo"},
			[]*Option{{Opt: 'a'}, {Opt: 'b', Arg: optArg("foo")}},
			nil,
		},
		{
			"ab:",
			[]string{"getopt", "-ab42"},
			[]*Option{{Opt: 'a'}, {Opt: 'b', Arg: optArg("42")}},
			nil,
		},
		// 'a' consumes "b" as its argument
		{
			"a:b",
			[]string{"getopt", "-ab"},
			[]*Option{{Opt: 'a', Arg: optArg("b")}},
			nil,
		},
		// 'a' consumes "-b" as its argument
		{
			"a:b",
			[]string{"getopt", "-a", "-b"},
			[]*Option{{Opt: 'a', Arg: optArg("-b")}},
			nil,
		},
		{
			"ab:",
			[]string{"getopt", "-a", "-b"},
			[]*Option{{Opt: 'a'}},
			[]error{MissingArgumentError('b')},
		},
		// 'a' consumes "-b" as its argument
		{
			"a:b:",
			[]string{"getopt", "-a", "-b"},
			[]*Option{{Opt: 'a', Arg: optArg("-b")}},
			nil,
		},
		{
			"a:b",
			[]string{"getopt", "-abcd"},
			[]*Option{{Opt: 'a', Arg: optArg("bcd")}},
			nil,
		},
		{
			"a:b",
			[]string{"getopt", "-a-bcd"},
			[]*Option{{Opt: 'a', Arg: optArg("-bcd")}},
			nil,
		},
		{
			"ab:",
			[]string{"getopt", "-abcd"},
			[]*Option{{Opt: 'a'}, {Opt: 'b', Arg: optArg("cd")}},
			nil,
		},
		// optional arguments
		{
			":ab:",
			[]string{"getopt", "-ab"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}},
			nil,
		},
		{
			":a:b:",
			[]string{"getopt", "-ab"},
			[]*Option{{Opt: 'a', Arg: optArg("b")}},
			nil,
		},
		{
			":ab:c:",
			[]string{"getopt", "-a", "-bc"},
			[]*Option{{Opt: 'a'}, {Opt: 'b', Arg: optArg("c")}},
			nil,
		},
		{
			":ab:c:",
			[]string{"getopt", "-a", "-b", "c"},
			[]*Option{{Opt: 'a'}, {Opt: 'b', Arg: optArg("c")}},
			nil,
		},
		{
			":ab:c:",
			[]string{"getopt", "-a", "-b", "-c"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c'}},
			nil,
		},
		{
			"ab::",
			[]string{"getopt", "-ab"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}},
			nil,
		},
		{
			"a::bc:",
			[]string{"getopt", "-a", "foo", "-bc42"},
			[]*Option{{Opt: 'a', Arg: optArg("foo")}, {Opt: 'b'}, {Opt: 'c', Arg: optArg("42")}},
			nil,
		},
		{
			"a::bc:",
			[]string{"getopt", "-a", "-bc42"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c', Arg: optArg("42")}},
			nil,
		},
		{
			"a::bc::",
			[]string{"getopt", "-a", "-bc42"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c', Arg: optArg("42")}},
			nil,
		},
		{
			"a::bc::",
			[]string{"getopt", "-a", "-bc"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}, {Opt: 'c'}},
			nil,
		},
	}

	for i, ex := range examples {
		actual, errors, _ := parseOptions(t, ex.optstring, ex.argv)
		if len(errors) > 0 || len(ex.errors) > 0 {
			if len(errors) > 0 && len(ex.errors) == 0 {
				t.Errorf("example %d: expected no errors, got\n%s", i+1, dumpErrors(errors))
			} else if len(errors) == 0 && len(ex.errors) > 0 {
				t.Errorf("example %d: expected errors\n%s\ngot none", i+1, dumpErrors(ex.errors))
			} else {
				expectedErrors := dumpErrors(ex.errors)
				actualErrors := dumpErrors(errors)
				if expectedErrors != actualErrors {
					t.Errorf("example %d: expected errors\n%s\ngot\n%s", i+1, expectedErrors, actualErrors)
				}
			}
		} else {
			if !reflect.DeepEqual(ex.expected, actual) {
				t.Errorf("example %d: expected options\n%s\ngot\n%s", i+1, dumpOptions(ex.expected), dumpOptions(actual))
			}
		}
	}
}

func TestOptionsRemainingArgs(t *testing.T) {
	examples := []struct {
		optstring string
		argv      []string
		expected  []*Option
		remaining []string
		errors    []error
	}{
		{
			"ab",
			[]string{"getopt", "-a", "-b"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}},
			nil,
			nil,
		},
		{
			"ab",
			[]string{"getopt", "-a", "-b", "arg1", "arg2"},
			[]*Option{{Opt: 'a'}, {Opt: 'b'}},
			[]string{"arg1", "arg2"},
			nil,
		},
		{
			"ab",
			[]string{"getopt", "-a", "--", "-b", "--", "arg1", "arg2"},
			[]*Option{{Opt: 'a'}},
			[]string{"-b", "--", "arg1", "arg2"},
			nil,
		},
		{
			":a:b",
			[]string{"getopt", "-a", "--", "-b", "--", "arg1", "arg2"},
			[]*Option{{Opt: 'a'}},
			[]string{"-b", "--", "arg1", "arg2"},
			nil,
		},
		{
			"a::b",
			[]string{"getopt", "-a", "--", "-b", "--", "arg1", "arg2"},
			[]*Option{{Opt: 'a'}},
			[]string{"-b", "--", "arg1", "arg2"},
			nil,
		},
	}

	for i, ex := range examples {
		actual, errors, remaining := parseOptions(t, ex.optstring, ex.argv)
		if len(errors) > 0 || len(ex.errors) > 0 {
			if len(errors) > 0 && len(ex.errors) == 0 {
				t.Errorf("example %d: expected no errors, got\n%s", i+1, dumpErrors(errors))
			} else if len(errors) == 0 && len(ex.errors) > 0 {
				t.Errorf("example %d: expected errors\n%s\ngot none", i+1, dumpErrors(ex.errors))
			} else {
				expectedErrors := dumpErrors(ex.errors)
				actualErrors := dumpErrors(errors)
				if expectedErrors != actualErrors {
					t.Errorf("example %d: expected errors\n%s\ngot\n%s", i+1, expectedErrors, actualErrors)
				}
			}
		} else {
			if !reflect.DeepEqual(ex.expected, actual) {
				t.Errorf("example %d: expected options\n%s\ngot\n%s", i+1, dumpOptions(ex.expected), dumpOptions(actual))
			}
			if !reflect.DeepEqual(ex.remaining, remaining) {
				t.Errorf("example %d: expected options\n%s\ngot\n%s", i+1, dumpRemaining(ex.remaining), dumpRemaining(remaining))
			}
		}
	}
}

func parseOptions(t *testing.T, optstring string, argv []string) ([]*Option, []error, []string) {
	var options []*Option
	var errors []error

	scanner, err := NewArgv(optstring, argv)
	if err != nil {
		t.Fatal(err)
		return nil, nil, nil
	}
	for scanner.Scan() {
		opt, err := scanner.Option()
		if err != nil {
			errors = append(errors, err)
		} else {
			options = append(options, opt)
		}
	}

	return options, errors, scanner.Args()
}

func dumpErrors(errors []error) string {
	res := make([]string, len(errors))
	for i, err := range errors {
		res[i] = fmt.Sprintf("\t%s", err.Error())
	}
	return strings.Join(res, "\n")
}

func dumpOptions(options []*Option) string {
	res := make([]string, len(options))
	for i, opt := range options {
		res[i] = fmt.Sprintf("\t{Opt: %q, Arg: %q}", opt.Opt, opt)
	}
	return strings.Join(res, "\n")
}

func dumpRemaining(args []string) string {
	res := make([]string, len(args))
	for i, arg := range args {
		res[i] = fmt.Sprintf("\t%s", arg)
	}
	return strings.Join(res, "\n")
}
