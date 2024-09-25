package main

import (
	"context"
	"flag"
	"log"

	"github.com/nqxcode/auth_microservice/internal/app"
)

var logLevel = flag.String("l", "debug", "log level")

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
