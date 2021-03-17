package main

import (
//  "bytes"
//  "encoding/gob"
"bufio"
 "fmt"
 "net"
 "encoding/json"
)

type Message struct {
 Name     string
 Quantity string
}

func send(conn net.Conn, host string) {
	fmt.Fprintf(conn, createHeader())
}

func recv(conn net.Conn) {
	fmt.Println("reading")
	message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print(message)
}

func main() {
	host := "localhost:8080"
 conn, _ := net.Dial("tcp", ":8080")
 send(conn, host)
 recv(conn)
  
}

func createHeader() string{
	method := "GET"
	path := "/"
	host := "127.0.0.1:8080"
	contentLength := 20
	contentType := "application/json"
	jsonStr := Message{Name:"mos", Quantity:"2"}
	jsonData, err := json.Marshal(jsonStr)
	if err != nil {
		fmt.Println(err)
	}
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", method, path, host, contentLength, contentType, string(jsonData))
	return headers
}
