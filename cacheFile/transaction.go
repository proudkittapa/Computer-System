package cacheFile

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
	TotalTime float64
	Success   bool
)

func GetQuantity(tx *sql.Tx, t chan int, id int) {
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
	//fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)

}

func Decrement(tx *sql.Tx, t chan int, transactionC chan string, orderQuantity int, id int) {

	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		fmt.Println("the order is out of stock")
		transactionC <- "not complete"
		return
	}
	//fmt.Println(newQuantity)

	_, err := tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	if err != nil {
		fmt.Println("decrement fail")
		tx.Rollback()
		transactionC <- "rollback"
		return
	}
	transactionC <- "done"
}

func Insert(tx *sql.Tx, transactionC chan string, user string, id int, q int) {
	tx.Exec("set transaction isolation level SERIALIZABLE")
	_, err := tx.ExecContext(ctx, "INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		fmt.Println("insert fail")
		tx.Rollback()
		return
	}
	transactionC <- "finish"
}

func Preorder(end chan bool, user string, productId int, orderQuantity int) {

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	start := time.Now()

	transactionC := make(chan string)
	t := make(chan int)
	go GetQuantity(tx, t, productId)
	go Decrement(tx, t, transactionC, orderQuantity, productId)
	if <-transactionC == "rollback" {
		fmt.Println("rollback")
		Preorder(end, user, productId, orderQuantity)
		return
	}
	go Insert(tx, transactionC, user, productId, orderQuantity)
	if err = tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)
	}
	if <-transactionC == "finish" {
		Success = true
	}
	//fmt.Println("success")
	//fmt.Println("-----------------------------------")
	elapsed := time.Since(start)
	tt := float64(elapsed)
	fmt.Printf("time: %v\n", elapsed)
	fmt.Printf("tt: %v\n", tt)
	TotalTime += tt
	fmt.Printf("total time: %v\n", TotalTime)
	end <- Success
}
func PostPreorder(id int, quantity int) bool {
	db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	ctx = context.Background()
	//n := 100
	end := make(chan bool)
	go Preorder(end, "1", id, quantity)

	Success = <-end
	fmt.Println("before return")
	return Success
}
