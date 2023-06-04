package config

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

func Must(namespace, appName string, providers ...interface{}) *Client {
	client, err := New(namespace, appName, providers...)
	if err != nil {
		panic(err)
	}

	return client
}

func New(namespace, appName string, providers ...interface{}) (*Client, error) {
	client := &Client{
		namespace: namespace,
		appName:   appName,
		providers: make([]Provider, len(providers)),
	}

	for idx, prov := range providers {
		switch current := prov.(type) {
		case Provider:
			client.providers[idx] = current
		case Factory:
			client.providers[idx] = &provider{
				factory: func(ctx context.Context) (Provider, error) {
					return current(ctx, client)
				},
			}
		default:
			return nil, fmt.Errorf("provier[%d]: %w %T", idx, ErrUnknowType, prov)
		}
	}

	return client, nil
}

type provider struct {
	mu       sync.Mutex
	done     uint32
	provider Provider
	factory  func(ctx context.Context) (Provider, error)
}

func (p *provider) init(ctx context.Context) error {
	if atomic.LoadUint32(&p.done) == 0 {
		if !p.mu.TryLock() {
			return fmt.Errorf("%w", ErrInitFactory)
		}
		defer atomic.StoreUint32(&p.done, 1)
		defer p.mu.Unlock()

		var err error
		if p.provider, err = p.factory(ctx); err != nil {
			return fmt.Errorf("init provider factory:%w", err)
		}
	}

	return nil
}

func (p *provider) Watch(ctx context.Context, key Key, callback WatchCallback) error {
	if err := p.init(ctx); err != nil {
		return fmt.Errorf("init read:%w", err)
	}

	watch, ok := p.provider.(WatchProvider)
	if !ok {
		return nil
	}

	if err := watch.Watch(ctx, key, callback); err != nil {
		return fmt.Errorf("factory provider: %w", err)
	}

	return nil
}

func (p *provider) Read(ctx context.Context, key Key) (Variable, error) {
	if err := p.init(ctx); err != nil {
		return Variable{}, fmt.Errorf("init read:%w", err)
	}

	variable, err := p.provider.Read(ctx, key)
	if err != nil {
		return Variable{}, fmt.Errorf("factory provider: %w", err)
	}

	return variable, nil
}

type Client struct {
	providers []Provider
	appName   string
	namespace string
}

func (c *Client) key(name string) Key {
	return Key{
		Name:      name,
		AppName:   c.appName,
		Namespace: c.namespace,
	}
}

// Value get value by name.
// nolint: ireturn
func (c *Client) Value(ctx context.Context, name string) (Value, error) {
	variable, err := c.Variable(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("variable:%w", err)
	}

	return variable.Value, nil
}

func (c *Client) Variable(ctx context.Context, name string) (Variable, error) {
	var (
		variable Variable
		err      error
	)

	key := c.key(name)

	for _, provider := range c.providers {
		variable, err = provider.Read(ctx, key)
		if err == nil || !(errors.Is(err, ErrVariableNotFound) || errors.Is(err, ErrInitFactory)) {
			break
		}
	}

	if err != nil {
		return variable, fmt.Errorf("client failed get variable: %w", err)
	}

	return variable, nil
}

func (c *Client) Watch(ctx context.Context, name string, callback WatchCallback) error {
	key := c.key(name)

	for idx, prov := range c.providers {
		provider, ok := prov.(WatchProvider)
		if !ok {
			continue
		}

		err := provider.Watch(ctx, key, callback)
		if err != nil {
			if errors.Is(err, ErrVariableNotFound) || errors.Is(err, ErrInitFactory) {
				continue
			}

			return fmt.Errorf("client: failed watch by provider[%d]: %w", idx, err)
		}
	}

	return nil
}
