package watcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"gitoa.ru/go-4devs/config"
)

var (
	_ config.Provider      = (*Provider)(nil)
	_ config.WatchProvider = (*Provider)(nil)
)

func New(duration time.Duration, provider config.NamedProvider, opts ...Option) *Provider {
	prov := &Provider{
		NamedProvider: provider,
		ticker:        time.NewTicker(duration),
		logger: func(_ context.Context, msg string) {
			log.Print(msg)
		},
	}

	for _, opt := range opts {
		opt(prov)
	}

	return prov
}

func WithLogger(l func(context.Context, string)) Option {
	return func(p *Provider) {
		p.logger = l
	}
}

type Option func(*Provider)

type Provider struct {
	config.NamedProvider
	ticker *time.Ticker
	logger func(context.Context, string)
}

func (p *Provider) Watch(ctx context.Context, callback config.WatchCallback, key ...string) error {
	oldVar, err := p.NamedProvider.Value(ctx, key...)
	if err != nil {
		return fmt.Errorf("failed watch variable: %w", err)
	}

	go func() {
		for {
			select {
			case <-p.ticker.C:
				newVar, err := p.NamedProvider.Value(ctx, key...)
				if err != nil {
					p.logger(ctx, err.Error())
				} else if !newVar.IsEquals(oldVar) {
					callback(ctx, oldVar, newVar)
					oldVar = newVar
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
