package vault_test

import (
	"testing"
	"time"

	"gitoa.ru/go-4devs/config/provider/vault"
	"gitoa.ru/go-4devs/config/test"
	"gitoa.ru/go-4devs/config/test/require"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	cl, err := NewVault()
	require.NoError(t, err)

	provider := vault.New("fdevs", "config", cl)

	read := []test.Read{
		test.NewReadConfig("database"),
		test.NewRead(test.DSN, "db", "dsn"),
		test.NewRead(time.Minute, "db", "timeout"),
	}
	test.Run(t, provider, read)
}
