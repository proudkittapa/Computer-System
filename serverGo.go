package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

type result struct {
	Name string `json:"name"`
}

func main() {
	li, err := net.Listen("tcp", ":8080")
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
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	req(conn)
}

func req(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			mux(conn, ln)
		}
		if ln == "" {
			//headers are done
			break
		}
		i++
	}

}
func mux(conn net.Conn, ln string) {
	m := strings.Fields(ln)[0] //method
	u := strings.Fields(ln)[1] //url
	fmt.Println("***METHOD", m)
	fmt.Println("***URL", u)

	if m == "GET" && u == "/" {
		index(conn)
	}
	if m == "GET" && u == "/products" {
		product(conn)
	}

	if m == "GET" && u[:10] == "/products/" {
		id := u[10:]
		// id := u
		productID(conn, id)
	}
	if m == "POST" && u[:10] == "/products/" {
		id := u[10:]
		// id := u
		productPost(conn, id)
	}

}

func index(conn net.Conn) {

	body, err := ioutil.ReadFile("about_us.html")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(body))
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Conten-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func product(conn net.Conn) {
	body := "products"
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Conten-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func productID(conn net.Conn, id string) {
	body := id
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Conten-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func productPost(conn net.Conn, id string) {
	// body := id
	data, err := bufio.NewReader(conn).ReadBytes(0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("read %v bytes from the server\n", len(data))
	fmt.Println("data: ", string(data))
	var obj result
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println(err)
	}
	// conn.Close()
	//fmt.Println(msg)
	// fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	// fmt.Fprintf(conn, "Conten-Length: %d\r\n", jsonValue)
	// fmt.Fprint(conn, "Content-Type: text/json\r\n")
	// fmt.Fprint(conn, "\r\n")
	// fmt.Fprint(conn, jsonValue)
}
