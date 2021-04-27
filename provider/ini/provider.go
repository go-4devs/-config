package ini

import (
	"context"
	"fmt"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/ini.v1"
)

var _ config.Provider = (*Provider)(nil)

func New(data *ini.File) *Provider {
	return &Provider{
		data: data,
		resolve: func(ctx context.Context, key config.Key) (string, string) {
			keys := strings.SplitN(key.Name, "/", 2)
			if len(keys) == 1 {
				return "", keys[0]
			}

			return keys[0], keys[1]
		},
	}
}

type Provider struct {
	data    *ini.File
	resolve func(ctx context.Context, key config.Key) (string, string)
}

func (p *Provider) IsSupport(ctx context.Context, key config.Key) bool {
	section, name := p.resolve(ctx, key)

	return section != "" && name != ""
}

func (p *Provider) Name() string {
	return "ini"
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	section, name := p.resolve(ctx, key)

	iniSection, err := p.data.GetSection(section)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: %s: %v", config.ErrVariableNotFound, p.Name(), err)
	}

	iniKey, err := iniSection.GetKey(name)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: %s: %v", config.ErrVariableNotFound, p.Name(), err)
	}

	return config.Variable{
		Name:     section + ":" + name,
		Provider: p.Name(),
		Value:    value.JString(iniKey.String()),
	}, nil
}
