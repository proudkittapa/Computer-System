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
	"strconv"
	// "time"
	// "bytes"
)

type result struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var count int = 0

func main() {

	li, err := net.Listen("tcp", ":8081")
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
	// conn.Write([]byte("Recieved\n"))
	// count++
	// fmt.Printf("%v", conn)

	
	// // fmt.Println("break")
	// defer conn.Close()
	
	req(conn)
	// if _, err := conn.Write([]byte("Recieved\n")); err != nil {
	// 	log.Printf("failed to respond to client: %v\n", err)
	// }
	
	
}

func req(conn net.Conn) {
	var data result
	// defer conn.Close()
	buffer := make([]byte, 1024)
	fmt.Printf("buffer %T", buffer)
	message := ""
	m := ""
	for {
		// fmt.Println("hihihiihihihi")

		n, err := conn.Read(buffer)
		if err != nil{
			fmt.Println(err)
		}
		// fmt.Printf("n value: %v, %T\n",n, n)
		
		// fmt.Println("hihihiihihihi")

		message = string(buffer[:n])
		// fmt.Println(n)
		fmt.Println("mess", message)
		if ! strings.Contains(message, "HTTP"){
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
			json.Unmarshal([]byte(match), &data)
			fmt.Println("data", data)
			fmt.Println("Name", data.Name)
			fmt.Println("Quantity", data.Quantity)

			fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			fmt.Fprintf(conn, "Content-Length: %d\r\n", len(data.Name)+1)
			fmt.Fprint(conn, "Content-Type: text/html\r\n")
			fmt.Fprint(conn, "\r\n")
			// q := strconv.Itoa(data.Quantity)
			fmt.Fprint(conn, data.Name)

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
			fmt.Println("Heyyyyyy")
			id = u[10:]
			// id := u
			productID(conn, id)
		}	
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
