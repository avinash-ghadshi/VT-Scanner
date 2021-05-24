package gogetopt

import (
	"fmt"
	"strconv"
	"unicode"
)

// Option represents an option identified on the command line.
// For options which take an argument, the original text is stored
// in Value. For numeric options, the parsed numeric value is also
// stored in the corresponding member. For boolean options (the
// default), the struct will have all zero values; the indicator
// of the option's presence is that it's present in the Options
// map.
type Option struct {
	Value string
	Int   int
	Float float64
}

// Options is a map of the options found by GetOpt. An option
// is present only if it was specified in the arguments parsed.
type Options map[string]*Option

// Seen is a convenience function to allow a query for whether
// an option has been set without using the ", ok" format, thus
// allowing a check inline in an expression.
func (o Options) Seen(s string) bool {
	_, ok := o[s]
	return ok
}

type optType int

type opt struct {
	typ optType
}

func (o *opt) setType(typ optType) error {
	if o == nil {
		return fmt.Errorf("option type specifier without option")
	}
	o.typ = typ
	return nil
}

const (
	optBool optType = iota
	optString
	optInt
	optFloat
	optCount
)

type optMap map[string]*opt

// parseOpt converts an option string to a map of options
func parseOpt(optstring string) (optMap, error) {
	o := make(optMap)
	var prev *opt
	for _, r := range optstring {
		c := string(r)
		if o[c] != nil {
			return nil, fmt.Errorf("duplicate option specifiers for '%s'", c)
		}
		if unicode.IsLetter(r) {
			o[c] = &opt{optBool}
			prev = o[c]
			continue
		}
		var err error
		switch r {
		case ':':
			err = prev.setType(optString)
		case '#':
			err = prev.setType(optInt)
		case '.':
			err = prev.setType(optFloat)
		case '+':
			err = prev.setType(optCount)
		default:
			return nil, fmt.Errorf("invalid option specifier '%s'", c)
		}
		if err != nil {
			return nil, err
		}
		// you can only specify a type once per option
		prev = nil
	}
	return o, nil
}

// GetOpt interprets the provided arguments according to optstring,
// which lists individual flags. All flags are single characters
// which are letters or numbers.
func GetOpt(args []string, optstring string) (opts Options, remaining []string, err error) {
	known, err := parseOpt(optstring)
	if err != nil {
		return nil, args, err
	}
	opts = make(Options)
	next := 0
ArgLoop:
	for next < len(args) {
		if args[next] == "--" {
			// skip this, return the rest without looking at them
			next++
			break
		}
		if args[next][0] != '-' {
			break
		}
		// we're now looking at an argument which started with a hyphen
		flags := args[next][1:]
		next++
		for _, f := range flags {
			flag := string(f)
			if known[flag] == nil {
				err = fmt.Errorf("unknown option '%s'", flag)
				break ArgLoop
			}
			opt, ok := opts[flag]
			if !ok {
				opt = &Option{}
				opts[flag] = opt
			} else {
				if known[flag].typ != optCount {
					err = fmt.Errorf("duplicate option '%s'", flag)
					break ArgLoop
				}
			}
			if known[flag].typ == optBool {
				continue
			}
			if known[flag].typ == optCount {
				opt.Int++
				continue
			}
			if next >= len(args) {
				err = fmt.Errorf("option '%s' requires an argument", flag)
				break ArgLoop
			}
			opts[flag].Value = args[next]
			next++
			switch known[flag].typ {
			case optInt:
				opts[flag].Int, err = strconv.Atoi(opts[flag].Value)
			case optFloat:
				opts[flag].Float, err = strconv.ParseFloat(opts[flag].Value, 64)
			}
			if err != nil {
				break ArgLoop
			}
		}
	}
	return opts, args[next:], err
}
