package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}

func getNames(c echo.Context) error {
	n := c.QueryParam("name")
	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("name: %s\n", n))
	}
	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": n,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "json or string",
	})
}

func main() {
	fmt.Println("Server started")
	e := echo.New()
	e.GET("/", hello)
	e.GET("/users/:data", getNames)
	e.Start(":8080")
}
