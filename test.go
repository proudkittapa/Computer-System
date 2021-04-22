package main

import (
	"context"
	"database/sql"
	"pin2pre/cacheFile"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
	//mutex sync.Mutex
	totalTime float64
)

func main() {
	c := cacheFile.Cache_cons(10)
	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	ctx = context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		panic(err)
	}
	//transactionC := make(chan string)
	t := make(chan int)
	cacheFile.GetQuantity(tx, t, 1)
	c.Display()

	// n := 2
	// end := make(chan int)
	// for i := 1; i <= n; i++ {
	// 	go cacheFile.Preorder(end, strconv.Itoa(i), 1, 5)
	// }
	// for i := 1; i <= n; i++ {
	// 	<-end
	// }
	return
}
