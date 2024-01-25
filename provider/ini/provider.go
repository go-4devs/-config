package ini

import (
	"context"
	"fmt"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/ini.v1"
)

const (
	Name      = "ini"
	Separator = "."
)

var _ config.Provider = (*Provider)(nil)

func New(data *ini.File) *Provider {
	return &Provider{
		data: data,
		resolve: func(path []string) (string, string) {
			if len(path) == 1 {
				return "", path[0]
			}

			return strings.Join(path[:len(path)-1], Separator), strings.ToUpper(path[len(path)-1])
		},
		name: Name,
	}
}

type Provider struct {
	data    *ini.File
	resolve func(path []string) (string, string)
	name    string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Value(ctx context.Context, path ...string) (config.Value, error) {
	section, name := p.resolve(path)

	iniSection, err := p.data.GetSection(section)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", config.ErrValueNotFound, p.Name(), err)
	}

	iniKey, err := iniSection.GetKey(name)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %w", config.ErrValueNotFound, p.Name(), err)
	}

	return value.JString(iniKey.String()), nil
}
