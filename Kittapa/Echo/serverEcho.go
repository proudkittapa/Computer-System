package main

import (
	"database/sql"
	"fmt"
	"github.com/JonathanMH/goClacks/echo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	// "strconv"
)

type (
	Excuse struct {
		Error    string `json:"error"`
		Id       string `json:"id"`
		Quote    string `json:"quote"`
		Quantity int    `json:"quantity`
	}
)

func main() {
	// Echo instance
	e := echo.New()
	e.Use(goClacks.Terrify)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		db, err := sql.Open("mysql", "root:62011139@tcp(localhost:3306)/test")

		if err != nil {
			fmt.Println(err.Error())
			response := Excuse{Id: "", Error: "true", Quote: ""}
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer db.Close()

		var quote string
		var id string
		var quantity int
		err = db.QueryRow("SELECT id, quote, quantity FROM excuses ORDER BY RAND() LIMIT 1").Scan(&id, &quote, &quantity)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(quote)
		response := Excuse{Id: id, Error: "false", Quote: quote, Quantity: quantity}
		return c.JSON(http.StatusOK, response)
	})

	e.GET("/id/:id", func(c echo.Context) error {
		requested_id := c.Param("id")
		fmt.Println(requested_id)
		db, err := sql.Open("mysql", "root:62011139@tcp(localhost:3306)/test")

		if err != nil {
			fmt.Println(err.Error())
			response := Excuse{Id: "", Error: "true", Quote: ""}
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer db.Close()

		var quote string
		var id string
		var quantity int
		err = db.QueryRow("SELECT id, quote, quantity FROM excuses WHERE id = ?", requested_id).Scan(&id, &quote, &quantity)

		if err != nil {
			fmt.Println(err)
		}

		response := Excuse{Id: id, Error: "false", Quote: quote, Quantity: quantity}
		return c.JSON(http.StatusOK, response)
	})

	e.POST("/id/:id", func(c echo.Context) error {
		requested_id := c.Param("id")
		db, err := sql.Open("mysql", "root:62011139@tcp(localhost:3306)/test")
		// q := 0
		if err != nil {
			fmt.Println(err.Error())
			response := Excuse{Id: "", Error: "true", Quote: ""}
			return c.JSON(http.StatusInternalServerError, response)
		}

		var quote string
		var id string
		var quantity int
		err = db.QueryRow("SELECT id, quote, quantity FROM excuses WHERE id = ?", requested_id).Scan(&id, &quote, &quantity)
		if err != nil {
			fmt.Println(err)
		}
		db.Exec("update excuses set quantity = ? where id = ? ", quantity+1, requested_id)
		response := Excuse{Id: id, Error: "false", Quote: quote, Quantity: quantity + 1}
		return c.JSON(http.StatusOK, response)

	})

	e.Logger.Fatal(e.Start(":4000"))
}
