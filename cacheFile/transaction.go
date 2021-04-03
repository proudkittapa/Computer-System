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
	var name string
	var quantity int
	var price int
	err := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = "+strconv.Itoa(1)).Scan(&name, &quantity, &price)
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

func Insert(tx *sql.Tx, wg *sync.WaitGroup, user string, id int, q int) {
	_, err := tx.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		fmt.Println("insert fail")
		tx.Rollback()
		return
	}
	wg.Done()
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
	var wg sync.WaitGroup
	wg.Add(1)
	go Insert(tx, &wg, user, productId, orderQuantity)
	wg.Wait()
	if err = tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)
	}
	Success = true
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
	// db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	ctx = context.Background()
	//n := 100
	end := make(chan bool)
	go Preorder(end, "1", id, quantity)

	Success = <-end
	// fmt.Println("before return")
	return Success
}
