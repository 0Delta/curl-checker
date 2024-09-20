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

			code, err := getFinalStatusCode(url)
			responseCodes[name] = code
			if err != nil {
				fmt.Println(err)
				responseCodes[name] = -1
			}
		}

		return c.JSON(http.StatusOK, responseCodes)
	})

	e.Logger.Fatal(e.Start(":" + port))
}

func getFinalStatusCode(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location == "" {
			return -1, fmt.Errorf("3xx status code received but no Location header found")
		}
		return getFinalStatusCode(location)
	}

	return resp.StatusCode, nil
}
