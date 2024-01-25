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
		test.NewRead("config title", "app.name.title"),
		test.NewRead(time.Minute, "app.name.timeout"),
		test.NewReadUnmarshal(&[]string{"name"}, &sl, "app.name.var"),
		test.NewReadConfig("cfg"),
		test.NewRead(true, "app", "name", "success"),
	}

	test.Run(t, prov, read)
}
