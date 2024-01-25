package generate

import (
	"fmt"
	"io"

	"gitoa.ru/go-4devs/config/definition"
)

type Generator struct {
	pkg string
	ViewOption
	Imp           Imports
	errs          []error
	defaultErrors []string
}

func (g Generator) Pkg() string {
	return g.pkg
}

func (g Generator) Imports() []Import {
	return g.Imp.Imports()
}

func (g Generator) Handle(w io.Writer, data Handler, opt definition.Option) error {
	handle := get(opt.Kind())

	return handle(w, data, opt)
}

func (g Generator) StructName() string {
	return FuncName(g.Prefix + "_" + g.Struct + "_" + g.Suffix)
}

func (g Generator) Options() ViewOption {
	return g.ViewOption
}

func (g Generator) Keys() []string {
	return nil
}

func (g Generator) DefaultErrors() []string {
	if len(g.defaultErrors) > 0 {
		return g.defaultErrors
	}

	if len(g.ViewOption.Errors.Default) > 0 {
		g.Imp.Adds("errors")
	}

	g.defaultErrors = make([]string, len(g.ViewOption.Errors.Default))
	for idx, name := range g.ViewOption.Errors.Default {
		short, err := g.AddType(name)
		if err != nil {
			g.errs = append(g.errs, fmt.Errorf("add default error[%d]:%w", idx, err))

			return nil
		}

		g.defaultErrors[idx] = short
	}

	return g.defaultErrors
}

func (g *Generator) AddType(pkg string) (string, error) {
	return g.Imp.AddType(pkg)
}
