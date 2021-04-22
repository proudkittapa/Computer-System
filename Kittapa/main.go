package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Handler func()

func Handle(conn net.Conn) {
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
		Receive(string(buffer[:n]))
	}

}

func main() {
	li, err := net.Listen("tcp", ":8080")
	// db.SetMaxIdleConns(200000)

	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()
	for {
		conn, err := li.Accept()

		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		// fmt.Println("connections:", count)
		go Handle(conn)
	}
}

func Receive(message string) (method string, pathS string, pathL []string) {
	// if !strings.Contains(message, "HTTP") {
	// 	if _, err := conn.Write([]byte("Recieved\n")); err != nil {
	// 		log.Printf("failed to respond to client: %v\n", err)
	// 	}
	// 	break
	// }
	fmt.Println(message)
	headers := strings.Split(message, "\n")
	method = (strings.Split(headers[0], " "))[0]
	pathS = (strings.Split(headers[0], " "))[1]
	pathL = strings.Split(pathS, "/")
	// fmt.Println("headers:", headers)
	fmt.Println("method:", method)
	fmt.Println("path:", pathS)
	fmt.Println("p:", pathL)
	return method, pathS, pathL
}

func GET(path string, do Handler) {
	do()
}

func d(d string) {
	fmt.Println("helooo")
}

func send(conn net.Conn, d string, c string) {
	fmt.Fprintf(conn, createHeader(d, c))
}

func createHeader(d string, contentType string) string {
	contentLength := len(d)
	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, d)
	return headers
}
