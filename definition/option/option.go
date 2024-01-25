package option

import (
	"gitoa.ru/go-4devs/config/definition"
)

var _ definition.Option = Option{}

const (
	Kind = "option"
)

const (
	TypeString   = "string"
	TypeInt      = "int"
	TypeInt64    = "int64"
	TypeUint     = "uint"
	TypeUint64   = "uint64"
	TypeFloat64  = "float64"
	TypeBool     = "bool"
	TypeTime     = "time.Time"
	TypeDuration = "time.Duration"
)

func Default(v any) func(*Option) {
	return func(o *Option) {
		o.Default = v
	}
}

func New(name, desc string, vtype any, opts ...func(*Option)) Option {
	option := Option{
		Name:        name,
		Description: desc,
		Type:        vtype,
	}

	for _, opt := range opts {
		opt(&option)
	}

	return option
}

type Option struct {
	Name        string
	Description string
	Type        any
	Default     any
	Params      definition.Params
}

func (o Option) WithParams(params ...definition.Param) Option {
	return Option{
		Name:        o.Name,
		Description: o.Description,
		Type:        o.Type,
		Params:      append(params, o.Params...),
	}
}

func (o Option) Kind() string {
	return Kind
}

func Time(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeTime, opts...)
}

func Duration(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeDuration, opts...)
}

func String(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeString, opts...)
}

func Int(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeInt, opts...)
}

func Int64(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeInt64, opts...)
}

func Uint(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeUint, opts...)
}

func Uint64(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeUint64, opts...)
}

func Float64(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeFloat64, opts...)
}

func Bool(name, desc string, opts ...func(*Option)) Option {
	return New(name, desc, TypeBool, opts...)
}
