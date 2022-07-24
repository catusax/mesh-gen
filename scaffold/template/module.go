package template

// Module is the go.mod template used for new projects.
var Module = Template{
	Path: "go.mod",
	Value: `module {{.Vendor}}{{.Service}}

go 1.16

require (
)

`,
}
