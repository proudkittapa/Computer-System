package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net"
	"os"
	"pin2pre/cacheFile"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/profile"
)

//helloproudd
type display struct {
	Product []string `json:"Product"`
}

const BUFFERSIZE = 1024

var mp map[int]string = make(map[int]string)
var cacheObject cacheFile.Cache = cacheFile.NewCache()
var lru cacheFile.Lru_cache = cacheFile.Cache_cons(10)

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

var (
	db          *sql.DB
	q           int
	newQuantity int
	mutex       sync.Mutex
)

type respond struct {
	Msg string `json:"msg"`
}

var count int = 0

//178.128.94.63

func main() {
	defer profile.Start().Stop()
	li, err := net.Listen("tcp", ":8080")
	db, _ = sql.Open("mysql", "root:62011139@tcp(localhost:3306)/prodj")
	// db.SetMaxIdleConns(200000)
	db.SetMaxOpenConns(200000)

	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()
	// defer db.Close()

	for {
		conn, err := li.Accept()

		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		count++
		fmt.Println("connections:", count)
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	req(conn)

}

func req(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		message := string(buffer[:n])
		// fmt.Println("mess", message)
		if !strings.Contains(message, "HTTP") {
			if _, err := conn.Write([]byte("Recieved\n")); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
			break
		}
		headers := strings.Split(message, "\n")
		method := (strings.Split(headers[0], " "))[0]
		path := (strings.Split(headers[0], " "))[1]
		p := strings.Split(path, "/")
		// fmt.Println(message)
		if p[1] == "" {
			home(conn, method, "pre-order/index.html", "text/html")
			break
		} else if p[1] == "text" {
			sendText(conn)
			break
		} else if p[1] == "payment" {
			receiveFile(conn)
			break
		} else if p[1] == "products" {
			if (len(p) > 2) && (p[2] != "") {
				fmt.Println("message", message)
				result := getJson(message)
				// fmt.Println("P2", p[2])
				productWithID(conn, method, p[2], result)
				break
			} else {
				// fmt.Println("HI")
				products(conn, method)
				break
			}
		} else if p[1] == "style.css" {
			home(conn, method, "pre-order/style.css", "text/css")
			break
		} else if p[1] == "images" {
			f := p[2]
			nf := "pre-order/images/" + f
			homeImg(conn, method, nf, "image/apng")
			break
		}
	}

}

func getJson(message string) data {
	var result data
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
func receiveFile(connection net.Conn) {
	// defer connection.Close()
	fmt.Println("Connected to server, start receiving file size")
	// bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	// fmt.Println("connection", connection)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 10)
	fmt.Println("fileSize", fileSize)
	mutex.Lock()
	newFile, err := os.Create("new.jpg")
	mutex.Unlock()
	if err != nil {
		panic(err)
	}
	defer newFile.Close()
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	if fileSize != 0 {
		send(connection, "Received file completely!", "text")
	} else {
		send2(connection, "429")
	}
	// fmt.Println("Received file completely!")
}

func homeImg(conn net.Conn, method string, filename string, t string) {
	if method == "GET" {
		c := t
		d, _ := getImageFromFilePath(filename)
		sendFile(conn, d, c)
	}
}

func sendText(conn net.Conn) {
	c := "text"
	d := "send text"
	send(conn, d, c)
	count--
	fmt.Println("connections after return:", count)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func home(conn net.Conn, method string, filename string, t string) {
	if method == "GET" {
		c := t
		d := cacheFile.Call_cache(filename)
		send(conn, d, c)
	}
}

func products(conn net.Conn, method string) {
	if method == "GET" {
		d := display_pro()
		c := "application/json"
		send(conn, d, c)
	}
}

func productWithID(conn net.Conn, method string, id string, result data) {
	// fmt.Println("ID")
	i, _ := strconv.Atoi(id)
	if method == "GET" {
		mutex.Lock()
		d := lru.Cache(i)
		mutex.Unlock()
		c := "application/json"
		send(conn, d, c)

	} else if method == "POST" {
		// fmt.Println("here")
		// fmt.Println(result.Quantity)
		success := cacheFile.PostPreorder(i, result.Quantity)
		// fmt.Println("sjdkfa;sd")
		msg := ""
		if success == true {
			msg = "success"
		} else {
			msg = "fail"
		}
		jsonStr := respond{Msg: msg}
		jsonData, err := json.Marshal(jsonStr)
		if err != nil {
			fmt.Println("error post", err)
		}
		d := string(jsonData)
		c := "application/json"
		send(conn, d, c)
	}
}

func sendFile(conn net.Conn, d image.Image, c string) {
	fmt.Fprintf(conn, createHeaderFile(d, c))
}

func createHeaderFile(d image.Image, contentType string) string {

	contentLength := 0

	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, d)
	// fmt.Println(headers)
	return headers
}

func send(conn net.Conn, d string, c string) {
	fmt.Fprintf(conn, createHeader(d, c))
}

//create header function
func createHeader(d string, contentType string) string {
	contentLength := len(d)
	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, d)
	return headers
}

func send2(conn net.Conn, h string) {
	fmt.Fprintf(conn, createHeader2(h))
}

//create header function
func createHeader2(httpStatus string) string {
	// contentLength := len(d)
	headers := fmt.Sprintf("HTTP/1.1 %s\r\n", httpStatus)
	return headers
}

func checkErr(err error) (a bool) {
	a = true
	if err != nil {
		fmt.Println("check err", err)
		a = false
	}
	return
}

func display_pro() (val string) {
	var l []string
	for i := 1; i <= 1; i++ {
		val := cacheFile.Db_query(i)
		l = append(l, val)
	}

	result := display{Product: l}

	byteArray, err := json.Marshal(result)
	checkErr(err)

	val = string(byteArray)
	// fmt.Println(val)
	return
}

func getQuantity(t chan int, id int) {
	start := time.Now()
	info := lru.Cache(id)

	var quan data
	err := json.Unmarshal([]byte(info), &quan)
	checkErr(err)
	t <- quan.Quantity

	fmt.Printf("time query from cache: %v\n", time.Since(start))
	// fmt.Println("Quantity: ", quan.Quantity)

}
