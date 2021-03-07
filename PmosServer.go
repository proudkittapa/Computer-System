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
	// defer conn.Close()
	req(conn)
}

func req(conn net.Conn) {
	// defer conn.Close()
	// body := "asdasdasdads"
	// fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	// fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	// fmt.Fprint(conn, "Content-Type: text/html\r\n")
	// fmt.Fprint(conn, "\r\n")
	// fmt.Fprint(conn, string(body))
	// defer conn.Close()
	// dst := os.Stdout
	// bytes, err := io.Copy(dst, conn)
	// if err != nil {
	//     panic(err)
	// }

	// run loop forever (or until ctrl-c)
	// for {
	// 	// will listen for message to process ending in newline (\n)
	// 	message, _ := bufio.NewReader(conn).ReadString('\n')
	// 	// output message received
	// 	fmt.Print("Message Received:", string(message))
	// 	// sample process for string received
	// 	newmessage := strings.ToUpper(message)
	// 	// send new string back to client
	// 	conn.Write([]byte(newmessage + "\n"))
	//   }

	// var data result

	// body := "asdasdasdads"
	// fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	// fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	// fmt.Fprint(conn, "Content-Type: text/html\r\n")
	// fmt.Fprint(conn, "\r\n")
	// fmt.Fprint(conn, string(body))
	// conn.Close()
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
		fmt.Println(id)
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
		fmt.Println(buffer)
		n, err := conn.Read(buffer)
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
