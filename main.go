package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NEWJYH/learngo/scrapper"
	"github.com/labstack/echo/v4"
)

const FILE_NAME string = "jobs.csv"

func handleHome(c echo.Context) error {
	fmt.Println(time.Now())
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(FILE_NAME)
	term := strings.ToLower(scrapper.CleanString(c.FormValue("term")))
	scrapper.Scrape(term)
	return c.Attachment(FILE_NAME, term + time.Now().String() + FILE_NAME)
}

func main() {
	// scrapper.Scrape("python")
	
	e := echo.New()

	e.GET("/", handleHome)

	e.POST("/scrape", handleScrape)

	e.Logger.Fatal(e.Start(":1323"))
}