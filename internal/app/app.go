package app

import (
	"go.uber.org/zap"

	"github.com/wolf3w/tag_test/internal/api"
	"github.com/wolf3w/tag_test/internal/domain"
)

type App struct {
	logger *zap.Logger
	conf   *domain.Config
}

func NewApp(logger *zap.Logger, conf *domain.Config) *App {
	return &App{
		logger: logger,
		conf:   conf,
	}
}

func (app *App) Run() error {
	return api.RunServer(app.logger, app.conf)
}
