package key_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/key"
)

func TestLastIndex(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cases := map[string]struct {
		sep   string
		path  string
		field string
	}{
		"/secret/with/field/name": {
			sep:   "/",
			path:  "/secret/with/field",
			field: "name",
		},
		"/secret/database:username": {
			sep:   ":",
			path:  "/secret/database",
			field: "username",
		},
		"database:username": {
			sep:   ":",
			path:  "database",
			field: "username",
		},
		"/secret/database-dsn": {
			sep:   ":",
			path:  "/secret/database-dsn",
			field: "",
		},
		"/secret/database--dsn": {
			sep:   "--",
			path:  "/secret/database",
			field: "dsn",
		},
		"/secret/database:dsn": {
			sep:   "--",
			path:  "/secret/database:dsn",
			field: "",
		},
	}

	for path, data := range cases {
		k := config.Key{
			Name: path,
		}

		fn := key.LastIndex(data.sep, key.Name)
		ns, field := fn(ctx, k)
		assert.Equal(t, data.field, field, k)
		assert.Equal(t, data.path, ns, k)
	}
}
