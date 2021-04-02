package main

import (
	//  "bytes"
	//  "encoding/gob"

	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Messagee struct {
	Name     string
	Quantity int
}

//209.97.165.170
var host = "209.97.165.170:8080"
var users = 10000

func send6(conn net.Conn, host string, m string, p string, userId int) {
	//fmt.Println("sent")
	userid++
	if m == "GET" {
		// fmt.Println("sent GET")
		fmt.Fprintf(conn, createH(m, p, userId))
	} else {
		fmt.Println("sent POST")
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

func client6(wg *sync.WaitGroup, m string, p string, userId int) {
	// t0 := time.Now()

	conn, err := net.Dial("tcp", host)
	if err != nil {
		count_Fail++
		log.Fatalln(err)
	}
	fmt.Println("sent", userId)
	send6(conn, host, m, p, userId)
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
	for i := 0; i < users; i++ {
		wg.Add(1)

		//client6(&wg, "GET", "/", i)
		client6(&wg, "GET", "/products", i)
		// client6(&wg, "GET", "/products/1", i)
		// client6(&wg, "POST", "/products/1", i)
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

func createH(methodd string, pathh string, u int) string {
	userID := u
	method := methodd
	path := pathh
	// host := "127.0.0.1:8080"
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
	// host := "127.0.0.1:8080"
	contentLength := 20
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
