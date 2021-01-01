package app

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

type Application struct {
	Config    *Config
	Logger    *zap.Logger
}

func New(config *Config) (app *Application, err error) {
	app = &Application{
		Config:    config,
	}

	app.Logger, err = NewLogger(app.Config.Level)

	if err != nil {
		return nil, err
	}
	app.Logger.Debug("debug mode on")

	return app, nil
}

func (app *Application) Close() {
	app.Logger.Debug("Application stops")
}
