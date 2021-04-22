package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"math"
)

type Messagee struct {
	Name     string
	Quantity int
	Price    int
}
type PayInfo struct {
	Name      string
	ProductID int
	Date      string
	Time      string
	imageName string
}

var wg sync.WaitGroup

var img_name string = "IMG_4.jpg"

func send6(conn net.Conn, host string, m string, p string, userid int) {
	// fmt.Println("sent")
	userid++
	if m == "GET" {
		// fmt.Println("sent GET")
		fmt.Fprintf(conn, createHG(p, userid))
	} else if m == "POSE" && p == "/payment" {
		// fmt.Println("sent POST, img")
		fmt.Fprintf(conn, createHPimg(conn, userid))
		time.Sleep(1 * time.Millisecond)
		send_file(conn)
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
		log.Println("failed to read contents", message)
		return
	} else {
		count_Res++
	}
	fmt.Print(message)
}

func client6(m string, p string) {
	// t0 := time.Now()
	host := "localhost:8080"
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	send6(conn, host, m, p, userid)
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
	path := "/products/" + strconv.Itoa(rand.Intn(10))
	host := "127.0.0.1:8080"

	contentType := "application/json"
	jsonStr := Messagee{Name: "mos", Quantity: 2}
	jsonData, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
	}
	contentLength := len(string(jsonData))
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
	// send_file(conn)
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
	fmt.Println("Send filesize!")
	conn.Write([]byte(fileSize))
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

func onerun() {
	client6("GET", "/")
	client6("GET", "/products")
	client6("GET", "/products/1")
	client6("POST", "/products/1")
	client6("POST", "/payment")
}

func main() {
	// flag.Parse()
	start := time.Now()
	/*------------Cache check (1)------------*/
	for i := 1; i < 6; i++ {
		t0 := time.Now()
		client6("GET", "/products/"+strconv.Itoa(i))
		t00 = float64(time.Since(t0)))/1e6/5
		fmt.Printf("Latency Time:   %v ", t00)
	}
	for i := 6; i < 11; i++ {
		t1 := time.Now()
		client6("GET", "/products/"+strconv.Itoa(i))
		t01 = float64(time.Since(t1)))/1e6/5
		fmt.Printf("Latency Time:   %v \n", t01)
	}
	for i := 6; i < 11; i++ {
		t2 := time.Now()
		client6("GET", "/products/"+strconv.Itoa(i))
		t02 = float64(time.Since(t2)))/1e6/5
		fmt.Printf("Latency Time:   %v \n", t02)
	}
	if math.Abs(t00 - t01) <= 1 {
		fmt.Println("miss?\n")
	} else {
		fmt.Println("something is not right(1) :")
		fmt.Println(t00 - t01, "\n")
	} 
	if t02 <= t01 {
		fmt.Println("faster\n")
	} else {
		fmt.Println("cache not make faster maybe not hit\n")
	}
	/*------------Cache check (2)------------*/
	for i := 0; i < 200; i++ {
		t3 := time.Now()
		client6("GET", "/products/"+strconv.Itoa(i))
		t03 = float64(time.Since(t3)))/1e6/5
		fmt.Printf("Latency Time:   %v ", t03)
	}
	// wg.Wait()
	// time.Sleep(100 * time.Millisecond)
	t := time.Since(start)
	fmt.Printf("\n \nTotal TIME: %v\n", t)
	fmt.Printf("Number Response: %d\n", count_Res)
	fmt.Printf("Number fail: %d\n", count_Fail)
	tt := float64(t) / 1e6
	rate := float64(count_Res) / (tt / 1000)
	fmt.Printf("Rate per Sec: %f", rate)
}
