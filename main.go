package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yuseferi/anki-csv-translator/app"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/", "public")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":8086"))
}

func upload(c echo.Context) error{
	config, err := app.NewConfig()
	if err != nil {
		panic(err)
	}
	application, err := app.New(config)
	if err != nil {
		panic(err)
	}
	defer application.Close()

	return application.TranslateRequestHandler(c)
}
