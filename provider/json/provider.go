package json

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/tidwall/gjson"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
	"gitoa.ru/go-4devs/config/value"
)

var _ config.Provider = (*Provider)(nil)

func New(json []byte, opts ...Option) *Provider {
	provider := Provider{
		key:  key.Name,
		data: json,
	}

	for _, opt := range opts {
		opt(&provider)
	}

	return &provider
}

func NewFile(path string, opts ...Option) (*Provider, error) {
	file, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("%w: unable to read config file %#q: file not found or unreadable", err, path)
	}

	return New(file), nil
}

type Option func(*Provider)

type Provider struct {
	data []byte
	key  config.KeyFactory
}

func (p *Provider) IsSupport(ctx context.Context, key config.Key) bool {
	return p.key(ctx, key) != ""
}

func (p *Provider) Name() string {
	return "json"
}

func (p *Provider) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	path := p.key(ctx, key)

	if val := gjson.GetBytes(p.data, path); val.Exists() {
		return config.Variable{
			Name:     path,
			Provider: p.Name(),
			Value:    value.JString(val.String()),
		}, nil
	}

	return config.Variable{}, config.ErrVariableNotFound
}
