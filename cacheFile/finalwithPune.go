package cacheFile

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"regexp"
	"strconv"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx       context.Context
	db        *sql.DB
	mutex     sync.Mutex
	totalTime float64
	x         Data
	result    string
)

type product struct {
	Name     string
	Quantity int
	Price    int
}

// func InitCache() {
// 	c = Cache_cons(10)
// }
func InitDatabase() {
	// db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	db, _ = sql.Open("mysql", "root:62011139@tcp(127.0.0.1:3306)/prodj")
	db.SetMaxOpenConns(3000)
	// db.SetMaxIdleConns(32000)
	// db.SetConnMaxLifetime(10 * time.Second)
	for i := 1; i <= 5; i++ {
		db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, i)
	}

}

func getJson(message string) product {
	var result product
	if strings.ContainsAny(string(message), "}") {

		r, _ := regexp.Compile("{([^)]+)}")
		match := r.FindString(message)
		// fmt.Println(match)
		fmt.Printf("%T\n", match)
		json.Unmarshal([]byte(match), &result)
		// fmt.Println("data", result)
	}
	return result
}

func GetQuantity(tx *sql.Tx, transactionC chan string, t chan int, id int) {
	//fmt.Println("stop1")
	rows := tx.QueryRow("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	var name string
	var quantity int
	var price int
	err := rows.Scan(&name, &quantity, &price)
	if err != nil {
		//fmt.Println("get quantity fail")
		transactionC <- "rollback"
		tx.Rollback()
		return
	}
	x = Data{Name: name, Quantity: quantity, Price: price}

	//C.Set(id, x)
	//fmt.Println("stop2")
	//fmt.Println(val)
	//fmt.Printf("Name: %s, Quantity: %d\n", name, quantity)
	//fmt.Println("done")
	//fmt.Println(quantity)
	t <- quantity
}

func Decrement(tx *sql.Tx, t chan int, transactionC chan string, orderQuantity int, id int) {
	quantity := <-t // channel from getQuantity
	newQuantity := quantity - orderQuantity
	if quantity == 0 {
		transactionC <- "out of stock"
		return
	}
	if newQuantity < 0 {
		//fmt.Println("the order is out of stock")
		transactionC <- "not complete"
		return
	}
	_, err := tx.ExecContext(ctx, "update products set quantity_in_stock = ? where product_id = ? ", newQuantity, strconv.Itoa(id))
	if err != nil {
		//fmt.Println("decrement fail")
		tx.Rollback()
		transactionC <- "rollback"
		return
	}
	x = Data{Quantity: newQuantity}
	val := C.Set(id, x)
	fmt.Println(val)
	transactionC <- "done"
}

func Insert(wg *sync.WaitGroup, tx *sql.Tx, user string, id int, q int) {
	_, err := tx.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
	if err != nil {
		fmt.Println("insert fail")
		tx.Rollback()
		return
	}
	wg.Done()
}

func Preorder(end chan string, user string, productId int, orderQuantity int) {
	ctx = context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Println(err)
		// panic(err)
	}
	transactionC := make(chan string)
	t := make(chan int)
	//start := time.Now()
	go GetQuantity(tx, transactionC, t, productId)
	go Decrement(tx, t, transactionC, orderQuantity, productId)
	result2 := <-transactionC
	if result2 == "rollback" {
		//fmt.Println("rollback")
		Preorder(end, user, productId, orderQuantity)
		return
	} else if result2 == "not complete" {
		result = "order more than stock quantity"
		fmt.Println(result)
		tx.Commit()
		end <- result
		return
	} else if result2 == "out of stock" {
		result = "The order is out of stock"
		fmt.Println(result)
		tx.Commit()
		end <- result
		return
	} else {
		fmt.Println("user:", user, "productId:", productId, "orderQuantity:", orderQuantity)
		var wg sync.WaitGroup
		wg.Add(1)
		go Insert(&wg, tx, user, productId, orderQuantity)
		wg.Wait()
		if err := tx.Commit(); err != nil {
			// fmt.Printf("Failed to commit tx: %v\n", err)
		}

		result = "transaction successful"
		end <- result
	}

	C.Display()
	return
}
func PostPreorder(id int, quantity int) string {
	// InitDatabase()
	//InitCache()
	//db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
	//n := 5
	end := make(chan string)
	go Preorder(end, strconv.Itoa(1), id, quantity)
	fmt.Printf("quantityyyy: %d\n", quantity)
	result = <-end
	fmt.Println("hererreerere")
	return result
}

//fmt.Println("success")
//fmt.Println("-----------------------------------")
//elapsed := time.Since(start)
//tt := float64(elapsed)
//fmt.Printf("time: %v\n", elapsed)
//fmt.Printf("tt: %v\n", tt)
//totalTime += tt
//fmt.Printf("total time: %v\n", totalTime)
//num, _ := strconv.Atoi(user)
