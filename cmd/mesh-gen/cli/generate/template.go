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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "force",
				Aliases: []string{
					"f",
				},
				Usage: "force override existing files",
			},
		},
	})
}

func Template(ctx *cli.Context) error {
	force := ctx.Bool("force")
	if ctx.Args().First() == "global" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		return generator2.GenTemplate(home, force)
	} else {
		return generator2.GenTemplate("", force)
	}
}
