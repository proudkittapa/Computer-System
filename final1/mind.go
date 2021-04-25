package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB
	q           int
	newQuantity int
)

func getQuantity(id int) {
	row, err := db.Query("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	for row.Next() {
		var name string
		var quantity int
		var price int
		row.Scan(&name, &quantity, &price)
		q = quantity
		fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
	}
}
func decrement(orderQuantity int, id int) {
	newQuantity := q - orderQuantity
	if newQuantity < 0 {
		return
	}
	fmt.Println("new quantity: ", newQuantity)
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)

	return
}

func insert(user string, id int, q int) {
	db.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func preorder(user string, productId int, orderQuantity int) {
	start := time.Now()
	insert(user, productId, orderQuantity)
	getQuantity(productId)
	decrement(orderQuantity, productId)
	fmt.Printf("time: %v\n", time.Since(start))
	return
}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	start := time.Now()
	for i := 1; i < 10; i++ {
		go preorder(strconv.Itoa(i), 1, 1)
		fmt.Printf("Total time: %v\n", time.Since(start))
	}
	fmt.Printf("Total time: %v\n", time.Since(start))

}
