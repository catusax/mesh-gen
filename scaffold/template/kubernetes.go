package template

// KubernetesEnv is a Kubernetes configmap manifest template used for
// environment variables in new projects.
var KubernetesEnv = `---

apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Service}}-env
  namespace: {{.Namespace}}
data:
  CONTAINER: kubernetes
`

// KubernetesDeployment is a Kubernetes deployment manifest template used for
// new projects.
var KubernetesDeployment = `---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Service}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Service}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Service}}
  template:
    metadata:
      labels:
        app: {{.Service}}
    spec:
      containers:
      - name: {{.Service}}
        image: {{.ContainerTag}}:{{.Version}}
        envFrom:
        - configMapRef:
            name: {{.Service}}-env


---
apiVersion: v1
kind: Service
metadata:
  name: {{.Service}}-grpc
  namespace: {{.Namespace}}
  labels:
    app: {{.Service}}-grpc
  annotations:
    mesh.traefik.io/traffic-type: "http"
    mesh.traefik.io/scheme: "h2c"
spec:
  ports:
    - port: {{.Port}}
      name: grpc
  selector:
    app: {{.Service}}
`
