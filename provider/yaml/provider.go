package yaml

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/yaml.v3"
)

const (
	Name = "yaml"
)

var _ config.Provider = (*Provider)(nil)

func NewFile(name string, opts ...Option) (*Provider, error) {
	in, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("yaml_file: read error: %w", err)
	}

	return New(in, opts...)
}

func New(yml []byte, opts ...Option) (*Provider, error) {
	var data yaml.Node
	if err := yaml.Unmarshal(yml, &data); err != nil {
		return nil, fmt.Errorf("yaml: unmarshal err: %w", err)
	}

	return create(opts...).With(&data), nil
}

func create(opts ...Option) *Provider {
	prov := Provider{
		name: Name,
	}

	for _, opt := range opts {
		opt(&prov)
	}

	return &prov
}

type Option func(*Provider)

type Provider struct {
	data node
	name string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Value(_ context.Context, path ...string) (config.Value, error) {
	return p.data.read(p.Name(), path)
}

func (p *Provider) With(data *yaml.Node) *Provider {
	return &Provider{
		data: node{Node: data},
	}
}

type node struct {
	*yaml.Node
}

func (n *node) read(name string, keys []string) (config.Value, error) {
	val, err := getData(n.Node.Content[0].Content, keys)
	if err != nil {
		if errors.Is(err, config.ErrValueNotFound) {
			return nil, fmt.Errorf("%w: %s", config.ErrValueNotFound, name)
		}

		return nil, fmt.Errorf("%w: %s", err, name)
	}

	return value.Decode(val), nil
}

func getData(node []*yaml.Node, keys []string) (func(interface{}) error, error) {
	for idx := len(node) - 1; idx > 0; idx -= 2 {
		if node[idx-1].Value == keys[0] {
			if len(keys) > 1 {
				return getData(node[idx].Content, keys[1:])
			}

			return node[idx].Decode, nil
		}
	}

	return nil, config.ErrValueNotFound
}
