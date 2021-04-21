package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "You are on the admin page")
}

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Example/1.0")
		c.Response().Header().Set("notReallyHeader", "notHeader")
		return next(c)
	}
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "You are on the cookie page")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")
	//check username and password against DP after hashing the password
	if username == "proud" && password == "kittapa" {
		cookie := &http.Cookie{}
		//cookie := new(http.Cookie)
		cookie.Name = "sessionID"    //store in db
		cookie.Value = "some_string" //usually store id or email'
		cookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(cookie)
		return c.String(http.StatusOK, "You were logged in!")
	}
	return c.String(http.StatusUnauthorized, "Your username and password are wrong")
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not preseint") {
				return c.String(http.StatusUnauthorized, "you don't have any cookie")
			}
			log.Println(err)
			return err
		}
		if cookie.Value == "some_string" {
			return next(c)
		}
		return c.String(http.StatusUnauthorized, "You don't have the right cookie, cookie")
	}
}
func main() {
	fmt.Println("Server started")
	e := echo.New()

	e.Use(ServerHeader)

	//grouping
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")

	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))
	//three ways to add middleware
	//1 -  add to the Group -> e.Group("/admin", middleware.Logger())
	//2 - g.Use(middleware.Logger())
	//3 - add to the method -> g.GET("/main", mainAdmin, middleware.Logger())

	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "proud" && password == "password" {
			return true, nil
		}
		return false, nil
	}))
	cookieGroup.Use(checkCookie)
	cookieGroup.GET("/main", mainCookie)
	adminGroup.GET("/main", mainAdmin)
	e.GET("/login", login)
	e.GET("/", hello)
	e.GET("/users/:data", getUsers)

	e.POST("/users", addUser)
	e.POST("/dogs", addDog)
	e.POST("/products", addProducts)
	e.Start(":8080")
}
