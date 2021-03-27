package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"context"
	//"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	//mutex sync.Mutex
)

func getQuantity(t chan int, id int) {

	row, err := db.Query("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	for row.Next() {
		var name string
		var quantity int
		var price int
		row.Scan(&name, &quantity, &price)
		t <- quantity
		fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
	}
}
func decrement(t chan int, transactionC chan int, orderQuantity int, id int) {
	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		transactionC <- 0
		return
	}
	fmt.Println(newQuantity)

	db.Query("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)
	transactionC <- 0
}

func insert(user string, id int, q int) {

	db.Query("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func preorder(end chan int, user string, productId int, orderQuantity int) {
	start := time.Now()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	//insert
	_, err = tx.ExecContext(ctx, "INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, productId, orderQuantity)
	if err != nil {
		panic(err)
	}
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(productId))
	var name string
	var quantity int
	var price float32
	err = rows.Scan(&name, &quantity, &price)
	fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
	if err != nil {
		tx.Rollback()
		return
	}

	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		return
	}
	fmt.Println(newQuantity)
	_, err = tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(productId))
	if err != nil {
		panic(err)
		tx.Rollback()
	}
	fmt.Println("updated")
	fmt.Println("---------------------------")
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("time: %v\n", time.Since(start))
	return

}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	end := make(chan int)
	for i := 1; i < 100; i++ {
		go preorder(end, strconv.Itoa(i), 1, 5)
	}
	for i := range end {
		fmt.Println(i)
	}
}
