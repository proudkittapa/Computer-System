package final1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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

type Dis struct {
	Product []string
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

func Display_pro() (val string) {
	var l []string
	for i := 1; i <= 1; i++ {
		val := db_query(i)
		l = append(l, val)
	}

	result := Dis{Product: l}

	byteArray, err := json.Marshal(result)
	checkErr(err)

	val = string(byteArray)
	// fmt.Println(val)
	return
}

func ReCache(id int) string {
	// db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	// checkErr(err)

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

func GetFile() string {
	f, err := os.Open("/root/go/src/Computer-System/pre-order/index.html")

	if err != nil {
		fmt.Println("File reading error", err)

	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	chunksize := 512
	reader := bufio.NewReader(f)
	part := make([]byte, chunksize)
	buffer := bytes.NewBuffer(make([]byte, 0))
	var bufferLen int
	for {
		count, err := reader.Read(part)
		if err != nil {
			break
		}
		bufferLen += count
		buffer.Write(part[:count])
	}
	// fmt.Println("home")
	return buffer.String()
	// contentType = "text/html"
	// headers = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", bufferLen, contentType, buffer)

}

func db_query(id int) string {
	//db, err := sql.Open("mysql", "root:62011212@tcp(127.0.0.1:3306)/prodj")
	//checkErr(err)

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
