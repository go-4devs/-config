package env

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
)

const Name = "env"

var _ config.Provider = (*Provider)(nil)

type Option func(*Provider)

func WithKeyFactory(factory func(...string) string) Option {
	return func(p *Provider) { p.key = factory }
}

func New(namespace, appName string, opts ...Option) *Provider {
	provider := Provider{
		key: func(path ...string) string {
			return strings.ToUpper(strings.Join(path, "_"))
		},
		prefix: strings.ToUpper(namespace + "_" + appName + "_"),
	}

	for _, opt := range opts {
		opt(&provider)
	}

	return &provider
}

type Provider struct {
	key    func(...string) string
	name   string
	prefix string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Value(ctx context.Context, path ...string) (config.Value, error) {
	name := p.prefix + p.key(path...)
	if val, ok := os.LookupEnv(name); ok {
		return value.JString(val), nil
	}

	return nil, fmt.Errorf("%v:%w", p.Name(), config.ErrValueNotFound)
}
