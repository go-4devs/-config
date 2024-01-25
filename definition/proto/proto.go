package proto

import (
	"gitoa.ru/go-4devs/config/definition"
)

const Kind = "proto"

func New(name, desc string, opt definition.Option) Proto {
	pr := Proto{
		Name:        name,
		Description: desc,
		Option:      opt,
	}

	return pr
}

type Proto struct {
	Name        string
	Description string
	Option      definition.Option
}

func (p Proto) Kind() string {
	return Kind
}
