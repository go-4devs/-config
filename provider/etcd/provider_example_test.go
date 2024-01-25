package etcd_test

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/etcd"
)

func ExampleClient_Value() {
	const (
		namespace = "fdevs"
		appName   = "config"
	)

	ctx := context.Background()

	// configure etcd client
	etcdClient, err := NewEtcd(ctx)
	if err != nil {
		log.Print(err)

		return
	}

	config, err := config.New(
		etcd.New(namespace, appName, etcdClient),
	)
	if err != nil {
		log.Print(err)

		return
	}

	enabled, err := config.Value(ctx, "maintain")
	if err != nil {
		log.Print("maintain ", err)

		return
	}

	fmt.Printf("maintain from etcd: %v\n", enabled.Bool())
	// Output:
	// maintain from etcd: true
}

func ExampleClient_Watch() {
	const (
		namespace = "fdevs"
		appName   = "config"
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// configure etcd client
	etcdClient, err := NewEtcd(ctx)
	if err != nil {
		log.Print(err)

		return
	}

	_, err = etcdClient.KV.Put(ctx, "fdevs/config/example_db_dsn", "pgsql://user@pass:127.0.0.1:5432")
	if err != nil {
		log.Print(err)

		return
	}

	defer func() {
		cancel()

		if _, err = etcdClient.KV.Delete(context.Background(), "fdevs/config/example_db_dsn"); err != nil {
			log.Print(err)

			return
		}
	}()

	watcher, err := config.New(
		etcd.New(namespace, appName, etcdClient),
	)
	if err != nil {
		log.Print(err)

		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	err = watcher.Watch(ctx, func(_ context.Context, oldVar, newVar config.Value) error {
		fmt.Println("update example_db_dsn old: ", oldVar.String(), " new:", newVar.String())
		wg.Done()

		return nil
	}, "example_db_dsn")
	if err != nil {
		log.Print(err)

		return
	}

	time.AfterFunc(time.Second, func() {
		if _, err := etcdClient.KV.Put(ctx, "fdevs/config/example_db_dsn", "mysql://localhost:5432"); err != nil {
			log.Print(err)

			return
		}
	})

	wg.Wait()

	// Output:
	// update example_db_dsn old:  pgsql://user@pass:127.0.0.1:5432  new: mysql://localhost:5432
}
