package main

import (
	"bytes"
	"text/template"

	goimports "golang.org/x/tools/imports"
)

var tmpl = template.Must(template.New("").Parse(`// Code generated by ifacegen. DO NOT EDIT.

package {{ .PackageName }}

import (
  {{ range .Imports }}
    {{ . -}}
  {{ end}}
)

type {{ .InterfaceName }} interface {
  {{- range .Methods -}}
  {{- with .Docs -}}
  {{ range . }}
    {{ . -}}
  {{ end -}}
  {{ end }}
    {{ .Code -}}
  {{- end -}}
}`))

func GenerateFile(files []*File) ([]byte, error) {
	imports := make([]string, 0)
	methods := make([]*Method, 0)
	for _, f := range files {
		imports = append(imports, f.Imports...)
		methods = append(methods, f.Methods...)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]any{
		"Imports":       imports,
		"PackageName":   files[0].PackageName,
		"InterfaceName": opts.InterfaceName,
		"Methods":       methods,
	}); err != nil {
		return nil, err
	}

	return goimports.Process("", buf.Bytes(), &goimports.Options{
		TabIndent: true,
		AllErrors: true,
		Fragment:  true,
		Comments:  true,
	})
}
