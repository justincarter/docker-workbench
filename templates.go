package main

var templateAppHelp = `{{.Name}} v{{.Version}}
{{.Usage}}

Usage: 
  {{.Name}} {{if .Flags}}[options] {{end}}COMMAND
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
Run '{{.Name}} help COMMAND' for more information on a command.
`

var templateCommandHelp = `{{.Usage}}

Usage: 
  docker-workbench {{.Name}}{{if .Flags}} [options]{{end}}

{{- if .Description}}
Description:
  {{.Description}}
{{- end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}`
