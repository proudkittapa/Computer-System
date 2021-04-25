package final2

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() {
	// db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
	db, _ = sql.Open("mysql", "root:62011139@tcp(127.0.0.1:3306)/prodj")
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	//db.SetConnMaxLifetime(10 * time.Second)
	for i := 1; i <= 5; i++ {
		db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, i)
	}

}

var (
	db     *sql.DB
	mutex  sync.Mutex
	result string
)

func GetQuantity(t chan int, id int) {

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
func Decrement(t chan int, transactionC chan string, orderQuantity int, id int) {
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
	fmt.Println("Product left in stock: ", newQuantity)
	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)
	transactionC <- "done"
}

func insert(user string, id int, q int) {
	db.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func Preorder(end chan string, user string, productId int, orderQuantity int) {
	// fmt.Printf("start\n")
	start := time.Now()
	transactionC := make(chan string)
	t := make(chan int)
	mutex.Lock()
	go GetQuantity(t, productId)
	go Decrement(t, transactionC, orderQuantity, productId)
	result2 := <-transactionC
	mutex.Unlock()
	if result2 == "not complete" {
		result = "order more than stock quantity"
		fmt.Println(result)
		end <- result
		return
	} else if result2 == "out of stock" {
		result = "The order is out of stock"
		fmt.Println(result)
		end <- result
		return
	} else {
		go insert(user, productId, orderQuantity)
		fmt.Printf("time: %v\n", time.Since(start))
		result = "transaction successful"
		end <- result
		return
	}
}
func PostPreorder(id int, quantity int) string {
	end := make(chan string)
	go Preorder(end, strconv.Itoa(1), id, quantity)
	fmt.Printf("quantityyyy: %d\n", quantity)
	result = <-end
	fmt.Println("hererreerere")
	return result
}

// func main() {
// 	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
// 	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)

// 	//defer db.Close()
// 	n := 10
// 	end := make(chan string)
// 	start2 := time.Now()

// 	for i := 0; i < n; i++ {
// 		go Preorder(end, strconv.Itoa(i), 1, 1)
// 	}
// 	for i := 0; i < n; i++ {
// 		<-end
// 	}
// 	fmt.Printf("Total time: %v\n", time.Since(start2))
// 	return
// }
