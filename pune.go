package main

import (
	// "bufio"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	// "time"
	// "bytes"
)

type result struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var count int = 0

func main() {

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
		count++
		fmt.Println("count error:", count)
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
	fmt.Println(conn)
	// defer conn.Close()
	req(conn)
}

func req(conn net.Conn) {
	data := result{}
	defer conn.Close()
	buffer := make([]byte, 1024)
	message := ""
	m := ""
	for {
		n, err := conn.Read(buffer)
		message = string(buffer[:n])
		// fmt.Println(n)
		// fmt.Println("mess", message)
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
			fmt.Println("match", match)
			json.Unmarshal([]byte(match), &data)
			fmt.Println("data", data)
			fmt.Println("Name", data.Name)
			fmt.Println("Quantity", data.Quantity)

			fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(conn, "Content-Length: %d\r\n", len(data.Name))
			fmt.Fprint(conn, "Content-Type: application/json\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, data.Name)

			// fmt.Println("break")
			break
		}
	}

	if m != "POST" {
		i := 0
		scanner := bufio.NewScanner(strings.NewReader(message))
		// fmt.Println("scan", scanner)
		// fmt.Println("mess", message)

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

}
func mux(conn net.Conn, ln string) {
	m := strings.Fields(ln)[0] //method
	u := strings.Fields(ln)[1] //url
	fmt.Println("***METHOD", m)
	fmt.Println("***URL", u)
	id := ""
	defer conn.Close()

	if m == "GET" && u == "/" {
		index(conn)
	}
	if m == "GET" && u == "/products" {
		product(conn)
	}

	if len(u) >= 10{
		if m == "GET" && u[:10] == "/products/" {
			id = u[10:]
			// id := u
			productID(conn, id)
		}	
	}
	/*
		if m == "POST" && u[:10] == "/products/" {
			id = u[10:]

			body := "asdasd"
			fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
			fmt.Fprint(conn, "Content-Type: application/json\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, string(body))

			// productPost(conn, id)
		}
	*/

}

func index(conn net.Conn) {
	body, err := ioutil.ReadFile("about_us.html")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println("Contents of file:", string(body))
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func product(conn net.Conn) {
	body := "products"
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func productID(conn net.Conn, id string) {
	body := id
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
}

func productPost(conn net.Conn, id string) {
	fmt.Println("herehere")
	data := result{}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		//fmt.Println(buffer)
		// status, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Println("status:", status)
		n, err := conn.Read(buffer)
		fmt.Println(buffer)
		message := string(buffer[:n])

		// fmt.Println(n)
		// fmt.Println(message)
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
			json.Unmarshal([]byte(match), &data)
			fmt.Println(data)
			fmt.Println(data.Name)
			fmt.Println(data.Quantity)
			// fmt.Println("break")
			break
		}
	}

	body := "asdasd"
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: application/json\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))

}
