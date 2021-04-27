package config_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/arg"
	"gitoa.ru/go-4devs/config/provider/env"
	"gitoa.ru/go-4devs/config/provider/etcd"
	"gitoa.ru/go-4devs/config/provider/json"
	"gitoa.ru/go-4devs/config/provider/vault"
	"gitoa.ru/go-4devs/config/provider/watcher"
	"gitoa.ru/go-4devs/config/provider/yaml"
	"gitoa.ru/go-4devs/config/test"
)

func ExampleNew() {
	ctx := context.Background()
	_ = os.Setenv("FDEVS_CONFIG_LISTEN", "8080")
	_ = os.Setenv("FDEVS_CONFIG_HOST", "localhost")

	args := os.Args

	defer func() {
		os.Args = args
	}()

	os.Args = []string{"main.go", "--host=gitoa.ru"}

	// configure etcd client
	etcdClient, err := test.NewEtcd(ctx)
	if err != nil {
		log.Print(err)

		return
	}

	// configure vault client
	vaultClient, err := test.NewVault()
	if err != nil {
		log.Print(err)

		return
	}

	// read json config
	jsonConfig := test.ReadFile("config.json")

	providers := []config.Provider{
		arg.New(),
		env.New(),
		etcd.NewProvider(etcdClient),
		vault.NewSecretKV2(vaultClient),
		json.New(jsonConfig),
	}
	config := config.New(test.Namespace, test.AppName, providers)

	dsn, err := config.Value(ctx, "example:dsn")
	if err != nil {
		log.Print(err)

		return
	}

	port, err := config.Value(ctx, "listen")
	if err != nil {
		log.Print(err)

		return
	}

	enabled, err := config.Value(ctx, "maintain")
	if err != nil {
		log.Print(err)

		return
	}

	title, err := config.Value(ctx, "app.name.title")
	if err != nil {
		log.Print(err)

		return
	}

	cfgValue, err := config.Value(ctx, "cfg")
	if err != nil {
		log.Print(err)

		return
	}

	hostValue, err := config.Value(ctx, "host")
	if err != nil {
		log.Print(err)

		return
	}

	cfg := test.Config{}
	_ = cfgValue.Unmarshal(&cfg)

	fmt.Printf("dsn from vault: %s\n", dsn.String())
	fmt.Printf("listen from env: %d\n", port.Int())
	fmt.Printf("maintain from etcd: %v\n", enabled.Bool())
	fmt.Printf("title from json: %v\n", title.String())
	fmt.Printf("struct from json: %+v\n", cfg)
	fmt.Printf("replace env host by args: %v\n", hostValue.String())
	// Output:
	// dsn from vault: pgsql://user@pass:127.0.0.1:5432
	// listen from env: 8080
	// maintain from etcd: true
	// title from json: config title
	// struct from json: {Duration:21m0s Enabled:true}
	// replace env host by args: gitoa.ru
}

func ExampleNewWatch() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// configure etcd client
	etcdClient, err := test.NewEtcd(ctx)
	if err != nil {
		log.Print(err)

		return
	}

	_ = os.Setenv("FDEVS_CONFIG_EXAMPLE_ENABLE", "true")

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

	providers := []config.Provider{
		watcher.New(time.Microsecond, env.New()),
		watcher.New(time.Microsecond, yaml.NewFile("test/fixture/config.yaml")),
		etcd.NewProvider(etcdClient),
	}
	watcher := config.NewWatch(test.Namespace, test.AppName, providers)
	wg := sync.WaitGroup{}
	wg.Add(2)

	err = watcher.Watch(ctx, "example_enable", func(ctx context.Context, oldVar, newVar config.Variable) {
		fmt.Println("update ", oldVar.Provider, " variable:", oldVar.Name, ", old: ", oldVar.Value.Bool(), " new:", newVar.Value.Bool())
		wg.Done()
	})
	if err != nil {
		log.Print(err)

		return
	}

	_ = os.Setenv("FDEVS_CONFIG_EXAMPLE_ENABLE", "false")

	err = watcher.Watch(ctx, "example_db_dsn", func(ctx context.Context, oldVar, newVar config.Variable) {
		fmt.Println("update ", oldVar.Provider, " variable:", oldVar.Name, ", old: ", oldVar.Value.String(), " new:", newVar.Value.String())
		wg.Done()
	})
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
	// update  env  variable: FDEVS_CONFIG_EXAMPLE_ENABLE , old:  true  new: false
	// update  etcd  variable: fdevs/config/example_db_dsn , old:  pgsql://user@pass:127.0.0.1:5432  new: mysql://localhost:5432
}
