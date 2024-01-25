package etcd_test

import (
	"context"
	"os"
	"time"

	client "go.etcd.io/etcd/client/v3"
)

const ConfigJSON = `{"duration":1260000000000,"enabled":true}`

func NewEtcd(ctx context.Context) (*client.Client, error) {
	dsn, ok := os.LookupEnv("FDEVS_CONFIG_ETCD_HOST")
	if !ok {
		dsn = "127.0.0.1:2379"
	}

	et, err := client.New(client.Config{
		Endpoints:   []string{dsn},
		DialTimeout: time.Second,
	})
	if err != nil {
		return nil, err
	}

	values := map[string]string{
		"fdevs/config/db_dsn":   "pgsql://user@pass:127.0.0.1:5432",
		"fdevs/config/duration": "12m",
		"fdevs/config/port":     "8080",
		"fdevs/config/maintain": "true",
		"fdevs/config/start_at": "2020-01-02T15:04:05Z",
		"fdevs/config/percent":  "0.064",
		"fdevs/config/count":    "2020",
		"fdevs/config/int64":    "2021",
		"fdevs/config/uint64":   "2022",
		"fdevs/config/config":   ConfigJSON,
	}

	for name, val := range values {
		_, err = et.KV.Put(ctx, name, val)
		if err != nil {
			return nil, err
		}
	}

	return et, nil
}
