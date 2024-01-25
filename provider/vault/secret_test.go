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

	provider := vault.NewSecretKV2("fdevs", "config", cl)

	read := []test.Read{
		test.NewReadConfig("database"),
		test.NewRead(test.DSN, "db", "dsn"),
		test.NewRead(time.Minute, "db", "timeout"),
	}
	test.Run(t, provider, read)
}
