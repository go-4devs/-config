package group

import (
	"gitoa.ru/go-4devs/config/definition"
)

const Kind = "group"

var _ definition.Option = Group{}

func New(name, desc string, opts ...definition.Option) Group {
	return Group{
		Name:        name,
		Description: desc,
		Options:     opts,
	}
}

type Group struct {
	Options     definition.Options
	Name        string
	Description string
}

func (o Group) Kind() string {
	return Kind
}
