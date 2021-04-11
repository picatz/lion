# lion

🦁 CLI Application Framework for [`mitchellh/cli`](https://github.com/mitchellh/cli)

```go
const sayComamnd = "say"

func sayCommandFactory() cli.CommandFactory {
	return func() (cli.Command, error) {
		return command.New(
			app,
			sayComamnd,
			command.WithDescription("Say hello to someone"),
			command.WithArg(0, "who", "who to say hello to"),
			command.WithFlag("cowboy", "say hello like a cowboy does", false),
			command.WithAction(sayCommandAction),
		)
	}
}

func sayCommandAction(c *command.Object) error {
	arg, err := c.Args.StringValueForIndex(0)
	if err != nil {
		return err
	}

	if ok, _ := c.Flags.BoolValue("cowboy"); ok {
		fmt.Println("howdy", arg)
	} else {
		fmt.Println("hello", arg)
	}
	return nil
}
```

```console
$ hello
Usage: hello [--version] [--help] <command> [<args>]

Available commands are:
    say    Say hello

$ hello say
missing positional argument(s): expected 1, given 0

...
$ hello say kent
hello kent
$ hello say kent --cowboy
howdy kent
```
