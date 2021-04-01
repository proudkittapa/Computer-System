package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
	//mutex sync.Mutex
	//totalTime
)

func getQuantity(tx *sql.Tx, t chan int, id int) {
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	var name string
	var quantity int
	var price float32
	err := rows.Scan(&name, &quantity, &price)
	if err != nil {
		fmt.Println("get quantity fail")
		tx.Rollback()
		return
	}
	t <- quantity
	fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)

}

func decrement(tx *sql.Tx, t chan int, transactionC chan string, orderQuantity int, id int) {

	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		fmt.Println("the order is out of stock")
		transactionC <- "not complete"
		return
	}
	fmt.Println(newQuantity)

	_, err := tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	if err != nil {
		fmt.Println("decrement fail")
		tx.Rollback()
		transactionC <- "rollback"
		return
	}
	transactionC <- "done"
}

func insert(tx *sql.Tx, user string, id int, q int) {
	tx.Exec("set transaction isolation level SERIALIZABLE")
	_, err := tx.ExecContext(ctx, "INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		tx.Rollback()
		return
	}

}

func preorder(end chan int, user string, productId int, orderQuantity int) {

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	start := time.Now()

	transactionC := make(chan string)
	t := make(chan int)
	go getQuantity(tx, t, productId)
	go decrement(tx, t, transactionC, orderQuantity, productId)
	if <-transactionC == "rollback" {
		fmt.Println("rollback")
		preorder(end, user, productId, orderQuantity)
		return
	}
	// go insert(tx, user, productId, orderQuantity)
	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)
	}
	fmt.Println("success")
	fmt.Println("-----------------------------------")
	fmt.Printf("time: %v\n", time.Since(start))
	num, _ := strconv.Atoi(user)
	end <- num
	return

}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	ctx = context.Background()
	n := 50
	end := make(chan int)
	for i := 1; i <= n; i++ {
		go preorder(end, strconv.Itoa(i), 1, 5)
	}
	for i := 1; i <= n; i++ {
		<-end
	}
	return
}
