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
	return c.String(http.StatusOK, fmt.Sprintf("name: %s\n", n))
}

func main() {
	e := echo.New()
	e.GET("/", hello)
	e.GET("/user", getNames)
	e.Start("8080")
}
