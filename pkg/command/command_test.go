package command

import (
	"fmt"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/require"
)

func sayCommandFactory() cli.CommandFactory {
	return func() (cli.Command, error) {
		return New(
			"hello",
			"say",
			WithDescription("Say hello"),
			WithArg(0, "who", "who to say hello to"),
			WithFlag("cowboy", "say hello like a cowboy does", false),
			WithAction(func(o *Object) error {
				arg, err := o.Args.StringValueForIndex(0)
				if err != nil {
					return err
				}

				if ok, _ := o.Flags.BoolValue("cowboy"); ok {
					_, err = fmt.Println("howdy", arg)
				} else {
					_, err = fmt.Println("hello", arg)
				}
				return err
			}),
		)
	}
}

func TestHelloSayCommand(t *testing.T) {
	require := require.New(t)

	cmd, err := sayCommandFactory()()
	require.NoError(err)

	// note: do not use shell expansion quotes in flags with spaces as you normally would
	args := []string{"kent", `--cowboy`}
	fmt.Println()
	exitStatus := cmd.Run(args)
	require.Equal(0, exitStatus)
}
