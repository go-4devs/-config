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
	prov := SecretKV2{
		client:  client,
		resolve: key.LastIndexField(":", "value", key.PrefixName("secret/data/", key.NsAppName("/"))),
	}

	for _, opt := range opts {
		opt(&prov)
	}

	return &prov
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

	secret, err := p.client.Logical().Read(path)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", err, path, field, p.Name())
	}

	if secret == nil || len(secret.Data) == 0 {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrVariableNotFound, path, field, p.Name())
	}

	if len(secret.Warnings) > 0 {
		return config.Variable{},
			fmt.Errorf("%w: warn: %s, path:%s, field:%s, provider:%s", config.ErrVariableNotFound, secret.Warnings, path, field, p.Name())
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return config.Variable{}, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrVariableNotFound, path, field, p.Name())
	}

	if val, ok := data[field]; ok {
		return config.Variable{
			Name:     path + field,
			Provider: p.Name(),
			Value:    value.JString(fmt.Sprint(val)),
		}, nil
	}

	md, err := json.Marshal(data)
	if err != nil {
		return config.Variable{}, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return config.Variable{
		Name:     path + field,
		Provider: p.Name(),
		Value:    value.JBytes(md),
	}, nil
}
