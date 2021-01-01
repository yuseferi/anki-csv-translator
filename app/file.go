package app

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"os"
)

func (app *Application) saveFile(c echo.Context) error {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		app.Logger.Debug("Error on getting file", zap.Error(err))
		return err
	}
	src, err := file.Open()
	if err != nil {
		app.Logger.Debug("Error on open file to save", zap.Error(err))
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(app.Config.CSVWordInputFile)
	if err != nil {
		app.Logger.Debug("Error on creating file", zap.Error(err))
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		app.Logger.Debug("Error on copy form file content to disk file", zap.Error(err))
		return err
	}
	return nil
}

func (app *Application) export(c echo.Context) error {
	//return c.Attachment(app.Config.CSVWordOutputFile, "translated-words.csv")
	return c.File(app.Config.CSVWordOutputFile)
}

func (app *Application) cleanup() {
	_ = os.Remove(app.Config.CSVWordInputFile)
	_ = os.Remove(app.Config.CSVWordOutputFile)
}
