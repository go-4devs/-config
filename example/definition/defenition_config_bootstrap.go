//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"go/format"
	"os"

	"gitoa.ru/go-4devs/config/definition"
	"gitoa.ru/go-4devs/config/definition/generate"
	"gitoa.ru/go-4devs/config/eample"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	def := definition.New()
	if err := eample.Configure(ctx, &def); err != nil {
		return err
	}

	f, err := os.Create("eample/defenition_input.go")
	if err != nil {
		return err
	}

	gerr := generate.Run(f, "eample", def, generate.ViewOption{
		Struct: "Configure",
		Suffix: "Input",
		Errors: generate.ViewErrors{
			Default: []string{
				"gitoa.ru/go-4devs/config.ErrVariableNotFound",
			},
		},
	})

	if gerr != nil {
		return gerr
	}

	in, err := os.ReadFile(f.Name())
	if err != nil {
		return err
	}

	out, err := format.Source(in)
	if err != nil {
		return err
	}

	return os.WriteFile(f.Name(), out, 0644)
}
