package main

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
	db  *sql.DB
	//mutex sync.Mutex
	totalTime float64
)

func getQuantity(tx *sql.Tx, transactionC chan string, t chan int, id int) {
	fmt.Println("stop1")
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	var name string
	var quantity int
	var price float32
	err := rows.Scan(&name, &quantity, &price)
	if err != nil {
		//fmt.Println("get quantity fail")
		fmt.Println("stop8")
		transactionC <- "rollback"
		tx.Rollback()
		return
	}
	fmt.Println("stop2")
	t <- quantity
	//fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
}

func decrement(tx *sql.Tx, t chan int, transactionC chan string, orderQuantity int, id int) {
	fmt.Println("stop3")
	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		//fmt.Println("the order is out of stock")
		fmt.Println("stop4")
		transactionC <- "not complete"
		return
	}
	//fmt.Println(newQuantity)
	_, err := tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	fmt.Println("stop5")
	if err != nil {
		//fmt.Println("decrement fail")
		tx.Rollback()
		transactionC <- "rollback"
		return
	}
	fmt.Println("stop6")
	transactionC <- "done"
}

func insert(wg *sync.WaitGroup, tx *sql.Tx, user string, id int, q int) {
	_, err := tx.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		fmt.Println("insert fail")
		tx.Rollback()
		return
	}
	fmt.Println("stop7")
	wg.Done()
}

func preorder(end chan int, user string, productId int, orderQuantity int) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	num, _ := strconv.Atoi(user)
	if err != nil {
		panic(err)
	}
	start := time.Now()
	transactionC := make(chan string)
	t := make(chan int)
	go getQuantity(tx, transactionC, t, productId)
	go decrement(tx, t, transactionC, orderQuantity, productId)
	if <-transactionC == "rollback" {
		//fmt.Println("rollback")
		preorder(end, user, productId, orderQuantity)
		return
	}
	if <-transactionC == "not complete" {
		end <- num
		return
	}
	fmt.Println("user:", user, "productId:", productId, "orderQuantity:", orderQuantity)
	var wg sync.WaitGroup
	wg.Add(1)
	go insert(&wg, tx, user, productId, orderQuantity)
	wg.Wait()
	if err := tx.Commit(); err != nil {
		//fmt.Printf("Failed to commit tx: %v\n", err)
	}
	//fmt.Println("success")
	//fmt.Println("-----------------------------------")
	elapsed := time.Since(start)
	tt := float64(elapsed)
	fmt.Printf("time: %v\n", elapsed)
	fmt.Printf("tt: %v\n", tt)
	totalTime += tt
	fmt.Printf("total time: %v\n", totalTime)

	end <- num
	return
}

func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	//db.SetMaxOpenConns(200000)
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	ctx = context.Background()
	n := 1
	end := make(chan int)
	for i := 1; i <= n; i++ {
		go preorder(end, strconv.Itoa(i), 1, 5)
	}
	for i := 1; i <= n; i++ {
		<-end
	}
	return
}
