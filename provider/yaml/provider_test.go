package yaml_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	provider "gitoa.ru/go-4devs/config/provider/yaml"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	prov, err := provider.New(test.ReadFile("config.yaml"))
	require.Nil(t, err)

	read := []test.Read{
		test.NewRead("duration_var", 21*time.Minute),
		test.NewRead("app/name/bool_var", true),
		test.NewRead("time_var", test.Time("2020-01-02T15:04:05Z")),
		test.NewReadConfig("cfg"),
	}

	test.Run(t, prov, read)
}
