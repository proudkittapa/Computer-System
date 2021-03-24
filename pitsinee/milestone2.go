package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db    *sql.DB
	mutex sync.Mutex
)

func getQuantity(t chan int, id int) {

	row, err := db.Query("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	defer row.Close()
	for row.Next() {
		var name string
		var quantity int
		var price int
		row.Scan(&name, &quantity, &price)
		//a, _ := strconv.Atoi(quantity)
		//fmt.Println("a: ", a)
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
	fmt.Println("Product left in stock: ", newQuantity)
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)
	transactionC <- 0
}

func insert(user string, id int, q int) {
	db.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func preorder(end chan int, user string, productId int, orderQuantity int) {
	// fmt.Printf("start\n")
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	start := time.Now()
	transactionC := make(chan int)
	t := make(chan int)
	mutex.Lock()
	go getQuantity(t, productId)
	go decrement(t, transactionC, orderQuantity, productId)
	<-transactionC // wait for all go routines
	mutex.Unlock()
	go insert(user, productId, orderQuantity)
	fmt.Printf("time: %v\n", time.Since(start))

	num, _ := strconv.Atoi(user)
	end <- num
	return
}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	//defer db.Close()
	n := 150
	end := make(chan int, n)
	start2 := time.Now()

	for i := 0; i < n; i++ {
		go preorder(end, strconv.Itoa(i), 1, 1)
	}
	for i := 0; i < n; i++ {
		<-end
	}
	fmt.Printf("Total time: %v\n", time.Since(start2))

	return
}
