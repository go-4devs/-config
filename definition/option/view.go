package option

import (
	"bytes"
	"embed"
	"encoding"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/template"
	"time"

	"gitoa.ru/go-4devs/config/definition"
	"gitoa.ru/go-4devs/config/definition/generate"
)

//go:embed tpl/*
var tpls embed.FS

//nolint:gochecknoglobals
var tpl = template.Must(template.New("tpls").ParseFS(tpls, "tpl/*.tmpl"))

//nolint:gochecknoinits
func init() {
	generate.MustAdd(Kind, Handle(tpl.Lookup("option.tmpl")))
}

func Handle(tpl *template.Template) generate.Handle {
	return func(w io.Writer, h generate.Handler, o definition.Option) error {
		opt, _ := o.(Option)
		if err := tpl.Execute(w, View{Option: opt, Handler: h}); err != nil {
			return fmt.Errorf("option tpl:%w", err)
		}

		return nil
	}
}

type View struct {
	Option
	generate.Handler
}

func (v View) Context() bool {
	return v.Options().Context
}

func (v View) FuncName() string {
	if funcName, ok := v.Option.Params.Get(ViewParamFunctName); ok {
		name, _ := funcName.(string)

		return name
	}

	return generate.FuncName(v.Name)
}

func (v View) Description() string {
	if desc, ok := v.Option.Params.Get(ViewParamDescription); ok {
		description, _ := desc.(string)

		return description
	}

	return v.Option.Description
}

func (v View) Default() string {
	switch data := v.Option.Default.(type) {
	case time.Time:
		return fmt.Sprintf("time.Parse(%q,time.RFC3339Nano)", data.Format(time.RFC3339Nano))
	case time.Duration:
		return fmt.Sprintf("time.ParseDuration(%q)", data)
	default:
		return fmt.Sprintf("%#v, nil", data)
	}
}

func (v View) HasDefault() bool {
	return v.Option.Default != nil
}

func (v View) ParentKeys() string {
	if len(v.Handler.Keys()) > 0 {
		return `"` + strings.Join(v.Handler.Keys(), `","`) + `",`
	}

	return ""
}

func (v View) Type() string {
	slice := ""

	if vtype, ok := v.Option.Type.(string); ok {
		if strings.Contains(vtype, ".") {
			if name, err := v.AddType(vtype); err == nil {
				return slice + name
			}
		}

		return vtype
	}

	rtype := reflect.TypeOf(v.Option.Type)

	if rtype.PkgPath() == "" {
		return rtype.String()
	}

	if rtype.Kind() == reflect.Slice {
		slice = "[]"
	}

	short, err := v.AddType(rtype.PkgPath() + "." + rtype.Name())
	if err != nil {
		return err.Error()
	}

	return slice + short
}

func (v View) FuncType() string {
	return generate.FuncName(v.Type())
}

func (v View) Parse(valName string, value string, keys []string) string {
	h := parser(v.Option.Type)

	data, err := h(ParseData{
		Value:   value,
		ValName: valName,
		Keys:    keys,
		View:    v,
	})
	if err != nil {
		return err.Error()
	}

	return data
}

//nolint:gochecknoglobals,unparam
var parses = map[string]func(data ParseData) (string, error){
	typesIntreface[0].Name(): func(data ParseData) (string, error) {
		var b bytes.Buffer

		err := tpl.ExecuteTemplate(&b, "unmarshal_text.tmpl", data)
		if err != nil {
			return "", fmt.Errorf("execute unmarshal text:%w", err)
		}

		return b.String(), nil
	},
	typesIntreface[1].Name(): func(data ParseData) (string, error) {
		var b bytes.Buffer

		err := tpl.ExecuteTemplate(&b, "unmarshal_json.tmpl", data)
		if err != nil {
			return "", fmt.Errorf("execute unmarshal json:%w", err)
		}

		return b.String(), nil
	},
	TypeInt:     internal,
	TypeInt64:   internal,
	TypeBool:    internal,
	TypeString:  internal,
	TypeFloat64: internal,
	TypeUint:    internal,
	TypeUint64:  internal,
	"time.Duration": func(data ParseData) (string, error) {
		return fmt.Sprintf("return %s.ParseDuration()", data.ValName), nil
	},
	"time.Time": func(data ParseData) (string, error) {
		return fmt.Sprintf("return %s.ParseTime()", data.ValName), nil
	},
	"any": func(data ParseData) (string, error) {
		return fmt.Sprintf("return %[2]s, %[1]s.Unmarshal(&%[2]s)", data.ValName, data.Value), nil
	},
}

func internal(data ParseData) (string, error) {
	var b bytes.Buffer

	err := tpl.ExecuteTemplate(&b, "parse.tmpl", data)
	if err != nil {
		return "", fmt.Errorf("execute parse.tmpl:%w", err)
	}

	return b.String(), nil
}

//nolint:gochecknoglobals
var typesIntreface = [...]reflect.Type{
	reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem(),
	reflect.TypeOf((*json.Unmarshaler)(nil)).Elem(),
}

func parser(data any) func(ParseData) (string, error) {
	vtype := reflect.TypeOf(data)
	name := vtype.Name()

	if v, ok := data.(string); ok {
		name = v
	}

	if vtype.Kind() == reflect.Slice {
		return parses["any"]
	}

	if h, ok := parses[name]; ok {
		return h
	}

	for _, extypes := range typesIntreface {
		if vtype.Implements(extypes) {
			return parses[extypes.Name()]
		}

		if vtype.Kind() != reflect.Ptr && reflect.PointerTo(vtype).Implements(extypes) {
			return parses[extypes.Name()]
		}
	}

	return parses["any"]
}

type ParseData struct {
	Value   string
	ValName string
	Keys    []string
	View
}
