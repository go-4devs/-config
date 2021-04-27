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
	p := Provider{
		key: key.Name,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return &p
}

type Provider struct {
	args map[string][]string
	key  config.KeyFactory
}

// nolint: cyclop
func (p *Provider) parseOne(arg string) (name, val string, err error) {
	if arg[0] != '-' {
		return
	}

	numMinuses := 1

	if arg[1] == '-' {
		numMinuses++
	}

	name = strings.TrimSpace(arg[numMinuses:])
	if len(name) == 0 {
		return
	}

	if name[0] == '-' || name[0] == '=' {
		return "", "", fmt.Errorf("%w: bad flag syntax: %s", config.ErrInvalidValue, arg)
	}

	for i := 1; i < len(name); i++ {
		if name[i] == '=' || name[i] == ' ' {
			val = strings.TrimSpace(name[i+1:])
			name = name[0:i]

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

	p.args = make(map[string][]string, len(os.Args[1:]))

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
		return config.Variable{Provider: p.Name()}, err
	}

	k := p.key(ctx, key)
	if val, ok := p.args[k]; ok {
		switch {
		case len(val) == 1:
			return config.Variable{
				Name:     k,
				Provider: p.Name(),
				Value:    value.JString(val[0]),
			}, nil
		default:
			var n yaml.Node

			if err := yaml.Unmarshal([]byte("["+strings.Join(val, ",")+"]"), &n); err != nil {
				return config.Variable{}, fmt.Errorf("arg: failed unmarshal yaml:%w", err)
			}

			return config.Variable{
				Name:     k,
				Provider: p.Name(),
				Value:    value.Decode(n.Decode),
			}, nil
		}
	}

	return config.Variable{Name: k, Provider: p.Name()}, config.ErrVariableNotFound
}
