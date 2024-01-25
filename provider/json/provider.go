package json

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
)

const (
	Name      = "json"
	Separator = "."
)

var _ config.Provider = (*Provider)(nil)

func New(json []byte, opts ...Option) *Provider {
	provider := Provider{
		key: func(s ...string) string {
			return strings.Join(s, Separator)
		},
		data: json,
	}

	for _, opt := range opts {
		opt(&provider)
	}

	return &provider
}

func NewFile(path string, opts ...Option) (*Provider, error) {
	file, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read config file %#q: file not found or unreadable", err, path)
	}

	return New(file, opts...), nil
}

type Option func(*Provider)

type Provider struct {
	data []byte
	key  func(...string) string
	name string
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Value(_ context.Context, path ...string) (config.Value, error) {
	key := p.key(path...)
	if val := gjson.GetBytes(p.data, key); val.Exists() {
		return value.JString(val.String()), nil
	}

	return nil, fmt.Errorf("%v:%w", p.Name(), config.ErrValueNotFound)
}
