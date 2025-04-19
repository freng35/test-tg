package main

import (
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"tg_service/internal/app"
	"tg_service/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no path for config")
		os.Exit(1)
	}

	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("can not read file:", err)
		os.Exit(1)
	}

	var config config.Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("can not parse file:", err)
		os.Exit(1)
	}

	err = run(config)
	_, _ = fmt.Fprintf(os.Stdout, "application is finished: %s", err)
}

func run(cfg config.Config) error {
	ctx := context.Background()
	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		return err
	}

	return app.Run(ctx)
}
