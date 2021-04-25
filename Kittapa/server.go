package Kittapa

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ID int = 0

type (
	Server struct {
		routers map[string]*Route
	}

	Route struct {
		Method string      `json:"method"`
		Path   string      `json:"path"`
		Name   HandlerFunc `json:"name"`
	}

	LimitOffset struct {
		limit  int
		offset int
	}
)
type HandlerFunc func() string

type data struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type Message struct {
	Mess string `json:"mess"`
}

var count = 0
var Result data
var LF LimitOffset

func (s *Server) GET(path string, h HandlerFunc) *Route {
	m := "GET"
	return s.Add(m, path, h)
}

func (s *Server) POST(path string, h HandlerFunc) *Route {
	m := "POST"
	return s.Add(m, path, h)
}

func (s *Server) Add(m, path string, h HandlerFunc) *Route {
	r := &Route{
		Method: m,
		Path:   path,
		Name:   h,
	}
	s.routers[m+path] = r
	return r
}

func New() *Server {
	return &Server{
		routers: map[string]*Route{},
	}
}

func (s *Server) listen(port string) {
	li, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()

		if err != nil {
			log.Fatalln(err.Error())
			continue
		}
		count++
		fmt.Println("connections:", count)
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	s.req(conn)
}

func (s *Server) req(conn net.Conn) {
	buffer := make([]byte, 1024)
	var fc string
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(string(buffer[:n]))
		if !strings.Contains(string(buffer[:n]), "HTTP") {
			if _, err := conn.Write([]byte("Recieved\n")); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
			break
		}
		method, path, _ := getMessage(string(buffer[:n]))
		fmt.Println(method, path)
		r, yes := s.check(method, path)
		if yes {
			fmt.Println("yesssss")
			fc = r.Name()
			send(conn, fc, "application/json")
		} else {
			fmt.Println("no")
			// a := fmt.Sprintf("HTTP/1.0 404 Nof Found\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", 14, "text/html", "404 not found")
			// a := fmt.Sprintf("HTTP/1.0 404\r\n")
			// fmt.Fprintf(conn, a)
		}
		break
	}
}

func getMessage(message string) (string, string, []string) {
	Result = getJson(message)
	headers := strings.Split(message, "\n")
	// fmt.Println("headers", headers)
	// if len(headers) == 1 {
	// 	panic("len is 1")
	// }
	method := (strings.Split(headers[0], " "))[0]
	// fmt.Println("len:", len(headers))
	// fmt.Println("headers[0]", headers[0])
	path := (strings.Split(headers[0], " "))[1]
	p := strings.Split(path, "/")
	fmt.Println("len p:", len(p))
	fmt.Println("p[1]:", p[1])
	a, b := queryString(p[1])
	if b {
		fmt.Println(a)
		qString := strings.Split(a, "&")
		for j := 0; j < len(qString); j++ {
			fmt.Println("J: ", qString[j])
		}
		k := strings.Split(qString[0], "=")[1]
		LF.limit, _ = strconv.Atoi(k)
		k = strings.Split(qString[1], "=")[1]
		LF.offset, _ = strconv.Atoi(k)
		path = "/products"
	} else if p[1] == "products" && len(p) == 3 {
		fmt.Println("productsWithID")
		ID, _ = strconv.Atoi(p[2])
		path = "/" + p[1] + "/:id"
	}
	fmt.Println("path", path)
	return method, path, p
}

func queryString(m string) (string, bool) {
	index := 0
	b := false
	for i := 0; i < len(m); i++ {
		if m[i] == 63 {
			// fmt.Println("?")
			b = true
			index = i
		}
		// fmt.Println(m[i])
	}
	a := m[index+1:]
	return a, b
	// match, _ := regexp.MatchString(, m)
	// fmt.Println(match)
}

func (s *Server) check(method, path string) (*Route, bool) {
	value, exist := s.routers[method+path]
	return value, exist
}

func (s *Server) Start(port string) {
	fmt.Println(s)
	s.listen(port)
}

func send(conn net.Conn, d string, c string) {
	fmt.Fprintf(conn, createHeader(d, c))
}

//create header function
func createHeader(d string, contentType string) string {
	m := Message{Mess: d}
	a, _ := json.Marshal(m)
	contentLength := len(a)
	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, a)
	return headers
}

func getJson(message string) data {
	var result data
	if strings.ContainsAny(string(message), "}") {

		r, _ := regexp.Compile("{([^)]+)}")
		match := r.FindString(message)
		// fmt.Println(match)
		fmt.Printf("%T\n", match)
		json.Unmarshal([]byte(match), &result)
		// fmt.Println("data", result)
	}
	return result
}
