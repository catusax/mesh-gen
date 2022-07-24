package new

import (
	"fmt"
	generator2 "github.com/catusax/mesh-gen/scaffold"
	"github.com/catusax/mesh-gen/scaffold/template"
	"os"
	"path"
	"strings"

	mcli "github.com/catusax/mesh-gen/cmd/mesh-gen/cli"
	"github.com/urfave/cli/v2"
)

// NewCommand returns a new cli command.
func init() {
	mcli.Register(&cli.Command{
		Name:   "new",
		Usage:  "Create a service template",
		Action: Service,
	})
}

// Service creates a new service project template. Exits on error.
func Service(ctx *cli.Context) error {
	return createProject(ctx)
}

func createProject(ctx *cli.Context) error {
	arg := ctx.Args().First()
	if len(arg) == 0 {
		return cli.ShowSubcommandHelp(ctx)
	}
	name, vendor := getNameAndVendor(arg)

	dir := name

	var tag = name

	if defaultRegistryPrefix != "" {
		tag = defaultRegistryPrefix + "/" + name
	}

	if path.IsAbs(dir) {
		fmt.Println("must provide a relative path as service name")
		return nil
	}

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return fmt.Errorf("%s already exists", dir)
	}

	fmt.Printf("creating service %s\n", name)

	g := generator2.New(
		generator2.Service(name),
		generator2.Vendor(vendor),
		generator2.Directory(dir),
		generator2.ContainerTag(tag),
		generator2.ContainerVersion(defaultVersion),
		generator2.Port(defaultPort),
		generator2.Namespace(defaultNamespace),
		generator2.RegistryPrefix(defaultRegistryPrefix),
		generator2.Mesh(defaultMesh),
	)

	files := []generator2.File{
		{".dockerignore", generator2.GetTemplate(template.DockerIgnore)},
		{".gitignore", generator2.GetTemplate(template.GitIgnore)},
		{"Dockerfile", generator2.GetTemplate(template.Dockerfile)},
		{"Makefile", generator2.GetTemplate(template.Makefile)},
		{"go.mod", generator2.GetTemplate(template.Module)},
	}

	//service and skaffold files
	files = append(files, []generator2.File{
		{"handler/" + name + ".go", generator2.GetTemplate(template.HandlerSRV)},
		{"main.go", generator2.GetTemplate(template.MainSRV)},
		{"proto/" + name + ".proto", generator2.GetTemplate(template.ProtoSRV)},
		{"resources/configmap.yaml", generator2.GetTemplate(template.KubernetesEnv)},
		{"resources/deployment.yaml", generator2.GetTemplate(template.KubernetesDeployment)},
		{"skaffold.yaml", generator2.GetTemplate(template.SkaffoldCFG)},
	}...)

	if err := g.Generate(files); err != nil {
		return err
	}

	var comments = protoComments(name, dir)

	for _, comment := range comments {
		fmt.Println(comment)
	}

	return nil
}

func protoComments(name, dir string) []string {
	return []string{
		"\ndownload protoc zip packages (protoc-$VERSION-$PLATFORM.zip) and install:",
		"\nvisit https://github.com/protocolbuffers/protobuf/releases/latest",
		"\ncd " + dir,
		"make init proto update tidy",
	}
}

func getNameAndVendor(s string) (string, string) {
	var n string
	var v string

	if i := strings.LastIndex(s, "/"); i == -1 {
		n = s
		v = ""
	} else {
		n = s[i+1:]
		v = s[:i+1]
	}

	return n, v
}
