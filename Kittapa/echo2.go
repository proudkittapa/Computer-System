package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Name string `json: "name"`
}

type Dog struct {
	Name string `json: "name"`
	Type string `json: "type"`
}

type Product struct {
	Name string `json: "name"`
	Type string `json: "type"`
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}

func getUsers(c echo.Context) error {
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

//fastest one
func addUser(c echo.Context) error {
	n := User{}
	defer c.Request().Body.Close()
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("fail to read the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	err = json.Unmarshal(b, &n)
	if err != nil {
		log.Printf("Fail unmarshaling addUser: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}
	log.Printf("this is your name: %#v", n)
	return c.String(http.StatusOK, "new user added")

}

func addDog(c echo.Context) error {
	dog := Dog{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Fail processing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("this is your dog: %#v", dog)
	return c.String(http.StatusOK, "new dog added")
}

//slower than the two above
func addProducts(c echo.Context) error {
	p := Product{}
	err := c.Bind(&p)
	if err != nil {
		log.Printf("Fail processing addProduct request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	log.Printf("this is your dog: %#v", p)
	return c.String(http.StatusOK, "new product added")
}

func main() {
	fmt.Println("Server started")
	e := echo.New()
	e.GET("/", hello)
	e.GET("/users/:data", getUsers)
	e.POST("/users", addUser)
	e.POST("/dogs", addDog)
	e.POST("/products", addProducts)
	e.Start(":8080")
}
