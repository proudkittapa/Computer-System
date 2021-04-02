package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Messagee struct {
	Name     string
	Quantity int
}
type PayInfo struct {
	Name      string
	ProductID int
	Date      string
	Time      string
	imageName string
}

var img_name string = "IMG_4.jpg"

func send6(conn net.Conn, host string, m string, p string) {
	fmt.Println("sent")
	userid++
	if m == "GET" {
		// fmt.Println("sent GET")
		fmt.Fprintf(conn, createHG(p, userid))
	} else if m == "POST" && p == "/payment" {
		// fmt.Println("sent POST, img")
		fmt.Fprintf(conn, createHPimg(conn, userid))
	} else {
		// fmt.Println("sent POST")
		fmt.Fprintf(conn, createHP(userid))
	}
}

func recv(conn net.Conn) {
	defer conn.Close()
	// fmt.Println("reading")
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		count_Fail++
		log.Println("failed to read contents")
		return
	}
	count_Res++
	// conn.Close()
	fmt.Print(message)
}

func client6(wg *sync.WaitGroup, m string, p string) {
	// t0 := time.Now()
	host := "localhost:8080"
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send6(conn, host, m, p)
	recv(conn)
	// fmt.Printf("Latency Time:   %v ", time.Since(t0))
	wg.Done()
	// <-ch
}

var userid = 0
var count_Res = 0
var count_Fail = 0

// var n = flag.Int("n", 5, "Number of goroutines to create")
// var ch = make(chan byte)

func main() {
	// flag.Parse()
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 200; i++ {
		wg.Add(1)
		client6(&wg, "POST", "/payment")
	}
	wg.Wait()
	// time.Sleep(100 * time.Millisecond)
	t := time.Since(start)
	fmt.Printf("\n \nTotal TIME: %v\n", t)
	fmt.Printf("Number Response: %d\n", count_Res)
	fmt.Printf("Number fail: %d\n", count_Fail)
	tt := float64(t) / 1e6
	rate := float64(count_Res) / (tt / 1000)
	fmt.Printf("Rate per Sec: %f", rate)
}

func createHG(pathh string, u int) string {
	userID := u
	method := "GET"
	path := pathh
	host := "127.0.0.1:8080"
	contentLength := 0
	contentType := "text"
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n userID:%d",
		method, path, host, contentLength, contentType, userID)
	return headers
}

func createHP(u int) string {
	userID := u
	method := "POST"
	path := "/products/1"
	host := "127.0.0.1:8080"
	contentLength := len(string(jsonData))
	contentType := "application/json"
	jsonStr := Messagee{Name: "mos", Quantity: 2}
	jsonData, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
	}
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s userID:%d",
		method, path, host, contentLength, contentType, string(jsonData), userID)
	return headers
}

func createHPimg(conn net.Conn, u int) string {
	userID := u
	method := "POST"
	path := "/payment"
	host := "127.0.0.1:8080"

	contentType := "image/jpg"
	jsonStr := PayInfo{Name: "Kanga", ProductID: 1123, Date: "20/02/21", Time: "12.00", imageName: img_name}
	jsonData, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
	}
	contentLength := len(string(jsonData))

	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s userID:%d",
		method, path, host, contentLength, contentType, string(jsonData), userID)
	send_file(conn)
	return headers
}

const BUFFERSIZE = 1024

func send_file(conn net.Conn) {
	file, err := os.Open(img_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	// fileSize := strconv.FormatInt(fileInfo.Size(), 10)
	// fileName := fillString(fileInfo.Name(), 64)
	// var size int64 = fileInfo.Size()
	// fileSize := make([]byte, size)
	fmt.Println("Send filesize!")
	conn.Write([]byte(fileSize))
	// connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	fmt.Println("File has been sent")
	return
}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
