package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-glint"
	"github.com/picatz/lion/pkg/command"
)

const (
	app            = "clock"
	defaultCommand = "" // empty string is the default command
)

func main() {
	c := cli.NewCLI(app, "1.0.0")

	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		defaultCommand: commandFactory(),
	}

	// https://github.com/mitchellh/cli/blob/8c0c01154428e654cfad0336e01de9d24ea306c4/help.go#L17
	//
	// builtin helper doesn't handle cases where the app only has one default command
	c.HelpFunc = func(commands map[string]cli.CommandFactory) string {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(
			"Usage: %s [--version] [--help]",
			app))

		// removed ...

		return buf.String()
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v", err)
	}

	os.Exit(exitStatus)
}

func commandFactory() cli.CommandFactory {
	return func() (cli.Command, error) {
		return command.New(
			app,
			defaultCommand,
			command.WithDescription("Simple clock for your CLI"),
			command.WithAction(commandAction),
		)
	}
}

func commandAction(c *command.Object) error {
	d := glint.New()
	d.Append(
		glint.Style(
			glint.TextFunc(func(rows, cols uint) string {
				return fmt.Sprintf("%v", time.Now())
			}),
			glint.Bold(),
			glint.Color("green"),
		),
	)
	d.SetRefreshRate(100 * time.Millisecond)
	d.Render(context.Background())

	return nil
}
