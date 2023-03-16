package main

import (
	// "net/http"

	"github.com/NEWJYH/learngo/scrapper"
	// "github.com/labstack/echo/v4"
)


func main() {
	scrapper.Scrape("python")
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.Logger.Fatal(e.Start(":1323"))
}