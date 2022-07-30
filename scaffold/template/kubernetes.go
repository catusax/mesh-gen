package template

import "path/filepath"

// KubernetesEnv is a Kubernetes configmap manifest template used for
// environment variables in new projects.
var KubernetesEnv = Template{
	Path: filepath.Join("resources", "configmap.yaml"),
	Value: `---

apiVersion: v1
kind: ConfigMap
metadata:
  name: {{dash .Service}}-env
  namespace: {{.Namespace}}
data:
  CONTAINER: kubernetes
`,
}

// KubernetesDeployment is a Kubernetes deployment manifest template used for
// new projects.
var KubernetesDeployment = Template{
	Path: filepath.Join("resources", "deployment.yaml"),
	Value: `---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{dash .Service}}
  namespace: {{.Namespace}}
  labels:
    app: {{dash .Service}}
spec:
  replicas: {{ .Replica }}
  selector:
    matchLabels:
      app: {{dash .Service}}
  template:
    metadata:
      labels:
        app: {{dash .Service}}
{{- if eq .Mesh "istio" }}
      annotations:
        inject.istio.io/templates: grpc-agent
        proxy.istio.io/config: '{"holdApplicationUntilProxyStarts": true}'
{{ end -}}
    spec:
      containers:
        - name: {{dash .Service}}
          image: {{.ContainerTag}}:{{.Version}}
          ports:
            - containerPort: {{.Port}}
          envFrom:
            - configMapRef:
                name: {{dash .Service}}-env


---
apiVersion: v1
kind: Service
metadata:
  name: {{dash .Service}}-grpc
  namespace: {{.Namespace}}
  labels:
    app: {{dash .Service}}-grpc
{{- if eq .Mesh "traefik-mesh" }}
  annotations:
    mesh.traefik.io/traffic-type: "http"
    mesh.traefik.io/scheme: "h2c"
{{ end -}}
spec:
  ports:
    - port: {{.Port}}
      name: grpc
      appProtocol: grpc
  selector:
    app: {{dash .Service}}
`,
}
