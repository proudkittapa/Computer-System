package main

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"io/ioutil"
	"log"
	"net"
	"strings"
	// "time"
	// "bytes"
)

type result struct {
	Name string `json:"name"`
	Quantity int `json:"quantity"`
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
	data := result{}
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
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
		if strings.ContainsAny(string(message), "}"){
			
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
	
  
    // Prints output 
    // fmt.Printf("The number of bytes are: %d\n", bytes) 


	// scanner := bufio.NewScanner(conn)
	// i := 0
	// var contentLength int
	// for scanner.Scan() {
	// 	ln := scanner.Text()
	// 	fmt.Println(ln)
	// 	if i == 0 {
	// 		mux(conn, ln)
	// 	}
	// 	if i == 8{
	// 		fmt.Println("adsad")
	// 		fmt.Println(ln[16:])
	// 		contentLength, _ = strconv.Atoi(ln[16:])
	// 	}
	// 	if ln == "" {
	// 		//headers are done
	// 		fmt.Println("asdasdsss")
	// 		tmp := make([]byte, contentLength+100)     // using small tmo buffer for demonstrating
	// 		for {
	// 		    n, err := conn.Read(tmp)
	// 		    if err != nil {
	// 		        if err != io.EOF {
	// 		            fmt.Println("read error:", err)
	// 		        }
	// 		        break
	// 		    }
	// 			fmt.Println("got", n, "bytes.")
	// 			fmt.Println(string(tmp[:n]))
	// 		    // buf = append(buf, tmp[:n]...)

	// 		}
	// 		// fmt.Println("total size:", len(buf))
	// 		// fmt.Println(string(temp))
	// 		fmt.Println("break")
	// 		break
	// 	}
	// 	i++
	// }

    // fmt.Fprintf(conn, "GET / HTTP/2.0\r\n\r\n")
	// scanner := bufio.NewScanner(conn)
	// for i := 0; i < 15; i++ {
	// 	scanner.Scan()
	// 	ln := scanner.Text()
	// 	fmt.Println(ln)
	// 	// if strings.ContainsAny(ln, "}"){
	// 	// 	fmt.Println("asd")
	// 	// 	fmt.Println(ln)
	// 	// 	break
	// 	// }
	// }
	// fmt.Println("asdsadasd")

	// defer conn.Close()
	// bs, err := ioutil.ReadAll(conn)
	// fmt.Println(len(conn))

	// lr := io.LimitReader(conn, 254 + 24)
	// temp, err := ioutil.ReadAll(lr)
	// if err != nil{
	// 	log.Println(err)
	// }
	// fmt.Println(string(temp))

	// var data result
	// fmt.Println(json.NewDecoder(conn).Decode(&data))

	// defer conn.Close()
	// // fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	// // fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
    // buf := make([]byte, 0, 4096) // big buffer
    // tmp := make([]byte, 291)     // using small tmo buffer for demonstrating
    // for {
    //     n, err := conn.Read(tmp)
    //     if err != nil {
    //         if err != io.EOF {
    //             fmt.Println("read error:", err)
    //         }
    //         break
    //     }
    //     fmt.Println("got", n, "bytes.")
    //     buf = append(buf, tmp[:n]...)

    // }
	// fmt.Println("total size:", len(buf))
	// fmt.Println(string(buf))

	// defer conn.Close()
    // // fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
    // var buf bytes.Buffer
    // io.Copy(&buf, conn)
    // fmt.Println("total size:", buf.Len())


	// defer conn.Close()
	// buffer := make([]byte, 1024)
	// for {
	// 	n, err := conn.Read(buffer)
	// 	message := string(buffer[:n])

	// 	if message == "/quit" {
	// 		fmt.Println("quit command received. Bye.")
	// 		return
	// 	}

	// 	if n > 0 {
	// 		fmt.Println(message)
	// 	}

	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }

	// for {
	// 	message, err := bufio.NewReader(conn).ReadString('\n')
	// 	if err != nil {
	// 		log.Printf("Error: %+v", err.Error())
	// 		return
	// 	}

	// 	log.Println("Message:", string(message))
	// }


	// lr := io.LimitReader(conn, 277)
	// temp, err := ioutil.ReadAll(lr)

	// if err != nil{
	// 	log.Println(err)
	// }
	// fmt.Println(string(temp))

    //fmt.Println(string(buf))

	// scanner.Split(bufio.ScanBytes)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Bytes())
	// }
	// if scanner.Err() != nil {
	// 	fmt.Println(scanner.Err())
	// }
	
    // fmt.Println("Scanning ended")
	body := "asdasdasdads"
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))
	// conn.Close()
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
	// encoder := json.NewEncoder(conn)
	// decoder := json.NewDecoder(conn)

	// conn.Close()
	// var data result
	// decoder.Decode(&data)

	// encoder.Encode(data)
	// fmt.Println(data)
	// fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	// // fmt.Fprintf(conn, "Conten-Length: %d\r\n", len(body))
	// fmt.Fprint(conn, "Content-Type: text/html\r\n")
	// fmt.Fprint(conn, "\r\n")
	// // fmt.Fprint(conn, data)
	// fmt.Fprint(conn, string(body))

	body := "asdasd"
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, string(body))

}
