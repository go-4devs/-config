package config

import "context"

type Provider interface {
	Read(ctx context.Context, key Key) (Variable, error)
}

type WatchCallback func(ctx context.Context, oldVar, newVar Variable)

type WatchProvider interface {
	Watch(ctx context.Context, key Key, callback WatchCallback) error
}

type ReadConfig interface {
	Value(ctx context.Context, name string) (Value, error)
}

type Factory func(ctx context.Context, cfg ReadConfig) (Provider, error)
