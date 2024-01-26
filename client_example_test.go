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
	"gitoa.ru/go-4devs/config/provider/watcher"
	"gitoa.ru/go-4devs/config/test"
)

func ExampleClient_Value() {
	ctx := context.Background()
	_ = os.Setenv("FDEVS_CONFIG_LISTEN", "8080")
	_ = os.Setenv("FDEVS_CONFIG_HOST", "localhost")

	args := os.Args

	defer func() {
		os.Args = args
	}()

	os.Args = []string{"main.go", "--host=gitoa.ru"}

	// read json config

	config, err := config.New(
		arg.New(),
		env.New(test.Namespace, test.AppName),
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

	hostValue, err := config.Value(ctx, "host")
	if err != nil {
		log.Print("host ", err)

		return
	}

	fmt.Printf("listen from env: %d\n", port.Int())
	fmt.Printf("replace env host by args: %v\n", hostValue.String())
	// Output:
	// listen from env: 8080
	// replace env host by args: gitoa.ru
}

func ExampleClient_Watch() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = os.Setenv("FDEVS_CONFIG_EXAMPLE_ENABLE", "true")

	watcher, err := config.New(
		watcher.New(time.Microsecond, env.New(test.Namespace, test.AppName)),
	)
	if err != nil {
		log.Print(err)

		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	err = watcher.Watch(ctx, func(ctx context.Context, oldVar, newVar config.Value) error {
		fmt.Println("update example_enable old: ", oldVar.Bool(), " new:", newVar.Bool())
		wg.Done()

		return nil
	}, "example_enable")
	if err != nil {
		log.Print(err)

		return
	}

	_ = os.Setenv("FDEVS_CONFIG_EXAMPLE_ENABLE", "false")

	err = watcher.Watch(ctx, func(ctx context.Context, oldVar, newVar config.Value) error {
		fmt.Println("update example_db_dsn old: ", oldVar.String(), " new:", newVar.String())
		wg.Done()

		return nil
	}, "example_db_dsn")
	if err != nil {
		log.Print(err)

		return
	}

	wg.Wait()

	// Output:
	// update example_enable old:  true  new: false
}

func ExampleClient_Value_factory() {
	ctx := context.Background()
	_ = os.Setenv("FDEVS_CONFIG_LISTEN", "8080")
	_ = os.Setenv("FDEVS_CONFIG_HOST", "localhost")
	_ = os.Setenv("FDEVS_GOLANG_HOST", "go.dev")

	args := os.Args

	defer func() {
		os.Args = args
	}()

	os.Args = []string{"main.go", "--env=golang"}

	config, err := config.New(
		arg.New(),
		config.Factory(func(ctx context.Context, cfg config.Provider) (config.Provider, error) {
			val, err := cfg.Value(ctx, "env")
			if err != nil {
				return nil, fmt.Errorf("failed read config file:%w", err)
			}

			return env.New(test.Namespace, val.String()), nil
		}),
		env.New(test.Namespace, test.AppName),
	)
	if err != nil {
		log.Print(err)

		return
	}

	envName, err := config.Value(ctx, "env")
	if err != nil {
		log.Print("env ", err)

		return
	}

	host, err := config.Value(ctx, "host")
	if err != nil {
		log.Print("host ", err)

		return
	}

	listen, err := config.Value(ctx, "listen")
	if err != nil {
		log.Print("listen", err)

		return
	}

	fmt.Printf("envName from env: %s\n", envName.String())
	fmt.Printf("host from env with app name golang: %s\n", host.String())
	fmt.Printf("listen from env with default app name: %s\n", listen.String())
	// Output:
	// envName from env: golang
	// host from env with app name golang: go.dev
	// listen from env with default app name: 8080
}
