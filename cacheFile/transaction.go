package cacheFile

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	// db  *sql.DB
	//mutex sync.Mutex
	TotalTime float64
	Success   bool
)

func GetQuantity(tx *sql.Tx, t chan int, id int) {
	fmt.Println(id)
	// rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	rows := tx.QueryRow("select unit_price from products where product_id = " + strconv.Itoa(id))

	// var name string
	// var quantity int
	var price int
	// err := rows.Scan(&name, &quantity, &price)
	err := rows.Scan(&price)
	if err != nil {
		fmt.Println("get quantity fail")
		tx.Rollback()
		return
	}
	t <- 1
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

func Insert(tx *sql.Tx, wg *sync.WaitGroup, transactionC chan string, user string, id int, q int) {
	_, err := tx.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		fmt.Println("insert fail")
		tx.Rollback()
		return
	}
	wg.Done()
	transactionC <- "finish"
}

func Preorder(end chan bool, user string, productId int, orderQuantity int) {

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	start := time.Now()
	var name string
	fmt.Println(productId)
	tx.QueryRow("select name from products where product_id = " + strconv.Itoa(1)).Scan(&name)
	fmt.Println("name: ", name)
	transactionC := make(chan string)
	t := make(chan int)
	go GetQuantity(tx, t, productId)
	go Decrement(tx, t, transactionC, orderQuantity, productId)
	if <-transactionC == "rollback" {
		fmt.Println("rollback")
		Preorder(end, user, productId, orderQuantity)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go Insert(tx, &wg, transactionC, user, productId, orderQuantity)
	wg.Wait()
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
	var name string
	db.QueryRow("select name from products where product_id = " + strconv.Itoa(1)).Scan(&name)
	fmt.Println(name)
	ctx = context.Background()
	//n := 100
	end := make(chan bool)
	go Preorder(end, "1", id, quantity)

	Success = <-end
	fmt.Println("before return")
	return Success
}
