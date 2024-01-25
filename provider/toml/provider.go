package toml

import (
	"context"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
)

const (
	Name      = "toml"
	Separator = "."
)

var _ config.Provider = (*Provider)(nil)

func NewFile(file string, opts ...Option) (*Provider, error) {
	tree, err := toml.LoadFile(file)
	if err != nil {
		return nil, fmt.Errorf("toml: failed load file: %w", err)
	}

	return configure(tree, opts...), nil
}

type Option func(*Provider)

func configure(tree *toml.Tree, opts ...Option) *Provider {
	prov := &Provider{
		tree: tree,
		key: func(s []string) string {
			return strings.Join(s, Separator)
		},
	}

	for _, opt := range opts {
		opt(prov)
	}

	return prov
}

func New(data []byte, opts ...Option) (*Provider, error) {
	tree, err := toml.LoadBytes(data)
	if err != nil {
		return nil, fmt.Errorf("toml failed load data: %w", err)
	}

	return configure(tree, opts...), nil
}

type Provider struct {
	tree *toml.Tree
	key  func([]string) string
	name string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Value(_ context.Context, path ...string) (config.Value, error) {
	if k := p.key(path); p.tree.Has(k) {
		return Value{Value: value.Value{Val: p.tree.Get(k)}}, nil
	}

	return nil, config.ErrValueNotFound
}
