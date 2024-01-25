package generate

import (
	"errors"

	"github.com/iancoleman/strcase"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrWrongType    = errors.New("wrong type")
	ErrWrongFormat  = errors.New("wrong format")
)

func FuncName(in string) string {
	return strcase.ToCamel(in)
}
