package json_test

import (
	"testing"
	"time"

	provider "gitoa.ru/go-4devs/config/provider/json"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	js := test.ReadFile("config.json")

	prov := provider.New(js)
	sl := []string{}
	read := []test.Read{
		test.NewRead("app.name.title", "config title"),
		test.NewRead("app.name.timeout", time.Minute),
		test.NewReadUnmarshal("app.name.var", &[]string{"name"}, &sl),
		test.NewReadConfig("cfg"),
		test.NewRead("app.name.success", true),
	}

	test.Run(t, prov, read)
}
