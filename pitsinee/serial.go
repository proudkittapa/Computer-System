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

func getQuantity(t chan int, id int) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		tx.Rollback()
		return
	}
	var name string
	var quantity int
	var price float32
	err = rows.Scan(&name, &quantity, &price)
	t <- quantity
	fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
	if err := tx.Commit(); err != nil {
		fmt.Printf("Failed to commit tx: %v\n", err)

	}
}

func decrement(t chan int, transactionC chan int, orderQuantity int, id int) {
	tx2, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	tx2.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if newQuantity < 0 {
		transactionC <- 0
		return
	}
	fmt.Println(newQuantity)

	_, err = tx2.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	if err != nil {
		//panic(err)
		tx2.Rollback()
		return
	}
	//fmt.Println("updated")
	transactionC <- 0
	if err := tx2.Commit(); err != nil {
		fmt.Printf("Failed to commit tx2: %v\n", err)
	}
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
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	//start := time.Now()
	transactionC := make(chan int)
	t := make(chan int)
	go getQuantity(t, productId)
	go decrement(t, transactionC, orderQuantity, productId)
	<-transactionC // wait for all go routines
	go insert(user, productId, orderQuantity)
	//fmt.Printf("time: %v\n", time.Since(start))
	num, _ := strconv.Atoi(user)
	end <- num
	return

}
func main() {
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	ctx = context.Background()

	end := make(chan int)
	for i := 1; i < 10; i++ {
		go preorder(end, strconv.Itoa(i), 1, 5)
	}
	for i := range end {
		fmt.Println(i)
	}
}
