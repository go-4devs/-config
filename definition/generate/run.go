package generate

import (
	"bytes"
	"fmt"
	"io"

	"gitoa.ru/go-4devs/config/definition"
)

func Run(w io.Writer, pkgName string, defs definition.Definition, viewOpt ViewOption) error {
	gen := Generator{
		pkg:        pkgName,
		ViewOption: viewOpt,
		Imp:        NewImports(),
	}

	gen.Imp.Adds("gitoa.ru/go-4devs/config", "fmt", "context")

	var view bytes.Buffer
	err := defs.View(func(o definition.Option) error {
		return gen.Handle(&view, &gen, o)
	})
	if err != nil {
		return fmt.Errorf("render options:%w", err)
	}

	if err := tpl.Execute(w, gen); err != nil {
		return fmt.Errorf("render base:%w", err)
	}

	_, cerr := io.Copy(w, &view)
	if cerr != nil {
		return fmt.Errorf("copy error:%w", cerr)
	}

	return nil
}
