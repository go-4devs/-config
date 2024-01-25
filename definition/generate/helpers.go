package generate

import (
	"errors"

	"github.com/iancoleman/strcase"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

func FuncName(in string) string {
	return strcase.ToCamel(in)
}
