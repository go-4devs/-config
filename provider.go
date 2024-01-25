package config

import "context"

type Provider interface {
	Value(ctx context.Context, path ...string) (Value, error)
}

type NamedProvider interface {
	Name() string
	Provider
}

type WatchCallback func(ctx context.Context, oldVar, newVar Value)

type WatchProvider interface {
	Watch(ctx context.Context, callback WatchCallback, path ...string) error
}

type Factory func(ctx context.Context, cfg Provider) (Provider, error)
