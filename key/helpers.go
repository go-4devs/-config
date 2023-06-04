package key

import (
	"context"
	"strings"

	"gitoa.ru/go-4devs/config"
)

func LastIndex(sep string, factory config.KeyFactory) func(ctx context.Context, key config.Key) (string, string) {
	return func(ctx context.Context, key config.Key) (string, string) {
		name := factory(ctx, key)

		idx := strings.LastIndex(name, sep)
		if idx == -1 {
			return name, ""
		}

		return name[0:idx], name[idx+len(sep):]
	}
}

func LastIndexField(sep, def string, factory config.KeyFactory) func(ctx context.Context, key config.Key) (string, string) {
	return func(ctx context.Context, key config.Key) (string, string) {
		p, k := LastIndex(sep, factory)(ctx, key)
		if k == "" {
			return p, def
		}

		return p, k
	}
}
