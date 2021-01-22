package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/picatz/lion/pkg/command"
)

const (
	app        = "hello"
	SayComamnd = "say"
)

func main() {
	c := cli.NewCLI(app, "1.0.0")

	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		SayComamnd: SayCommandFactory(),
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v", err)
	}

	os.Exit(exitStatus)
}

func SayCommandFactory() cli.CommandFactory {
	return func() (cli.Command, error) {
		return command.New(
			app,
			SayComamnd,
			command.WithDescription("Say hello"),
			command.WithArg(0, "who", "who to say hello to"),
			command.WithFlag("cowboy", "say hello like a cowboy does", false),
			command.WithAction(SayCommandAction),
		)
	}
}

func SayCommandAction(c *command.Object) error {
	arg, err := c.Args.StringValueForIndex(0)
	if err != nil {
		return err
	}

	if ok, _ := c.Flags.BoolValue("cowboy"); ok {
		_, err = fmt.Println("howdy", arg)
	} else {
		_, err = fmt.Println("hello", arg)
	}
	return err
}
