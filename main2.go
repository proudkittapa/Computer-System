package main

import (
	"pin2pre/Kittapa"
)

func main() {
	s := Kittapa.New()
	// db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")

	s.GET("/", abc)
	s.GET("/products", product)
	s.Start(":8080")
}

func abc() string {
	a := "abc"
	return a
}

func product() string {
	return "all products"
}
