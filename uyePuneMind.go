package main

import (
	// "bufio"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	// "strconv"
	// "time"
	// "bytes"
	_ "github.com/go-sql-driver/mysql"
)

var u string
var ID int

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type respond struct {
	Msg string `json:"msg"`
}

type display struct {
	Product []string `json:"Product"`
}

var count int = 0

var (
	db          *sql.DB
	q           int
	newQuantity int
)

func main() {

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
		// fmt.Println("count error:", count)
	}
	defer li.Close()
	for {
		conn, err := li.Accept()

		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		go handle(conn)
	}

}

func handle(conn net.Conn) {
	defer conn.Close()

	req(conn)

}

func req(conn net.Conn) {
	var result data
	// defer conn.Close()
	buffer := make([]byte, 1024)
	fmt.Printf("buffer %T", buffer)
	message := ""
	m := ""
	for {

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Printf("n value: %v, %T\n",n, n)

		message = string(buffer[:n])
		// fmt.Println(n)
		fmt.Println("mess", message)
		if !strings.Contains(message, "HTTP") {
			if _, err := conn.Write([]byte("Recieved\n")); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
			break
		}

		if len(message) != 0 {
			m = message[:4]
			if m != "POST" {
				// fmt.Println(len(m))
				break
			}
			// fmt.Println("mess", m)
		}

		// totalBytes += n
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s", err)
			}
			break
		}
		if strings.ContainsAny(string(message), "}") {

			r, _ := regexp.Compile("{([^)]+)}")
			// match, _ := regexp.MatchString("{([^)]+)}", message)
			// fmt.Println(r.FindString(message))
			match := r.FindString(message)
			fmt.Println(match)
			// match = "`\n"+match+"\n`"
			fmt.Printf("%T\n", match)
			json.Unmarshal([]byte(match), &result)
			fmt.Println("data", result)
			fmt.Println("Name", result.Name)
			fmt.Println("Quantity", result.Quantity)
			fmt.Println("Price", result.Price)
			send(conn, "POST", "preorder")

			// fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			// fmt.Fprintf(conn, "Content-Length: %d\r\n", len(data.Name)+1)
			// fmt.Fprint(conn, "Content-Type: text/html\r\n")
			// fmt.Fprint(conn, "\r\n")
			// // q := strconv.Itoa(data.Quantity)
			// fmt.Fprint(conn, data.Name)

			// fmt.Println("break")
			break
		}
	}
	if m != "POST" {
		// fmt.Println("hihriehiehr")

		i := 0
		scanner := bufio.NewScanner(strings.NewReader(message))
		// fmt.Println("scan", scanner)
		// fmt.Println("mess", message)

		for scanner.Scan() {
			ln := scanner.Text()
			fmt.Println(ln)
			if i == 0 {
				fmt.Println("mux")
				mux(conn, ln)
			}
			if ln == "" {
				//headers are done
				break
			}
			i++
		}
	}

}

func mux(conn net.Conn, ln string) {
	m := strings.Fields(ln)[0] //method
	u = strings.Fields(ln)[1]  //url
	fmt.Println("***METHOD", m)
	fmt.Println("***URL", u)
	// id := ""
	defer conn.Close()

	if m == "GET" && u == "/" {
		// index(conn) //uye
		send(conn, "GET", "getHome")
	}
	if m == "GET" && u == "/products" {
		// product(conn) //mind return name, quantity, price
		//พราวต้องเอาใส่ json.Marshal
		send(conn, "GET", "getProducts")
	}

	if len(u) >= 10 {
		if m == "GET" && u[:10] == "/products/" {
			id := u[10:]
			fmt.Printf("%v, %T/n", id, id)
			i, _ := strconv.Atoi(id)
			ID = i
			send(conn, "GET", "getId")
		}
	}
}

func send(conn net.Conn, method string, res string) {
	fmt.Fprintf(conn, createHeader(method, res))
}

//create header function
func createHeader(method string, res string) string {

	d := ""
	headers := ""

	contentType := "application/json"
	if method == "POST" {
		fmt.Println("post")
		postPreorder(ID)
		jsonStr := respond{Msg: "POSTed"}
		jsonData, err := json.Marshal(jsonStr)
		if err != nil {
			fmt.Println(err)
		}
		d = string(jsonData)
	} else if (method == "GET") && (res == "getHome") {
		f, err := os.Open("about_us.html")

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
		fmt.Println("home")
		contentType = "text/html"
		headers = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", bufferLen, contentType, buffer)

		//uye
	} else if (res == "getId") && (method == "GET") {
		fmt.Println("id")
		d = db_query(ID)

	} else if res == "getProducts" {
		fmt.Println("products")
		d = display_pro()

		//Mind
	}

	contentLength := len(d)
	if res != "getHome" {
		headers = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, string(d))
	}
	fmt.Println(headers)
	return headers
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func db_query(id int) (val string) {
	db, err := sql.Open("mysql", "root:62011139@tcp(127.0.0.1:3306)/prodj")
	checkErr(err)

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

		val = string(byteArray)
		// fmt.Println(val)
	}
	return
}

func display_pro() (val string) {
	var l []string
	for i := 1; i <= 10; i++ {
		val := db_query(i)
		l = append(l, val)
	}

	result := display{Product: l}

	byteArray, err := json.Marshal(result)
	checkErr(err)

	val = string(byteArray)
	fmt.Println(val)
	return
}

func getQuantity(id int) {
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
func decrement(orderQuantity int, id int) {
	newQuantity := q - orderQuantity
	if newQuantity < 0 {
		return
	}
	fmt.Println("new quantity: ", newQuantity)
	db.Query("update products set quantity_in_stock = ? where product_id = ? ", newQuantity, id)

	return
}

func insert(user string, id int, q int) {
	db.Query("INSERT INTO order_items(username, product_id, quantity) VALUES (?, ?, ?)", user, id, q)
}

func preorder(user string, productId int, orderQuantity int) {
	//start := time.Now()
	insert(user, productId, orderQuantity)
	getQuantity(productId)
	decrement(orderQuantity, productId)
	//fmt.Printf("time: %v\n", time.Since(start))
	return
}

func postPreorder(id int) {
	db, _ = sql.Open("mysql", "root:62011139@tcp(127.0.0.1:3306)/prodj")
	for i := 1; i < 5; i++ {
		preorder(strconv.Itoa(i), id, 5) //userID, ID, quantity
	}
}
