package main

import (
	"context"
	"log"

	"profkom/config"
	"profkom/internal/app"

	cfgloader "profkom/pkg/config"
)

func main() {
	ctx := context.Background()
	cfg := &config.Config{}

	if err := cfgloader.LoadConfig(ctx, cfg); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
