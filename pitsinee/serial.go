package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	//"sync"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/jackc/pgx/v4"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
	//mutex sync.Mutex
)

func getQuantity(t chan int, tx *sql.Tx, id int) {
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	var name string
	var quantity int
	var price float32
	err := rows.Scan(&name, &quantity, &price)
	if err != nil {
		fmt.Println("get quantity rollback")
		tx.Rollback()
		// return
	}
	t <- quantity
	fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
}

func decrement(t chan int, tx *sql.Tx, transactionC chan string, orderQuantity int, id int) {
	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		transactionC <- "done"
		return
	}
	fmt.Println(newQuantity)
	_, err := tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	if err != nil {
		fmt.Println("decrement rollback")
		tx.Rollback()
		transactionC <- "rollback"
		return
	}
	transactionC <- "done"
}

func insert(user string, id int, q int) {
	tx3, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx3.Exec("set transaction isolation level SERIALIZABLE")
	_, err = tx3.ExecContext(ctx, "INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err := tx3.Commit(); err != nil {
		fmt.Printf("Failed to commit tx3: %v\n", err)
	}
}

func preorder(end chan int, user string, productId int, orderQuantity int) {
	// fmt.Printf("start\n")
	//start := time.Now()
	transactionC := make(chan string)
	t := make(chan int)
	ctx = context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	//xType := fmt.Sprintf("%T", tx)
	//fmt.Println(xType) // "[]int"
	go getQuantity(t, tx, productId)
	go decrement(t, tx, transactionC, orderQuantity, productId)
	if <-transactionC == "rollback" {
		preorder(end, user, productId, orderQuantity)
		return
	}
	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx2: %v\n", err)
	}
	// go insert(user, productId, orderQuantity)
	//fmt.Printf("time: %v\n", time.Since(start))
	num, _ := strconv.Atoi(user)
	end <- num
	return
}

func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	n := 10
	end := make(chan int, n-1)
	for i := 1; i <= n; i++ {
		go preorder(end, strconv.Itoa(i), 1, 1)
	}
	for i := 1; i < n; i++ {
		<-end
	}
}
