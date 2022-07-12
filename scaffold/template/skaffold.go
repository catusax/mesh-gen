package template

// SkaffoldCFG is the Skaffold config template used for new projects.
var SkaffoldCFG = `---

apiVersion: skaffold/v2beta21
kind: Config
metadata:
  name: {{.Service}}
build:
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
`
