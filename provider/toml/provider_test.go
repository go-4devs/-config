package toml_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitoa.ru/go-4devs/config/provider/toml"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	prov, err := toml.NewFile(test.FixturePath("config.toml"))
	require.NoError(t, err)

	m := []int{}

	read := []test.Read{
		test.NewRead("database.server", "192.168.1.1"),
		test.NewRead("title", "TOML Example"),
		test.NewRead("servers.alpha.ip", "10.0.0.1"),
		test.NewRead("database.enabled", true),
		test.NewRead("database.connection_max", 5000),
		test.NewReadUnmarshal("database.ports", &[]int{8001, 8001, 8002}, &m),
	}

	test.Run(t, prov, read)
}
