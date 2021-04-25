package final1

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB
	q           int
	newQuantity int
	result      string
)

func GetQuantity(id int) {
	row, err := db.Query("select name, quantity_in_stock, unit_price from products where product_id = " + strconv.Itoa(id))
	if err != nil {
		panic(err)
	}
	for row.Next() {
		var name string
		var quantity int
		var price int
		row.Scan(&name, &quantity, &price)
		q = quantity
		fmt.Println("name: ", name, " quantity: ", quantity, " price: ", price)
	}
}
func Decrement(orderQuantity int, id int) {
	newQuantity := q - orderQuantity
	if newQuantity < 0 {
		result = "order more than stock quantity"
		return
	} else if q == 0 {
		result = "The order is out of stock"
		return
	} else {
		result = "transaction successful"
		fmt.Println("new quantity: ", newQuantity)
		db.Exec("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)
	}
	return

}

func Insert(user string, id int, q int) {
	db.Exec("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func Preorder(user string, productId int, orderQuantity int) {
	start := time.Now()
	GetQuantity(productId)
	Decrement(orderQuantity, productId)
	Insert(user, productId, orderQuantity)
	fmt.Printf("time: %v\n", time.Since(start))
	return
}

func PostPreorder(id int, quantity int) string {
	Preorder(strconv.Itoa(1), id, quantity)
	return result
}

// func main() {
// 	db, _ = sql.Open("mysql", "root:mind10026022@tcp(127.0.0.1:3306)/prodj")
// 	db.Exec("update products set quantity_in_stock = ? where product_id = ? ", 1000, 1)
// 	start := time.Now()
// 	n := 10
// 	for i := 0; i < n; i++ {
// 		Preorder(strconv.Itoa(i), 1, 1)
// 	}

// 	fmt.Printf("Total time: %v\n", time.Since(start))
// 	return

// }
