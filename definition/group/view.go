package group

import (
	"fmt"
	"io"
	"text/template"

	"gitoa.ru/go-4devs/config/definition"
	"gitoa.ru/go-4devs/config/definition/generate"
)

func init() {
	generate.MustAdd(Kind, handle)
}

func handle(w io.Writer, data generate.Handler, option definition.Option) error {
	group, ok := option.(Group)
	if !ok {
		return fmt.Errorf("uexepected type:%T", option)
	}
	viewData := View{
		Group:      group,
		ParentName: data.StructName(),
		ViewOption: data.Options(),
	}

	err := tpl.Execute(w, viewData)

	if err != nil {
		return fmt.Errorf("render group:%w", err)
	}

	childData := ChildData{
		Handler:    data,
		structName: viewData.StructName(),
		keys:       append(data.Keys(), group.Name),
	}
	for idx, child := range group.Options {
		if cerr := data.Handle(w, childData, child); cerr != nil {
			return fmt.Errorf("render group child[%d]:%w", idx, cerr)
		}
	}

	return nil
}

type ChildData struct {
	generate.Handler
	structName string
	keys       []string
}

func (c ChildData) StructName() string {
	return c.structName
}

func (v ChildData) Keys() []string {
	return v.keys
}

type View struct {
	Group
	ParentName string
	generate.ViewOption
}

func (v View) FuncName() string {
	return generate.FuncName(v.Name)
}

func (v View) StructName() string {
	return generate.FuncName(v.Prefix + v.Name + v.Suffix)
}

var tpl = template.Must(template.New("tpls").Parse(tplw))

var tplw = `type {{.StructName}} struct {
    {{.ParentName}}
}

// {{.FuncName}} {{.Description}}.
func (i {{.ParentName}}) {{.FuncName}}() {{.StructName}} {
	return {{.StructName}}{i}
}
`
