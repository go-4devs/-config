package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/vault/api"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/value"
)

const (
	Name      = "vault"
	Separator = "/"
	Prefix    = "secret/data/"
	ValueName = "value"
)

var _ config.Provider = (*SecretKV2)(nil)

type SecretOption func(*SecretKV2)

func WithSecretResolve(f func(key []string) (string, string)) SecretOption {
	return func(s *SecretKV2) { s.resolve = f }
}

func NewSecretKV2(namespace, appName string, client *api.Client, opts ...SecretOption) *SecretKV2 {
	prov := SecretKV2{
		client: client,
		resolve: func(key []string) (string, string) {
			keysLen := len(key)
			if keysLen == 1 {
				return "", key[0]
			}

			return strings.Join(key[:keysLen-1], Separator), key[keysLen-1]
		},
		name:   Name,
		prefix: Prefix + namespace + Separator + appName,
	}

	for _, opt := range opts {
		opt(&prov)
	}

	return &prov
}

type SecretKV2 struct {
	client  *api.Client
	resolve func(key []string) (string, string)
	name    string
	prefix  string
}

func (p *SecretKV2) Name() string {
	return p.name
}
func (p *SecretKV2) Key(in []string) (string, string) {
	path, val := p.resolve(in)
	if path == "" {
		return p.prefix, val
	}

	return p.prefix + Separator + path, val
}
func (p *SecretKV2) read(path, key string) (*api.Secret, error) {
	secret, err := p.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}
	if secret == nil && key != ValueName {
		return p.read(path+Separator+key, ValueName)
	}

	return secret, err
}

func (p *SecretKV2) Value(ctx context.Context, key ...string) (config.Value, error) {
	path, field := p.Key(key)

	secret, err := p.read(path, field)
	if err != nil {
		return nil, fmt.Errorf("%w: path:%s, field:%s, provider:%s", err, path, field, p.Name())
	}

	if secret == nil || len(secret.Data) == 0 {
		log.Println(secret == nil)
		return nil, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrValueNotFound, path, field, p.Name())
	}

	if len(secret.Warnings) > 0 {
		return nil,
			fmt.Errorf("%w: warn: %s, path:%s, field:%s, provider:%s", config.ErrValueNotFound, secret.Warnings, path, field, p.Name())
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: path:%s, field:%s, provider:%s", config.ErrValueNotFound, path, field, p.Name())
	}

	if val, ok := data[field]; ok {
		return value.JString(fmt.Sprint(val)), nil
	}

	if val, ok := data[ValueName]; ok {
		return value.JString(fmt.Sprint(val)), nil
	}

	md, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", config.ErrInvalidValue, err)
	}

	return value.JBytes(md), nil
}
