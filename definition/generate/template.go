package generate

import "text/template"

//nolint:gochecknoglobals
var (
	tpl          = template.Must(template.New("tpls").Parse(baseTemplate))
	baseTemplate = `// Code generated gitoa.ru/go-4devs/config DO NOT EDIT.
package {{.Pkg}}

import (
    {{range .Imports}}
    {{- .Alias }}"{{ .Package }}"
    {{end}}
)

func With{{.StructName}}Log(log func(context.Context, string, ...any)) func(*{{.StructName}}) {
	return func(ci *{{.StructName}}) {
		ci.log = log
	}
}

func New{{.StructName}}(prov config.Provider, opts ...func(*{{.StructName}})) {{.StructName}} {
	i := {{.StructName}}{
		Provider: prov,
		log: func(_ context.Context, format string, args ...any) {
			fmt.Printf(format, args...)
		},
	}

	for _, opt := range opts {
		opt(&i)
	}

	return i
}

type {{.StructName}} struct {
	config.Provider
	log   func(context.Context, string, ...any)
}
`
)
