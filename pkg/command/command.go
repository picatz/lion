package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/cli"
)

type Object struct {
	cli.Command
	App      string
	Name     string
	Desc     string
	Args     Args
	Flags    Flags
	Examples Examples
	Action   Action
	Context  context.Context
}

type Examples []string

func (c Examples) HelpString(appName, cmdName string) string {
	all := []string{}
	for _, ex := range c {
		all = append(all, fmt.Sprintf("\t$ %s %s %s", appName, cmdName, ex))
	}
	return strings.TrimRight(strings.Join(all, "\n"), "\n")
}

type Action = func(*Object) error

type Option = func(*Object) error

func New(app, name string, opts ...Option) (*Object, error) {
	c := &Object{Name: name, Flags: Flags{}, Args: Args{}, Examples: Examples{}, Context: context.Background()}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func WithFlag(name, desc string, dflt interface{}) Option {
	return func(c *Object) error {
		dlftStr := fmt.Sprintf("%v", dflt)
		if dlftStr == "" {
			dlftStr = `""`
		}
		c.Flags[name] = &Flag{Name: name, Desc: fmt.Sprintf("%s (default: %s)", desc, dlftStr), Default: dflt}
		return nil
	}
}

func WithArg(index int, name, desc string) Option {
	return func(c *Object) error {
		c.Args[name] = &Arg{Name: name, Index: index, Desc: desc}
		return nil
	}
}

func WithDescription(desc string) Option {
	return func(c *Object) error {
		c.Desc = desc
		return nil
	}
}

func WithExample(example string) Option {
	return func(c *Object) error {
		c.Examples = append(c.Examples, example)
		return nil
	}
}

func WithAction(action Action) Option {
	return func(c *Object) error {
		c.Action = action
		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Object) error {
		c.Context = ctx
		return nil
	}
}

func (c *Object) Help() string {
	var lines = []string{}

	var addSection = func(name, content string) {
		lines = append(lines, fmt.Sprintf("%s:", name))
		lines = append(lines, content)
	}

	var usage string

	if len(c.Args) == 0 {
		usage = fmt.Sprintf("%s %s [options]", c.App, c.Name)
	} else {
		usage = fmt.Sprintf("%s %s %s [options]", c.App, c.Name, c.Args.HelpString())
	}

	addSection("Usage", "\t"+usage)

	if c.Desc != "" {
		addSection("Description", "\t"+c.Desc)
	}

	// if len(c.Flags) != 0 {
	addSection("Flags", c.Flags.HelpString())
	//}

	if len(c.Examples) != 0 {
		addSection("Examples", c.Examples.HelpString(c.App, c.Name))
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func (c *Object) Synopsis() string {
	return c.Desc
}

func (c *Object) newFlagSet() *flag.FlagSet {
	cmdFlags := flag.NewFlagSet(strings.Join(strings.Split(c.Name, " "), "-"), flag.ExitOnError)

	for _, f := range c.Flags {
		switch df := f.Default.(type) {
		case string:
			f.Value = cmdFlags.String(f.Name, df, f.Desc)
		case bool:
			f.Value = cmdFlags.Bool(f.Name, df, f.Desc)
		case int:
			f.Value = cmdFlags.Int(f.Name, df, f.Desc)
		case int64:
			f.Value = cmdFlags.Int64(f.Name, df, f.Desc)
		case time.Duration:
			f.Value = cmdFlags.Duration(f.Name, df, f.Desc)
		}
	}

	return cmdFlags
}

func (c *Object) Run(args []string) int {
	cmdFlags := c.newFlagSet()

	if len(args) < len(c.Args) {
		fmt.Fprintf(os.Stderr, "missing positional argument(s): expected %d, given %d\n", len(c.Args), len(args))
		fmt.Fprintf(os.Stderr, "\n%s\n", c.Help())
		return 1
	}

	if len(c.Args) > 0 {
		// set arguments, if there are enough
		for _, arg := range c.Args {
			if len(args) < arg.Index {
				fmt.Fprintf(os.Stderr, "missing argument %q at index %d\n", arg.Name, arg.Index)
				return 1
			}
			// prevent flags given as indexes
			if strings.HasPrefix(args[arg.Index], "-") {
				fmt.Fprintf(os.Stderr, "unexpected flag %q given at argument index %d, not %q argument\n", args[arg.Index], arg.Index, arg.Name)
				return 1
			}
			arg.Value = args[arg.Index]
		}

		if len(args) < len(c.Args) {
			missingArgIndex := len(args) + 1
			arg, err := c.Args.ForIndex(missingArgIndex)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
			fmt.Fprintf(os.Stderr, "missing argument %q at index %d\n", arg.Name, arg.Index)
		}
	}

	if len(args) >= len(c.Args)+1 {
		err := cmdFlags.Parse(args[len(c.Args):])
		if err != nil {
			fmt.Println(err)
			return 1
		}
	}

	if c.Action != nil {
		err := c.Action(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	} else {
		fmt.Println("command not implemented")
	}

	return 0
}
