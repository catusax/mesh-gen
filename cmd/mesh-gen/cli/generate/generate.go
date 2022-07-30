package generate

import (
	"bufio"
	"errors"
	"fmt"
	generator2 "github.com/catusax/mesh-gen/scaffold"
	"github.com/catusax/mesh-gen/scaffold/template"
	"github.com/urfave/cli/v2"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	mcli "github.com/catusax/mesh-gen/cmd/mesh-gen/cli"
)

func init() {
	cmd := &cli.Command{
		Name:    "generate",
		Aliases: []string{"gen"},
		Usage:   "Generate project Skaffold & Kubernetes resource template files after the fact",
		Action:  Skaffold,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "all",
				Aliases: []string{
					"A",
				},
				Usage: "Regenerate all files except ./handler ./proto ./go.mod",
			},
			&cli.BoolFlag{
				Name: "fmt",
				Aliases: []string{
					"f",
				},
				Usage: "Format code using go fmt .",
			},
		},
	}
	mcli.Register(cmd)
}

// Skaffold generates Skaffold template files in the current working directory.
// Exits on error.
func Skaffold(ctx *cli.Context) error {

	service, err := getService()
	if err != nil {
		return err
	}

	vendor, err := getServiceVendor(service)
	if err != nil {
		return err
	}

	tag, err := getContainerTag()
	if err != nil {
		return err
	}

	version, err := getVersion()
	if err != nil {
		return err
	}

	port, err := getPort()
	if err != nil {
		return err
	}

	namespace, err := getNamespace()
	if err != nil {
		return err
	}

	registryPrefix, _ := getRegistryPrefix()

	replica, _ := getReplica()

	var replicaNumber int
	replicaNumber, err = strconv.Atoi(replica)
	if err != nil {
		replicaNumber = 3
	}

	mesh, _ := getMesh()

	g := generator2.New(
		generator2.Service(service),
		generator2.Vendor(vendor),
		generator2.Directory("."),
		generator2.ContainerTag(tag),
		generator2.ContainerVersion(version),
		generator2.Port(port),
		generator2.Namespace(namespace),
		generator2.RegistryPrefix(registryPrefix),
		generator2.Mesh(mesh),
		generator2.Replica(replicaNumber),
	)

	files := []generator2.File{
		{filepath.Join("resources", "configmap.yaml"), generator2.GetTemplate(template.KubernetesEnv)},
		{filepath.Join("resources", "deployment.yaml"), generator2.GetTemplate(template.KubernetesDeployment)},
		{"skaffold.yaml", generator2.GetTemplate(template.SkaffoldCFG)},
	}

	if ctx.Bool("all") {
		files = append(files,
			generator2.File{Path: ".gitignore", Template: generator2.GetTemplate(template.GitIgnore)},
			generator2.File{Path: "Makefile", Template: generator2.GetTemplate(template.Makefile)},
			generator2.File{Path: "main.go", Template: generator2.GetTemplate(template.MainSRV)},
			generator2.File{Path: ".dockerignore", Template: generator2.GetTemplate(template.DockerIgnore)},
			generator2.File{Path: "Dockerfile", Template: generator2.GetTemplate(template.Dockerfile)},
		)
	}

	if err := g.Generate(files); err != nil {
		return err
	}

	if ctx.Bool("fmt") {
		exec.Command("go", "fmt", ".")
	}

	fmt.Println("skaffold project template files generated")

	return nil
}

func getService() (string, error) {
	return ReadKey("NAME")
}

func getServiceVendor(s string) (string, error) {
	f, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer f.Close()

	line := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "module ") {
			line = scanner.Text()
			break

		}
	}
	if line == "" {
		return "", nil
	}

	module := line[strings.LastIndex(line, " ")+1:]
	if module == s {
		return "", nil
	}

	return module[:strings.LastIndex(module, "/")] + "/", nil
}

func getContainerTag() (string, error) {
	name, err := ReadKey("NAME")
	if err != nil {
		return "", errors.New("could not get container tag")
	}

	registryPrefix, _ := ReadKey("REGISTRY_PREFIX")
	if registryPrefix == "" {
		return name, nil
	}

	return registryPrefix + "/" + name, nil
}

func getVersion() (string, error) {
	return ReadKey("VERSION")
}

func getPort() (string, error) {
	return ReadKey("PORT")
}

func getNamespace() (string, error) {
	return ReadKey("NAMESPACE")
}

func getRegistryPrefix() (string, error) {
	return ReadKey("REGISTRY_PREFIX")
}

func getMesh() (string, error) {
	return ReadKey("MESH")
}

func getReplica() (string, error) {
	return ReadKey("REPLICA")
}

func ReadKey(key string) (string, error) {
	if env, ok := os.LookupEnv("SRV_" + key); ok {
		return env, nil
	}

	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	value := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), key+"=") {
			value = scanner.Text()
		}
	}

	if value == "" {
		return "", errors.New("could not get " + key)
	}

	position := strings.Index(value, "=")
	if position == -1 {
		return "", errors.New("could not get " + key)
	}

	value = value[position+1:]

	return value, nil
}
