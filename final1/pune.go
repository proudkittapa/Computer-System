package final1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var mp map[int]string = make(map[int]string)

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

// func cache(id int) string {
// 	// cache(1) // if true return val else return -1
// 	if val, ok := mp[id]; ok {
// 		fmt.Println("-----------HIT----------")
// 		fmt.Println(val)
// 		return val
// 	} else {
// 		return db_query(id)
// 		// return ""
// 	}
// }

func GetCache(id int) string {
	db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

	// fmt.Println("----------MISS----------")

	rows, err := db.Query("SELECT name, quantity_in_stock, unit_price FROM products WHERE product_id = " + strconv.Itoa(id))
	checkErr(err)

	for rows.Next() {
		var name string
		var quantity int
		var price int
		err = rows.Scan(&name, &quantity, &price)

		result := data{Name: name, Quantity: quantity, Price: price}
		byteArray, err := json.Marshal(result)
		checkErr(err)
		// fmt.Println(len(byteArray))

		mp[id] = string(byteArray)

	}
	val := mp[id]
	fmt.Println(val)
	return val

}

// func main() {

// 	for i := 0; i < 10; i++ {
// 		for j := 0; j < 2; j++ {
// 			db_query(i + 1)

// 		}
// 	}

// }
