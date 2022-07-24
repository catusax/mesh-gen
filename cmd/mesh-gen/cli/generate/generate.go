package generate

import (
	"bufio"
	"errors"
	"fmt"
	generator2 "github.com/catusax/mesh-gen/scaffold"
	"github.com/catusax/mesh-gen/scaffold/template"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"

	mcli "github.com/catusax/mesh-gen/cmd/mesh-gen/cli"
)

func init() {
	mcli.Register(&cli.Command{
		Name:   "generate",
		Usage:  "Generate project Skaffold & Kubernetes resource template files after the fact",
		Action: Skaffold,
	})
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

	g := generator2.New(
		generator2.Service(service),
		generator2.Vendor(vendor),
		generator2.Directory("."),
		generator2.ContainerTag(tag),
		generator2.ContainerVersion(version),
		generator2.Port(port),
		generator2.Namespace(namespace),
	)

	files := []generator2.File{
		{".dockerignore", generator2.GetTemplate(template.DockerIgnore)},
		{"Dockerfile", generator2.GetTemplate(template.Dockerfile)},
		{filepath.Join("resources", "configmap.yaml"), generator2.GetTemplate(template.KubernetesEnv)},
		{filepath.Join("resources", "deployment.yaml"), generator2.GetTemplate(template.KubernetesDeployment)},
		{"skaffold.yaml", generator2.GetTemplate(template.SkaffoldCFG)},
	}

	if err := g.Generate(files); err != nil {
		return err
	}

	fmt.Println("skaffold project template files generated")

	return nil
}

func getService() (string, error) {

	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	name := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		if strings.HasPrefix(scanner.Text(), "NAME=") {
			name = scanner.Text()
		}
	}

	if name == "" {
		fmt.Println("Makefile is missing NAME or VERSION variables")
		return "", errors.New("could not get container tag")
	}

	name = name[strings.Index(name, "=")+1:]

	return name, nil
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
	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	registryPrefix := ""
	name := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "REGISTRY_PREFIX=") {
			registryPrefix = scanner.Text()
		}

		if strings.HasPrefix(scanner.Text(), "NAME=") {
			name = scanner.Text()
		}

	}

	name = name[strings.Index(name, "=")+1:]
	registryPrefix = registryPrefix[strings.Index(registryPrefix, "=")+1:]

	if name == "" {
		fmt.Println("Makefile is missing NAME variables")
		return "", errors.New("could not get container tag")
	}

	if registryPrefix == "" {
		return name, nil
	}

	return registryPrefix + "/" + name, nil
}

func getVersion() (string, error) {
	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	version := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "VERSION=") {
			version = scanner.Text()
			break
		}
	}

	version = version[strings.Index(version, "=")+1:]

	if version == "" {
		fmt.Println("Makefile is missing VERSION variable")
		return "", errors.New("could not get container tag")
	}

	return version, nil
}

func getPort() (string, error) {
	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	port := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "PORT=") {
			port = scanner.Text()
			break
		}
	}

	port = port[strings.Index(port, "=")+1:]

	if port == "" {
		fmt.Println("Makefile is missing VERSION variable")
		return "", errors.New("could not get container tag")
	}

	return port, nil
}

func getNamespace() (string, error) {
	f, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer f.Close()

	namespace := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "NAMESPACE=") {
			namespace = scanner.Text()
			break
		}
	}

	namespace = namespace[strings.Index(namespace, "=")+1:]

	if namespace == "" {
		fmt.Println("Makefile is missing NAMESPACE variable")
		return "", errors.New("could not get container tag")
	}

	return namespace, nil
}
