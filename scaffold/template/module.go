package template

// Module is the go.mod template used for new projects.
var Module = `module {{.Vendor}}{{.Service}}

go 1.16

require (
)

`
