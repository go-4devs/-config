package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/test/require"
)

const (
	DSN       = "pgsql://user@pass:127.0.0.1:5432"
	Namespace = "fdevs"
	AppName   = "config"
)

func Run(t *testing.T, provider config.Provider, read []Read) {
	t.Helper()

	ctx := context.Background()

	for idx, read := range read {
		t.Run(fmt.Sprintf("%v:%v", idx, read.Key), func(t *testing.T) {
			val, err := provider.Value(ctx, read.Key...)
			require.NoError(t, err, read.Key)
			read.Assert(t, val)
		})
	}
}

type Read struct {
	Key    []string
	Assert func(t *testing.T, v config.Value)
}

type Config struct {
	Duration time.Duration
	Enabled  bool
}

func NewReadConfig(key ...string) Read {
	ex := &Config{
		Duration: 21 * time.Minute,
		Enabled:  true,
	}

	return NewReadUnmarshal(ex, &Config{}, key...)
}

func NewReadUnmarshal(expected, target interface{}, key ...string) Read {
	return Read{
		Key: key,
		Assert: func(t *testing.T, v config.Value) {
			t.Helper()
			require.NoErrorf(t, v.Unmarshal(target), "unmarshal")
			require.Equal(t, expected, target, "unmarshal")
		},
	}
}

func Time(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)

	return t
}

// nolint: cyclop
func NewRead(expected interface{}, key ...string) Read {
	return Read{
		Key: key,
		Assert: func(t *testing.T, v config.Value) {
			t.Helper()
			var (
				val   interface{}
				err   error
				short interface{}
			)
			switch expected.(type) {
			case bool:
				val, err = v.ParseBool()
				short = v.Bool()
			case int:
				val, err = v.ParseInt()
				short = v.Int()
			case int64:
				val, err = v.ParseInt64()
				short = v.Int64()
			case uint:
				val, err = v.ParseUint()
				short = v.Uint()
			case uint64:
				val, err = v.ParseUint64()
				short = v.Uint64()
			case string:
				val, err = v.ParseString()
				short = v.String()
			case float64:
				val, err = v.ParseFloat64()
				short = v.Float64()
			case time.Duration:
				val, err = v.ParseDuration()
				short = v.Duration()
			case time.Time:
				val, err = v.ParseTime()
				short = v.Time()
			default:
				require.Fail(t, "unexpected type:%+T", expected)
			}

			require.Equalf(t, val, short, "type:%T", expected)
			require.NoErrorf(t, err, "type:%T", expected)
			require.Equalf(t, expected, val, "type:%T", expected)
		},
	}
}
