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
	ctx context.Context
	db  *sql.DB
	//mutex sync.Mutex
)

func preorder(user string, productId int, orderQuantity int) {
	start := time.Now()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	tx2, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx2.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	//getQuantity
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(productId))
	if err != nil {
		tx.Rollback()
		return
	}
	var name string
	var quantity int
	var price float32
	err = rows.Scan(&name, &quantity, &price)
	fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)

	//decrement
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		return
	}
	fmt.Println(newQuantity)
	_, err = tx2.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(productId))
	if err != nil {
		panic(err)
		tx2.Rollback()
	}
	fmt.Println("updated")

	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)
	}

	if err := tx2.Commit(); err != nil {
		fmt.Printf("Failed to commit tx2: %v\n", err)
	}

	//insert
	tx3, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx3.Exec("set transaction isolation level SERIALIZABLE")
	_, err = tx3.ExecContext(ctx, "INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, productId, orderQuantity)

	if err := tx3.Commit(); err != nil {
		fmt.Printf("Failed to commit tx3: %v\n", err)
	}

	fmt.Println("---------------------------")
	fmt.Printf("time: %v\n", time.Since(start))
	return

}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")

	end := make(chan int)
	for i := 1; i < 100; i++ {
		go preorder(strconv.Itoa(i), 1, 5)
	}
	for i := range end {
		fmt.Println(i)
	}
}
