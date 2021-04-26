package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// func add_db() {

// }

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:62011139@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

	for i := 0; i < 1000; i++ {
		insert, err := db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shirt', '10000', 20) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Swater', '10000', 30) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Hoodie', '10000', 25) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shoes', '10000', 10) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Pants', '10000', 20) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shirt', '10000', 20) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Swater', '10000', 30) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Hoodie', '10000', 25) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shoes', '10000', 10) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Pants', '10000', 20) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shirt', '10000', 20) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Swater', '10000', 30) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Hoodie', '10000', 25) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Shoes', '10000', 10) ")
		insert, err = db.Query("INSERT INTO products (name, quantity_in_stock, unit_price) VALUE ('Pants', '10000', 20) ")

		insert.Close()
		checkErr(err)

		fmt.Println("done")
	}
}
