package main

import (
	"github.com/wolf3w/tag_test/internal/app"
	"github.com/wolf3w/tag_test/internal/domain"

	"go.uber.org/zap"
)

func main() {
	config, err := domain.NewFromEnv()
	if err != nil {
		panic(err)
	}

	logger := zap.Must(zap.NewDevelopment())

	fsApp := app.NewApp(logger, config)

	err = fsApp.Run()
	if err != nil {
		panic(err)
	}
}
