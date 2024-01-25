package value_test

import (
	"errors"
	"testing"
	"time"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/test/require"
	"gitoa.ru/go-4devs/config/value"
)

func TestStringDuration(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		raw value.JString
		exp time.Duration
		err error
	}{
		"1m": {
			raw: value.JString("1m"),
			exp: time.Minute,
		},
		"number error": {
			raw: value.JString("100000000"),
			err: config.ErrInvalidValue,
		},
	}

	for name, data := range tests {
		require.Equal(t, data.exp, data.raw.Duration(), name)
		d, err := data.raw.ParseDuration()
		require.Truef(t, errors.Is(err, data.err), "%[1]s: expect:%#[2]v, got:%#[3]v", name, data.err, err)
		require.Equal(t, data.exp, d, name)
	}
}
