package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-glint"
	"github.com/mitchellh/go-glint/components"
	"github.com/picatz/lion/pkg/command"
)

const (
	app         = "copy"
	fileCommand = "file"
)

func main() {
	c := cli.NewCLI(app, "1.0.0")

	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		fileCommand: commandFactory(),
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
			fileCommand,
			command.WithDescription("Simple copy command"),
			command.WithArg(0, "source", "source file path"),
			command.WithArg(1, "destination", "destination file path"),
			command.WithAction(commandAction),
		)
	}
}

type readMonitor struct {
	io.Reader
	progressBar *components.ProgressElement
}

func (cm *readMonitor) Read(p []byte) (int, error) {
	n, err := cm.Reader.Read(p)
	if err != nil {
		return n, err
	}
	cm.progressBar.Add(n)
	return n, nil
}

func commandAction(c *command.Object) error {
	src, err := c.Args.StringValueForIndex(0)
	if err != nil {
		return err
	}

	dst, err := c.Args.StringValueForIndex(1)
	if err != nil {
		return err
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	srcFh, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFh.Close()

	dstFh, err := os.OpenFile(dst, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.FileMode(0655))
	if err != nil {
		return err
	}
	defer dstFh.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	d := glint.New()

	progressBar := components.Progress(int(srcStat.Size()))

	d.Append(
		glint.Style(
			progressBar,
			glint.Bold(),
			glint.Color("green"),
		),
	)

	d.RenderFrame()

	go func() {
		defer cancel()
		defer progressBar.Finish()
		_, err := io.Copy(dstFh, &readMonitor{Reader: srcFh, progressBar: progressBar})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}()

	d.Render(ctx)

	return nil
}
