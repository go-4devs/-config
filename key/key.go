package key

import (
	"context"
	"strings"

	"gitoa.ru/go-4devs/config"
)

func NsAppName(sep string) config.KeyFactory {
	return func(_ context.Context, key config.Key) string {
		return strings.Join([]string{key.Namespace, key.AppName, key.Name}, sep)
	}
}

func AppName(sep string) config.KeyFactory {
	return func(_ context.Context, key config.Key) string {
		return strings.Join([]string{key.AppName, key.Name}, sep)
	}
}

func PrefixName(prefix string, factory config.KeyFactory) config.KeyFactory {
	return func(ctx context.Context, key config.Key) string {
		return prefix + factory(ctx, key)
	}
}

func Name(_ context.Context, key config.Key) string {
	return key.Name
}

func AliasName(name string, alias string, def config.KeyFactory) config.KeyFactory {
	return func(ctx context.Context, key config.Key) string {
		if name == key.Name {
			return alias
		}

		return def(ctx, key)
	}
}

func ReplaceAll(oldVal, newVal string, parent config.KeyFactory) config.KeyFactory {
	return func(ctx context.Context, key config.Key) string {
		return strings.ReplaceAll(parent(ctx, key), oldVal, newVal)
	}
}
