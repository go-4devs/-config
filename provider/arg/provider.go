package arg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
	"gopkg.in/yaml.v3"
)

const Name = "arg"

var _ config.Provider = (*Provider)(nil)

type Option func(*Provider)

func WithKeyFactory(factory func(s ...string) string) Option {
	return func(p *Provider) { p.key = factory }
}

func New(opts ...Option) *Provider {
	prov := Provider{
		key: func(s ...string) string {
			return strings.Join(s, "-")
		},
		args: make(map[string][]string, len(os.Args[1:])),
		name: Name,
	}

	for _, opt := range opts {
		opt(&prov)
	}

	return &prov
}

type Provider struct {
	args map[string][]string
	key  func(...string) string
	name string
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
	return p.name
}

func (p *Provider) Value(ctx context.Context, path ...string) (config.Value, error) {
	if err := p.parse(); err != nil {
		return nil, err
	}

	name := p.key(path...)
	if val, ok := p.args[name]; ok {
		switch {
		case len(val) == 1:
			return value.JString(val[0]), nil
		default:
			var yNode yaml.Node

			if err := yaml.Unmarshal([]byte("["+strings.Join(val, ",")+"]"), &yNode); err != nil {
				return nil, fmt.Errorf("arg: failed unmarshal yaml:%w", err)
			}

			return value.Decode(yNode.Decode), nil
		}
	}

	return nil, fmt.Errorf("%s:%w", p.Name(), config.ErrValueNotFound)
}
