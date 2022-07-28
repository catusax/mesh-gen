package template

// SkaffoldCFG is the Skaffold config template used for new projects.
var SkaffoldCFG = Template{
	Path: "skaffold.yaml",
	Value: `---

apiVersion: skaffold/v2beta21
kind: Config
metadata:
  name: {{.Service}}
build:
  tagPolicy:
    envTemplate:
      template: "{{.Version}}"
  artifacts:
    - image: {{.ContainerTag}}
      docker:
        noCache: false
  local:
    useBuildkit: true
deploy:
  kubectl:
    manifests:
      - resources/*.yaml
`,
}
