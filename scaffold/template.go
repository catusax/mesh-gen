package scaffold

import (
	"github.com/catusax/mesh-gen/scaffold/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

var templateDir = ".template"

func GenTemplate(path string) error {
	var dir = filepath.Join(path, templateDir)

	err := os.MkdirAll(dir, 0750)
	if err != nil {
		return err
	}
	return WriteTemplate(dir)
}

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

func WriteTemplate(path string) error {
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

		err := ioutil.WriteFile(filepath.Join(path, templates[i].Path+".tmpl"), []byte(templates[i].Value), 0644)
		if err != nil {
			return err
		}
	}

	return nil

}
