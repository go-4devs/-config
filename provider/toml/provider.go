package toml

import (
	"context"
	"fmt"

	"github.com/pelletier/go-toml"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
	"gitoa.ru/go-4devs/config/value"
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
		key:  key.Name,
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
	key  config.KeyFactory
}

func (p *Provider) IsSupport(ctx context.Context, key config.Key) bool {
	return p.key(ctx, key) != ""
}

func (p *Provider) Name() string {
	return "toml"
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	if k := p.key(ctx, key); p.tree.Has(k) {
		return config.Variable{
			Name:     k,
			Provider: p.Name(),
			Value:    Value{Value: value.Value{Val: p.tree.Get(k)}},
		}, nil
	}

	return config.Variable{}, config.ErrVariableNotFound
}
