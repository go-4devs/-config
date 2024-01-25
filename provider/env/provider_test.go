package env_test

import (
	"os"
	"testing"

	"gitoa.ru/go-4devs/config/provider/env"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	os.Setenv("FDEVS_CONFIG_DSN", test.DSN)
	os.Setenv("FDEVS_CONFIG_PORT", "8080")

	provider := env.New("fdevs", "config")

	read := []test.Read{
		test.NewRead(test.DSN, "dsn"),
		test.NewRead(8080, "port"),
	}
	test.Run(t, provider, read)
}
