package json_test

import (
	"context"
	"fmt"
	"log"

	"gitoa.ru/go-4devs/config"
	"gitoa.ru/go-4devs/config/provider/json"
	"gitoa.ru/go-4devs/config/test"
)

func ExampleClient_Value() {
	ctx := context.Background()

	// read json config
	jsonConfig, jerr := fixture.ReadFile("fixture/config.json")
	if jerr != nil {
		log.Printf("failed load file:%v", jerr)

		return
	}

	config, err := config.New(
		json.New(jsonConfig),
	)
	if err != nil {
		log.Print(err)

		return
	}

	title, err := config.Value(ctx, "app.name.title")
	if err != nil {
		log.Print("app.name.title", err)

		return
	}

	cfgValue, err := config.Value(ctx, "cfg")
	if err != nil {
		log.Print("cfg ", err)

		return
	}

	cfg := test.Config{}
	_ = cfgValue.Unmarshal(&cfg)

	fmt.Printf("title from json: %v\n", title.String())
	fmt.Printf("struct from json: %+v\n", cfg)
	// Output:
	// title from json: config title
	// struct from json: {Duration:21m0s Enabled:true}
}
