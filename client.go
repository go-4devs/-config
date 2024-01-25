package config

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

func Must(providers ...interface{}) *Client {
	client, err := New(providers...)
	if err != nil {
		panic(err)
	}

	return client
}

func New(providers ...interface{}) (*Client, error) {
	client := &Client{
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

func (p *provider) Watch(ctx context.Context, callback WatchCallback, path ...string) error {
	if err := p.init(ctx); err != nil {
		return fmt.Errorf("init read:%w", err)
	}

	watch, ok := p.provider.(WatchProvider)
	if !ok {
		return nil
	}

	if err := watch.Watch(ctx, callback, path...); err != nil {
		return fmt.Errorf("factory provider: %w", err)
	}

	return nil
}

func (p *provider) Value(ctx context.Context, path ...string) (Value, error) {
	if err := p.init(ctx); err != nil {
		return nil, fmt.Errorf("init read:%w", err)
	}

	variable, err := p.provider.Value(ctx, path...)
	if err != nil {
		return nil, fmt.Errorf("factory provider: %w", err)
	}

	return variable, nil
}

type Client struct {
	providers []Provider
}

func (c *Client) Name() string {
	return "client"
}

// Value get value by name.
func (c *Client) Value(ctx context.Context, path ...string) (Value, error) {
	var (
		value Value
		err   error
	)

	for _, provider := range c.providers {
		value, err = provider.Value(ctx, path...)
		if err == nil || !(errors.Is(err, ErrValueNotFound) || errors.Is(err, ErrInitFactory)) {
			break
		}
	}

	if err != nil {
		return value, fmt.Errorf("client failed get value: %w", err)
	}

	return value, nil
}

func (c *Client) Watch(ctx context.Context, callback WatchCallback, path ...string) error {
	for idx, prov := range c.providers {
		provider, ok := prov.(WatchProvider)
		if !ok {
			continue
		}

		err := provider.Watch(ctx, callback, path...)
		if err != nil {
			if errors.Is(err, ErrValueNotFound) || errors.Is(err, ErrInitFactory) {
				continue
			}

			return fmt.Errorf("client: failed watch by provider[%d]: %w", idx, err)
		}
	}

	return nil
}
