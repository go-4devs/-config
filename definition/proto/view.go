package proto

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"gitoa.ru/go-4devs/config/definition"
	"gitoa.ru/go-4devs/config/definition/generate"
	"gitoa.ru/go-4devs/config/definition/option"
)

func init() {
	generate.MustAdd(Kind, handle)
}

func handle(w io.Writer, data generate.Handler, opt definition.Option) error {
	proto, ok := opt.(Proto)
	if !ok {
		return fmt.Errorf("uexepected type:%T", opt)
	}
	if viewOpt, ok := proto.Option.(option.Option); ok {
		viewOpt = viewOpt.WithParams(
			definition.Param{
				Name:  option.ViewParamFunctName,
				Value: generate.FuncName(proto.Name) + generate.FuncName(viewOpt.Name),
			},
			definition.Param{
				Name:  option.ViewParamDescription,
				Value: proto.Description + " " + viewOpt.Description,
			},
		)

		return option.Handle(tpl)(w, data, viewOpt)
	}

	return fmt.Errorf("not support option type")
}

var tpl = template.Must(template.New("tpls").Funcs(template.FuncMap{"join": strings.Join}).Parse(templateOption))

var templateOption = `// read{{.FuncName}} {{.Description}}.
func (i {{.StructName}}) read{{.FuncName}}(ctx context.Context, key string) (v {{.Type}},e error) {
	val, err := i.Value(ctx, {{ .ParentKeys }} key, "{{.Name}}")
	if err != nil {
        {{if .HasDefault}}
        {{$default := .Default}}
        {{range .DefaultErrors}}
        if errors.Is(err,{{.}}){
            return {{$default}}
		}
        {{end}}
        {{end}}
        return v, fmt.Errorf("read {{.Keys}}:%w",err)
	}

	{{.Parse "val" "v" .Keys }}
}

// Read{{.FuncName}} {{.Description}}.
func (i {{.StructName}}) Read{{.FuncName}}(ctx context.Context, key string) ({{.Type}}, error) {
    return i.read{{.FuncName}}(ctx, key)
}

// {{.FuncName}} {{.Description}}.
func (i {{.StructName}}) {{.FuncName}}({{if .Context}} ctx context.Context, {{end}} key string) {{.Type}} {
    {{if not .Context}} ctx := context.Background() {{end}}
	val, err := i.read{{.FuncName}}(ctx, key)
	if err != nil {
		i.log(ctx, "get {{.Keys}}: %v", err)
	}

	return val
}
`
