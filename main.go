package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {
		urlMap := make(map[string]string)
		if err := c.Bind(&urlMap); err != nil {
			return err
		}

		responseCodes := make(map[string]int)
		for name, url := range urlMap {
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				url = "http://" + url
			}

			resp, err := http.Get(url)
			fmt.Println(name)
			fmt.Println(url)
			if err != nil {
				responseCodes[name] = 0
			} else {
				responseCodes[name] = resp.StatusCode
				resp.Body.Close()
			}
		}

		return c.JSON(http.StatusOK, responseCodes)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
