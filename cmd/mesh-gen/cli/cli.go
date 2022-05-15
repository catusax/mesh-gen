package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	// DefaultCLI is the default, unmodified root command.
	DefaultCLI CLI = NewCLI()

	name        string = "mesh-gen"
	description string = "The Go Micro CLI tool"
	version     string = "latest"
)

// CLI is the interface that wraps the cli app.
//
// CLI embeds the Cmd interface from the go-mesh-gen.dev/v4/cmd
// package and adds a Run method.
//
// Run runs the cli app within this command and exits on error.
type CLI interface {
	App() *cli.App
	Run() error
}

type cmd struct {
	app *cli.App
}

// App returns the cli app within this command.
func (c *cmd) App() *cli.App {
	return c.app
}

// App returns the cli app within the default command.
func App() *cli.App {
	return DefaultCLI.App()
}

// Register appends commands to the default app.
func Register(cmds ...*cli.Command) {
	app := DefaultCLI.App()
	app.Commands = append(app.Commands, cmds...)
}

// Run runs the cli app within this command and exits on error.
func (c *cmd) Run() error {
	return c.app.Run(os.Args)
}

// Run runs the cli app within the default command. On error, it prints the
// error message and exits.
func Run() {
	if err := DefaultCLI.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// NewCLI returns a new command.
func NewCLI() CLI {

	c := new(cmd)
	c.app = cli.NewApp()
	c.app.Name = name
	c.app.Usage = description
	c.app.Version = version

	return c
}
