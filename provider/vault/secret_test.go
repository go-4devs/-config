package vault_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gitoa.ru/go-4devs/config/provider/vault"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	cl, err := test.NewVault()
	require.NoError(t, err)

	provider := vault.NewSecretKV2(cl)

	read := []test.Read{
		test.NewReadConfig("database"),
		test.NewRead("db:dsn", test.DSN),
		test.NewRead("db:timeout", time.Minute),
	}
	test.Run(t, provider, read)
}
