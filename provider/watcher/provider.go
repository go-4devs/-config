package watcher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gitoa.ru/go-4devs/config"
)

var (
	_ config.Provider      = (*Provider)(nil)
	_ config.WatchProvider = (*Provider)(nil)
)

func New(duration time.Duration, provider config.Provider, opts ...Option) *Provider {
	prov := &Provider{
		Provider: provider,
		duration: duration,
		logger:   slog.ErrorContext,
	}

	for _, opt := range opts {
		opt(prov)
	}

	return prov
}

func WithLogger(l func(context.Context, string, ...any)) Option {
	return func(p *Provider) {
		p.logger = l
	}
}

type Option func(*Provider)

type Provider struct {
	config.Provider
	duration time.Duration
	logger   func(context.Context, string, ...any)
}

func (p *Provider) Watch(ctx context.Context, callback config.WatchCallback, key ...string) error {
	old, err := p.Provider.Value(ctx, key...)
	if err != nil {
		return fmt.Errorf("failed watch variable: %w", err)
	}

	go func(oldVar config.Value) {
		ticker := time.NewTicker(p.duration)
		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-ticker.C:
				newVar, err := p.Provider.Value(ctx, key...)
				if err != nil {
					p.logger(ctx, "get value%v:%v", key, err.Error())
				} else if !newVar.IsEquals(oldVar) {
					if err := callback(ctx, oldVar, newVar); err != nil {
						if errors.Is(err, config.ErrStopWatch) {
							return
						}
						p.logger(ctx, "callback %v:%v", key, err)

					}
					oldVar = newVar
				}
			case <-ctx.Done():
				return
			}
		}
	}(old)

	return nil
}
