package main

import (
	"context"
	"database/sql"
	"pin2pre/cacheFile"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
	//mutex sync.Mutex
	totalTime float64
)

func main() {
	cacheFile.InitDatabase()
	cacheFile.InitCache()
	//db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)

	n := 5
	end := make(chan int)
	for i := 1; i <= n; i++ {
		go cacheFile.Preorder(end, strconv.Itoa(i), 1, 1)
		//time.Sleep(1 * time.Second)
	}

	for i := 1; i <= n; i++ {
		<-end
	}

	return
}
