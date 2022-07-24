package scaffold

import (
	"fmt"
	"github.com/catusax/mesh-gen/scaffold/template"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"os"
	"path/filepath"
)

var templateDir = ".template"

func GenTemplate(path string, force bool) error {
	var dir = filepath.Join(path, templateDir)

	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}
	return WriteTemplate(dir, force)
}

//GetTemplate search template in config dir,ordered by:
// current dir > .template > ../.template > ~/.template > default
func GetTemplate(t template.Template) string {
	var file []byte
	var searchFiles = []string{
		t.Path + ".tmpl",
		filepath.Join(templateDir, t.Path+".tmpl"),
		filepath.Join("..", templateDir, t.Path+".tmpl"),
	}
	home, err := os.UserHomeDir()
	if err == nil {
		searchFiles = append(searchFiles, filepath.Join(home, templateDir, t.Path+".tmpl"))
	}

	for _, filePath := range searchFiles {
		file, err = ioutil.ReadFile(filePath)
		if err == nil {
			break
		}
	}

	if len(file) == 0 {
		return t.Value
	}

	return string(file)

}

//WriteTemplate Write Template to dir ,so user can edit
func WriteTemplate(path string, force bool) error {
	var templates = []template.Template{
		template.Dockerfile,
		template.DockerIgnore,
		template.GitIgnore,
		template.HandlerSRV,
		template.KubernetesDeployment,
		template.KubernetesEnv,
		template.MainSRV,
		template.Makefile,
		template.Module,
		template.ProtoSRV,
		template.SkaffoldCFG,
	}

	for i := range templates {
		dir, _ := filepath.Split(templates[i].Path)
		if dir != "" {
			err := os.MkdirAll(filepath.Join(path, dir), 0777)
			if err != nil {
				return err
			}
		}

		err := writeFile(filepath.Join(path, templates[i].Path+".tmpl"), &templates[i].Value, force)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(file string, content *string, force bool) error {
	if force {
		return ioutil.WriteFile(file, []byte(*content), 0644)
	}
	old, err := ioutil.ReadFile(file)
	if err == nil {
		diffMatch := diffmatchpatch.New()
		diff := diffMatch.DiffMain(string(old), *content, true)

		if len(diff) == 0 || len(diff) == 1 && diff[0].Type == diffmatchpatch.DiffEqual {
			return nil
		}
		fmt.Println(file, " is different from latest version! ", len(diff))
		fmt.Println(diffMatch.DiffPrettyText(diff))
		fmt.Println()
		for {
			var yes string
			fmt.Print(file, " is different from latest version!", "override? y/n :")
			fmt.Scanf("%s", &yes)
			if yes == "y" || yes == "Y" {
				break
			}
			if yes == "n" || yes == "N" {
				return nil
			}
		}

	} else {
		err = ioutil.WriteFile(file, []byte(*content), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
