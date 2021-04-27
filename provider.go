package config

import "context"

type Provider interface {
	Read(ctx context.Context, key Key) (Variable, error)
	NamedProvider
}

type WatchCallback func(ctx context.Context, oldVar, newVar Variable)

type WatchProvider interface {
	Watch(ctx context.Context, key Key, callback WatchCallback) error
	NamedProvider
}

type NamedProvider interface {
	Name() string
}
