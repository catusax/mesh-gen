package main

import (
	"github.com/catusax/mesh-gen/cmd/mesh-gen/cli"
	// register commands
	_ "github.com/catusax/mesh-gen/cmd/mesh-gen/cli/generate"
	_ "github.com/catusax/mesh-gen/cmd/mesh-gen/cli/new"
)

func main() {
	cli.Run()
}
