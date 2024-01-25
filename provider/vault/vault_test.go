package vault_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
	"gitoa.ru/go-4devs/config/test"
)

const token = "dev"

func NewVault() (*api.Client, error) {
	address, ok := os.LookupEnv("VAULT_DEV_LISTEN_ADDRESS")
	if !ok {
		address = "http://127.0.0.1:8200"
	}

	tokenID, ok := os.LookupEnv("VAULT_DEV_ROOT_TOKEN_ID")
	if !ok {
		tokenID = token
	}

	cl, err := api.NewClient(&api.Config{
		Address: address,
	})
	if err != nil {
		return nil, err
	}

	cl.SetToken(tokenID)

	values := map[string]map[string]interface{}{
		"database": {
			"duration": 1260000000000,
			"enabled":  true,
		},
		"db": {
			"dsn":     test.DSN,
			"timeout": "60s",
		},
		"example": {
			"dsn":     test.DSN,
			"timeout": "60s",
		},
	}

	for name, val := range values {
		if err := create(address, tokenID, name, val); err != nil {
			return nil, err
		}
	}

	return cl, nil
}

func create(host, token, path string, data map[string]interface{}) error {
	type Req struct {
		Data interface{} `json:"data"`
	}

	b, err := json.Marshal(Req{Data: data})
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(b)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		host+"/v1/secret/data/fdevs/config/"+path,
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Set("X-Vault-Token", token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return res.Body.Close()
}
