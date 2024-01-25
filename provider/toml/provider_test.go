package toml_test

import (
	"embed"
	"testing"

	"gitoa.ru/go-4devs/config/provider/toml"
	"gitoa.ru/go-4devs/config/test"
	"gitoa.ru/go-4devs/config/test/require"
)

//go:embed fixture/*
var fixtures embed.FS

func TestProvider(t *testing.T) {
	t.Parallel()

	files, ferr := fixtures.ReadFile("fixture/config.toml")
	require.NoError(t, ferr)

	prov, err := toml.New(files)
	require.NoError(t, err)

	m := []int{}

	read := []test.Read{
		test.NewRead("192.168.1.1", "database.server"),
		test.NewRead("TOML Example", "title"),
		test.NewRead("10.0.0.1", "servers.alpha.ip"),
		test.NewRead(true, "database.enabled"),
		test.NewRead(5000, "database.connection_max"),
		test.NewReadUnmarshal(&[]int{8001, 8001, 8002}, &m, "database", "ports"),
	}

	test.Run(t, prov, read)
}
