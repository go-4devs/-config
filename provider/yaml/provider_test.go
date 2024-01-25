package yaml_test

import (
	"embed"
	"testing"
	"time"

	"gitoa.ru/go-4devs/config/provider/yaml"
	"gitoa.ru/go-4devs/config/test"
	"gitoa.ru/go-4devs/config/test/require"
)

//go:embed fixture/*
var fixture embed.FS

func TestProvider(t *testing.T) {
	t.Parallel()

	data, err := fixture.ReadFile("fixture/config.yaml")
	require.NoError(t, err)
	prov, err := yaml.New(data)
	require.NoError(t, err)

	read := []test.Read{
		test.NewRead(21*time.Minute, "duration_var"),
		test.NewRead(true, "app", "name", "bool_var"),
		test.NewRead(test.Time("2020-01-02T15:04:05Z"), "time_var"),
		test.NewReadConfig("cfg"),
	}

	test.Run(t, prov, read)
}
