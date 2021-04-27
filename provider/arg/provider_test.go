package arg_test

import (
	"os"
	"testing"
	"time"

	"gitoa.ru/go-4devs/config/provider/arg"
	"gitoa.ru/go-4devs/config/test"
)

func TestProvider(t *testing.T) {
	t.Parallel()

	args := os.Args

	defer func() {
		os.Args = args
	}()

	os.Args = []string{
		"main.go",
		"--listen=8080",
		"--config=config.hcl",
		"--url=http://4devs.io",
		"--url=https://4devs.io",
		"--timeout=1m",
		"--timeout=1h",
		"--start-at=2010-01-02T15:04:05Z",
		"--end-after=2009-01-02T15:04:05Z",
		"--end-after=2008-01-02T15:04:05+03:00",
	}
	read := []test.Read{
		test.NewRead("listen", 8080),
		test.NewRead("config", "config.hcl"),
		test.NewRead("start-at", test.Time("2010-01-02T15:04:05Z")),
		test.NewReadUnmarshal("url", &[]string{"http://4devs.io", "https://4devs.io"}, &[]string{}),
		test.NewReadUnmarshal("timeout", &[]time.Duration{time.Minute, time.Hour}, &[]time.Duration{}),
		test.NewReadUnmarshal("end-after", &[]time.Time{
			test.Time("2009-01-02T15:04:05Z"),
			test.Time("2008-01-02T15:04:05+03:00"),
		}, &[]time.Time{}),
	}

	prov := arg.New()

	test.Run(t, prov, read)
}
