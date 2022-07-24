package template

// GitIgnore is the .gitignore template used for new projects.
var GitIgnore = Template{
	Path: ".gitignore",
	Value: `# don't commit the service binary to vcs
{{.Service}}

.idea
.vscode
`,
}
