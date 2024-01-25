package generate

import (
	"fmt"
	"strconv"
	"strings"
)

func NewImports() Imports {
	return Imports{
		data: make(map[string]string),
	}
}

type Imports struct {
	data map[string]string
}

func (i Imports) Imports() []Import {
	imports := make([]Import, 0, len(i.data))
	for name, alias := range i.data {
		imports = append(imports, Import{
			Package: name,
			Alias:   alias,
		})
	}

	return imports
}

func (i *Imports) Short(fullType string) (string, error) {
	idx := strings.LastIndexByte(fullType, '.')
	if idx == -1 {
		return "", fmt.Errorf("unexpected")
	}
	if alias, ok := i.data[fullType[:idx]]; ok {
		return alias + fullType[idx:], nil
	}

	return "", fmt.Errorf("%w alias for pkg %v", ErrNotFound, fullType[:idx])
}

func (i *Imports) AddType(fullType string) (string, error) {
	idx := strings.LastIndexByte(fullType, '.')
	if idx == -1 {
		return "", fmt.Errorf("unexpected")
	}

	imp := i.Add(fullType[:idx])

	return imp.Alias + fullType[idx:], nil

}

func (i *Imports) Adds(pkgs ...string) {
	for _, pkg := range pkgs {
		i.Add(pkg)
	}
}

func (i *Imports) Add(pkg string) Import {
	alias := pkg
	idx := strings.LastIndexByte(pkg, '/')
	if idx != -1 {
		alias = pkg[idx+1:]
	}
	if al, ok := i.data[pkg]; ok {
		return Import{Package: pkg, Alias: al}
	}

	for _, al := range i.data {
		if al == alias {
			alias += strconv.Itoa(len(i.data))
		}
	}
	i.data[pkg] = alias

	return Import{
		Alias:   alias,
		Package: pkg,
	}
}

type Import struct {
	Alias   string
	Package string
}
