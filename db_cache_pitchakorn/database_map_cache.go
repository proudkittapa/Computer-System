package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

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

// func cache(id int) (val string, res int) {

// 	_, in_cache := mp[id]

// 	if in_cache == true {
// 		// fmt.Println(mp[id])
// 		fmt.Println("-----------HIT----------")
// 		val = mp[id]
// 		res = 1
// 		// fmt.Println(val)
// 		return val, 0
// 	} else {
// 		val = "no"
// 		res = -1
// 		return
// 	}
// }

func cache(id int) string {
	// cache(1) // if true return val else return -1
	if val, ok := mp[id]; ok {
		fmt.Println("-----------HIT----------")
		fmt.Println(val)
		return val
	} else {
		return db_query(id)
		// return ""
	}
}

func db_query(id int) string {
	db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

	fmt.Println("----------MISS----------")

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

	////////// size cal ///////////
	// buf := &bytes.Buffer{}
	// // Any writer, not just a buffer!
	// var out io.Writer = buf
	// cw := &CounterWr{Writer: out}
	// if err := json.NewEncoder(cw).Encode(val); err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Count: %d bytes, %d bits\n", cw.Count, cw.Count*8)
	// fmt.Printf("Verif: %d bytes, %d bits\n", buf.Len(), buf.Len()*8)
	//////////size cal////////////

	fmt.Println(val)
	return val
	// fmt.Println(mp)
	// fmt.Println("from data")

}

type CounterWr struct {
	io.Writer
	Count int
}

func (cw *CounterWr) Write(p []byte) (n int, err error) {
	n, err = cw.Writer.Write(p)
	cw.Count += n
	return
}

func main() {
	var intSlice [10]int
	sum := 0

	for i := 0; i < 10; i++ {
		for j := 0; j < 2; j++ {
			start := time.Now()
			cache(i + 1)
			end := time.Since(start)
			fmt.Printf("%v\n", end)
			intSlice[i] = int(end)
		}
	}

	for k := 0; k < 10; k++ {
		sum += (intSlice[k])
	}

	// fmt.Println(intSlice)
	avg := (float64(sum)) / (float64(10))
	fmt.Println("Average = ", avg, "µs")
}
