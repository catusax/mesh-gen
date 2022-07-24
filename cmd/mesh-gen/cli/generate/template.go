package generate

import (
	mcli "github.com/catusax/mesh-gen/cmd/mesh-gen/cli"
	generator2 "github.com/catusax/mesh-gen/scaffold"
	"github.com/urfave/cli/v2"
	"os"
)

func init() {
	mcli.Register(&cli.Command{
		Name:   "template",
		Usage:  "Generate project template on this folder",
		Action: Template,
	})
}

func Template(ctx *cli.Context) error {
	if ctx.Args().First() == "global" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return generator2.GenTemplate(home)
	} else {
		return generator2.GenTemplate("")
	}
}
