package app

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (app *Application) TranslateRequestHandler(c echo.Context) error {
	app.Logger.Info("on TranslatorRequest Handler *******")
	dictionary := c.FormValue("dictionary")
	output_type := c.FormValue("output_type")
	err := app.saveFile(c)
	if err != nil {
		app.Logger.Debug("Error on getting file", zap.Error(err))
		return err
	}
	switch dictionary {
	case LONGMAN:
		dictionary = LONGMAN
		//app.Config.BaseUrl = "https://www.ldoceonline.com/dictionary/"
		app.Config.BaseUrl = "https://cors.yuseferi.workers.dev/dictionary/"
	case LINGUEE:
		dictionary = LINGUEE
		app.Config.BaseUrl = "https://www.linguee.com/english-german/search?source=auto&query="
	default:
		app.Logger.Panic("Wrong dictionary type")
		return nil
	}
	switch output_type {
	case OUTPUT_TYPE_FULL_HTML:
	case OUTPUT_TYPE_IFRAME:
	default:
		app.Logger.Panic("Wrong output type")
		return nil
	}
	if err := app.ScrapeAll(dictionary,output_type); err != nil {
		app.Logger.Panic("HTTP Server start error", zap.Error(err))
	}

	if err := app.export(c); err != nil {
		app.Logger.Panic("HTTP Server start error", zap.Error(err))
	}
	app.cleanup()

	return nil
}
