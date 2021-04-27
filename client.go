package config

import (
	"context"
	"errors"
	"fmt"
)

func New(namespace, appName string, providers []Provider) *Client {
	return &Client{
		namespace: namespace,
		appName:   appName,
		providers: providers,
	}
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

func (c *Client) Value(ctx context.Context, name string) (Value, error) {
	variable, err := c.Variable(ctx, name)
	if err != nil {
		return nil, err
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
		if err == nil || !errors.Is(err, ErrVariableNotFound) {
			break
		}
	}

	if err != nil {
		return variable, fmt.Errorf("client failed get variable: %w", err)
	}

	return variable, nil
}

func NewWatch(namespace, appName string, providers []Provider) *WatchClient {
	cl := WatchClient{
		Client: New(namespace, appName, providers),
	}

	for _, provider := range providers {
		if watch, ok := provider.(WatchProvider); ok {
			cl.providers = append(cl.providers, watch)
		}
	}

	return &cl
}

type WatchClient struct {
	*Client
	providers []WatchProvider
}

func (wc *WatchClient) Watch(ctx context.Context, name string, callback WatchCallback) error {
	key := wc.key(name)

	for _, provider := range wc.providers {
		err := provider.Watch(ctx, key, callback)
		if err != nil {
			if errors.Is(err, ErrVariableNotFound) {
				continue
			}

			return fmt.Errorf("client: failed watch by provider %s: %w", provider.Name(), err)
		}
	}

	return nil
}
