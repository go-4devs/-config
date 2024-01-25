package generate

import (
	"fmt"
	"io"
	"sync"

	"gitoa.ru/go-4devs/config/definition"
)

var handlers = sync.Map{}

func Add(kind string, h Handle) error {
	_, ok := handlers.Load(kind)
	if ok {
		return fmt.Errorf("kind %v: %w", kind, ErrAlreadyExist)
	}

	handlers.Store(kind, h)
	return nil
}

func get(kind string) Handle {
	h, ok := handlers.Load(kind)
	if !ok {
		return func(w io.Writer, h Handler, o definition.Option) error {
			return fmt.Errorf("handler by %v:%w", kind, ErrNotFound)
		}
	}

	return h.(Handle)
}

func MustAdd(kind string, h Handle) {
	if err := Add(kind, h); err != nil {
		panic(err)
	}
}

type Handle func(io.Writer, Handler, definition.Option) error

type Handler interface {
	StructName() string
	Handle(io.Writer, Handler, definition.Option) error
	Options() ViewOption
	Keys() []string
	AddType(fullName string) (string, error)
	DefaultErrors() []string
}

type ViewOption struct {
	Prefix, Suffix string
	Context        bool
	Struct         string
	Errors         ViewErrors
}

type ViewErrors struct {
	Default []string
}
