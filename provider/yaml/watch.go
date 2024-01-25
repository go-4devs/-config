package yaml

import (
	"context"
	"fmt"
	"os"

	"gitoa.ru/go-4devs/config"
	"gopkg.in/yaml.v3"
)

const NameWatch = "yaml_watch"

func NewWatch(name string, opts ...Option) *Watch {
	f := Watch{
		file: name,
		prov: create(opts...),
		name: NameWatch,
	}

	return &f
}

type Watch struct {
	file string
	prov *Provider
	name string
}

func (p *Watch) Name() string {
	return p.name
}

func (p *Watch) Value(ctx context.Context, path ...string) (config.Value, error) {
	in, err := os.ReadFile(p.file)
	if err != nil {
		return nil, fmt.Errorf("yaml_file: read error: %w", err)
	}

	var yNode yaml.Node
	if err = yaml.Unmarshal(in, &yNode); err != nil {
		return nil, fmt.Errorf("yaml_file: unmarshal error: %w", err)
	}

	return p.prov.With(&yNode).Value(ctx, path...)
}
