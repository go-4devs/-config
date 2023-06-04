package yaml

import (
	"context"
	"fmt"
	"io/ioutil"

	"gitoa.ru/go-4devs/config"
	"gopkg.in/yaml.v3"
)

func NewWatch(name string, opts ...Option) *Watch {
	f := Watch{
		file: name,
		prov: create(opts...),
	}

	return &f
}

type Watch struct {
	file string
	prov *Provider
}

func (p *Watch) Name() string {
	return "yaml_watch"
}

func (p *Watch) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	in, err := ioutil.ReadFile(p.file)
	if err != nil {
		return config.Variable{}, fmt.Errorf("yaml_file: read error: %w", err)
	}

	var n yaml.Node
	if err = yaml.Unmarshal(in, &n); err != nil {
		return config.Variable{}, fmt.Errorf("yaml_file: unmarshal error: %w", err)
	}

	return p.prov.With(&n).Read(ctx, key)
}
