package arg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/yaml.v3"
)

var _ config.Provider = (*Provider)(nil)

type Option func(*Provider)

func WithKeyFactory(factory config.KeyFactory) Option {
	return func(p *Provider) { p.key = factory }
}

func New(opts ...Option) *Provider {
	prov := Provider{
		key:  key.Name,
		args: make(map[string][]string, len(os.Args[1:])),
	}

	for _, opt := range opts {
		opt(&prov)
	}

	return &prov
}

type Provider struct {
	args map[string][]string
	key  config.KeyFactory
}

// nolint: cyclop
// return name, value, error.
func (p *Provider) parseOne(arg string) (string, string, error) {
	if arg[0] != '-' {
		return "", "", nil
	}

	numMinuses := 1

	if arg[1] == '-' {
		numMinuses++
	}

	name := strings.TrimSpace(arg[numMinuses:])
	if len(name) == 0 {
		return name, "", nil
	}

	if name[0] == '-' || name[0] == '=' {
		return "", "", fmt.Errorf("%w: bad flag syntax: %s", config.ErrInvalidValue, arg)
	}

	var val string

	for idx := 1; idx < len(name); idx++ {
		if name[idx] == '=' || name[idx] == ' ' {
			val = strings.TrimSpace(name[idx+1:])
			name = name[0:idx]

			break
		}
	}

	if val == "" && numMinuses == 1 && len(arg) > 2 {
		name, val = name[:1], name[1:]
	}

	return name, val, nil
}

func (p *Provider) parse() error {
	if len(p.args) > 0 {
		return nil
	}

	for _, arg := range os.Args[1:] {
		name, value, err := p.parseOne(arg)
		if err != nil {
			return err
		}

		if name != "" {
			p.args[name] = append(p.args[name], value)
		}
	}

	return nil
}

func (p *Provider) Name() string {
	return "arg"
}

func (p *Provider) IsSupport(ctx context.Context, key config.Key) bool {
	return p.key(ctx, key) != ""
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	if err := p.parse(); err != nil {
		return config.Variable{
			Name:     "",
			Value:    nil,
			Provider: p.Name(),
		}, err
	}

	name := p.key(ctx, key)
	if val, ok := p.args[name]; ok {
		switch {
		case len(val) == 1:
			return config.Variable{
				Name:     name,
				Provider: p.Name(),
				Value:    value.JString(val[0]),
			}, nil
		default:
			var yNode yaml.Node

			if err := yaml.Unmarshal([]byte("["+strings.Join(val, ",")+"]"), &yNode); err != nil {
				return config.Variable{}, fmt.Errorf("arg: failed unmarshal yaml:%w", err)
			}

			return config.Variable{
				Name:     name,
				Provider: p.Name(),
				Value:    value.Decode(yNode.Decode),
			}, nil
		}
	}

	return config.Variable{}, fmt.Errorf("%w: %s", config.ErrVariableNotFound, name)
}
