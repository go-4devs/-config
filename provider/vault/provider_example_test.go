package vault_test

import (
	"context"
	"fmt"
	"log"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/vault"
)

func ExampleClient_Value() {
	const (
		namespace = "fdevs"
		appName   = "config"
	)

	ctx := context.Background()

	// configure vault client
	vaultClient, err := NewVault()
	if err != nil {
		log.Print(err)

		return
	}

	config, err := config.New(
		vault.New(namespace, appName, vaultClient),
	)
	if err != nil {
		log.Print(err)

		return
	}

	dsn, err := config.Value(ctx, "example", "dsn")
	if err != nil {
		log.Print("example:dsn ", err)

		return
	}

	fmt.Printf("dsn from vault: %s\n", dsn.String())
	// Output:
	// dsn from vault: pgsql://user@pass:127.0.0.1:5432
}
