## getopt(3) for go

I get what `flag` is doing, and it's pretty, but it's also a lot more
typing than I want in simple cases.

The `getopt(3)` interface is clunky in some ways related to C's
design; it relies on being called repeatedly, and sentinel values,
because allocation is expensive and there's no multiple returns.

But I really like the expressiveness of the option string, for the
common case with single-letter options, and I like the ability to
combine short options.

Usage:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"github.com/seebs/gogetopt"
)

func main() {
	opts, remaining, err := gogetopt.GetOpt(os.Args[1:], "ab:")
	if err != nil {
		log.Fatal("argument parsing failed: %s", err)
	}
	if opts["a"] != nil { // boolean option was set
		fmt.Printf("a: yes.\n")
	}
	if opts["b"] != nil { // string option
		fmt.Printf("b: %s\n", opts["b"].Value)
	}
	if len(remaining) > 0 {
		fmt.Println("additional args:")
		for _, word := range remaining {
			fmt.Printf("  - %s\n", word)
		}
	}
}
```

### Invocation

`opts, remaining, err := gogetopt.GetOpt(args, optstring)`

Takes a slice of strings and an option string, returns a map of provided
options, remaining arguments as a slice of strings, and an error if
an error occurred. If you are using `os.Args`, you probably want
to specify `os.Args[1:]` to avoid the command's name being treated
as a parameter.

```go
type Option {
    Value   string
    Int     int64
    Float   float64
}

type Options map[string]*Option
```

The Options type also provides a `.Seen(s)` method which checks
whether an option is present, because I think that's prettier than
`!= nil`.

### Option Strings

An option string is a list of allowed options. Letters are
allowed options.  A letter followed by punctuation indicates
an option which can take an additional parameter as a value,
or which can be specified more than once:

```
:	Any string
#	Integer number
.	Floating-point number
+	Can be specified multiple times (but takes no argument)
```

Flags which take a value consume additional arguments
starting immediately after the argument containing the flag
itself.  Thus, given `getopt.GetOpt(args, "a:b:c")`:

```
"-ab" "foo" "bar" "-c" "baz"
```

yields:

```
a: foo
b: bar
c: baz
```

Options with `+` specified may be present multiple times, and
in that case, `opts[c].Int` will be the number of times the
option was present. (This supports idioms like `-v` for
verbose output, and `-vv` for more-verbose output.)

The special option `--` indicates the end of option parsing,
as does any non-option encountered.

## Future Directions

I'm considering cool features like allowing additional
parameters to `GetOpt` which are used to specify where to
store values, but in the past I've found that this was
not especially *useful* to me even though I always think
it sounds neat, so I'm not doing it yet.

I sort of dislike the option map being a map of pointers,
but map members being non-addressable makes that seem less
bad than at least some alternatives.
