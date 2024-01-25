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
	"gitoa.ru/go-4devs/config/provider/watcher"
	"gitoa.ru/go-4devs/config/provider/yaml"
	"gitoa.ru/go-4devs/config/test"
)

func ExampleClient_Value() {
	const (
		namespace = "fdevs"
		appName   = "config"
	)

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

	// read json config
	jsonConfig := test.ReadFile("config.json")

	config, err := config.New(
		arg.New(),
		env.New(test.Namespace, test.AppName),
		etcd.NewProvider(namespace, appName, etcdClient),
		json.New(jsonConfig),
	)
	if err != nil {
		log.Print(err)

		return
	}

	port, err := config.Value(ctx, "listen")
	if err != nil {
		log.Print("listen", err)

		return
	}

	enabled, err := config.Value(ctx, "maintain")
	if err != nil {
		log.Print("maintain", err)

		return
	}

	title, err := config.Value(ctx, "app.name.title")
	if err != nil {
		log.Print("app.name.title", err)

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

	fmt.Printf("listen from env: %d\n", port.Int())
	fmt.Printf("maintain from etcd: %v\n", enabled.Bool())
	fmt.Printf("title from json: %v\n", title.String())
	fmt.Printf("struct from json: %+v\n", cfg)
	fmt.Printf("replace env host by args: %v\n", hostValue.String())
	// Output:
	// listen from env: 8080
	// maintain from etcd: true
	// title from json: config title
	// struct from json: {Duration:21m0s Enabled:true}
	// replace env host by args: gitoa.ru
}

func ExampleClient_Watch() {
	const (
		namespace = "fdevs"
		appName   = "config"
	)

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

	watcher, err := config.New(
		watcher.New(time.Microsecond, env.New(test.Namespace, test.AppName)),
		watcher.New(time.Microsecond, yaml.NewWatch("test/fixture/config.yaml")),
		etcd.NewProvider(namespace, appName, etcdClient),
	)
	if err != nil {
		log.Print(err)

		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	err = watcher.Watch(ctx, func(ctx context.Context, oldVar, newVar config.Value) {
		fmt.Println("update example_enable old: ", oldVar.Bool(), " new:", newVar.Bool())
		wg.Done()
	}, "example_enable")
	if err != nil {
		log.Print(err)

		return
	}

	_ = os.Setenv("FDEVS_CONFIG_EXAMPLE_ENABLE", "false")

	err = watcher.Watch(ctx, func(ctx context.Context, oldVar, newVar config.Value) {
		fmt.Println("update example_db_dsn old: ", oldVar.String(), " new:", newVar.String())
		wg.Done()
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
	// update example_enable old:  true  new: false
	// update example_db_dsn old:  pgsql://user@pass:127.0.0.1:5432  new: mysql://localhost:5432
}

func ExampleClient_Value_factory() {
	ctx := context.Background()
	_ = os.Setenv("FDEVS_CONFIG_LISTEN", "8080")
	_ = os.Setenv("FDEVS_CONFIG_HOST", "localhost")

	args := os.Args

	defer func() {
		os.Args = args
	}()

	os.Args = []string{"main.go", "--config-json=config.json", "--config-yaml=test/fixture/config.yaml"}

	config, err := config.New(
		arg.New(),
		env.New(test.Namespace, test.AppName),
		config.Factory(func(ctx context.Context, cfg config.Provider) (config.Provider, error) {
			val, err := cfg.Value(ctx, "config-json")
			if err != nil {
				return nil, fmt.Errorf("failed read config file:%w", err)
			}
			jsonConfig := test.ReadFile(val.String())

			return json.New(jsonConfig), nil
		}),
		config.Factory(func(ctx context.Context, cfg config.Provider) (config.Provider, error) {
			val, err := cfg.Value(ctx, "config-yaml")
			if err != nil {
				return nil, fmt.Errorf("failed read config file:%w", err)
			}

			provader, err := yaml.NewFile(val.String())
			if err != nil {
				return nil, fmt.Errorf("failed init by file %v:%w", val.String(), err)
			}

			return provader, nil
		}),
	)
	if err != nil {
		log.Print(err)

		return
	}

	port, err := config.Value(ctx, "listen")
	if err != nil {
		log.Print(err)

		return
	}

	title, err := config.Value(ctx, "app", "name", "title")
	if err != nil {
		log.Print(err)

		return
	}

	yamlTitle, err := config.Value(ctx, "app", "title")
	if err != nil {
		log.Print(err)

		return
	}

	cfgValue, err := config.Value(ctx, "cfg")
	if err != nil {
		log.Print(err)

		return
	}

	cfg := test.Config{}
	_ = cfgValue.Unmarshal(&cfg)

	fmt.Printf("listen from env: %d\n", port.Int())
	fmt.Printf("title from json: %v\n", title.String())
	fmt.Printf("yaml title: %v\n", yamlTitle.String())
	fmt.Printf("struct from json: %+v\n", cfg)
	// Output:
	// listen from env: 8080
	// title from json: config title
	// yaml title: yaml title
	// struct from json: {Duration:21m0s Enabled:true}
}
