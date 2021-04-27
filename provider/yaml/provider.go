package yaml

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/yaml.v3"
)

var _ config.Provider = (*Provider)(nil)

func keyFactory(ctx context.Context, key config.Key) []string {
	return strings.Split(key.Name, "/")
}

func New(yml []byte, opts ...Option) (*Provider, error) {
	var data yaml.Node
	if err := yaml.Unmarshal(yml, &data); err != nil {
		return nil, fmt.Errorf("yaml: unmarshal err: %w", err)
	}

	p := Provider{
		key:  keyFactory,
		data: node{Node: &data},
	}

	for _, opt := range opts {
		opt(&p)
	}

	return &p, nil
}

type Option func(*Provider)

type Provider struct {
	data node
	key  func(context.Context, config.Key) []string
}

func (p *Provider) Name() string {
	return "yaml"
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	k := p.key(ctx, key)

	return p.data.read(p.Name(), k)
}

type node struct {
	*yaml.Node
}

func (n *node) read(name string, k []string) (config.Variable, error) {
	val, err := getData(n.Node.Content[0].Content, k)
	if err != nil {
		if errors.Is(err, config.ErrVariableNotFound) {
			return config.Variable{}, fmt.Errorf("%w: %s", config.ErrVariableNotFound, name)
		}

		return config.Variable{}, fmt.Errorf("%w: %s", err, name)
	}

	return config.Variable{
		Name:     strings.Join(k, "."),
		Provider: name,
		Value:    value.Decode(val),
	}, nil
}

func getData(node []*yaml.Node, keys []string) (func(interface{}) error, error) {
	for i := len(node) - 1; i > 0; i -= 2 {
		if node[i-1].Value == keys[0] {
			if len(keys) > 1 {
				return getData(node[i].Content, keys[1:])
			}

			return node[i].Decode, nil
		}
	}

	return nil, config.ErrVariableNotFound
}
