package vault

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
	"gitoa.ru/go-4devs/config/value"
)

var _ config.Provider = (*SecretKV2)(nil)

type SecretOption func(*SecretKV2)

func WithSecretResolve(f func(context.Context, config.Key) (string, string)) SecretOption {
	return func(s *SecretKV2) { s.resolve = f }
}

func NewSecretKV2(client *api.Client, opts ...SecretOption) *SecretKV2 {
	s := SecretKV2{
		client:  client,
		resolve: key.LastIndexField(":", "value", key.PrefixName("secret/data/", key.NsAppName("/"))),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

type SecretKV2 struct {
	client  *api.Client
	resolve func(ctx context.Context, key config.Key) (string, string)
}

func (p *SecretKV2) IsSupport(ctx context.Context, key config.Key) bool {
	path, _ := p.resolve(ctx, key)

	return path != ""
}

func (p *SecretKV2) Name() string {
	return "vault"
}

func (p *SecretKV2) Read(ctx context.Context, key config.Key) (config.Variable, error) {
	path, field := p.resolve(ctx, key)

	s, err := p.client.Logical().Read(path)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", err, path, field, p.Name())
	}

	if s == nil || len(s.Data) == 0 {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrVariableNotFound, path, field, p.Name())
	}

	if len(s.Warnings) > 0 {
		return config.Variable{},
			fmt.Errorf("%w: warn: %s, path:%s, field:%s, provider:%s", config.ErrVariableNotFound, s.Warnings, path, field, p.Name())
	}

	d, ok := s.Data["data"].(map[string]interface{})
	if !ok {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrVariableNotFound, path, field, p.Name())
	}

	if val, ok := d[field]; ok {
		return config.Variable{
			Name:     path + field,
			Provider: p.Name(),
			Value:    value.JString(fmt.Sprint(val)),
		}, nil
	}

	md, err := json.Marshal(d)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: %s", config.ErrInvalidValue, err)
	}

	return config.Variable{
		Name:     path + field,
		Provider: p.Name(),
		Value:    value.JBytes(md),
	}, nil
}
