package Kittapa

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type (
	Server struct {
		routers map[string]*Route
	}

	Route struct {
		Method string      `json:"method"`
		Path   string      `json:"path"`
		Name   HandlerFunc `json:"name"`
	}
)
type HandlerFunc func() string

var count = 0

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
		method, path, _ := getMessage(string(buffer[:n]))
		fmt.Println(method, path)
		r, yes := s.check(method, path)
		if yes {
			fc = r.Name()
			send(conn, fc, "text/html")
		} else {
			// a := fmt.Sprintf("HTTP/1.0 404 Nof Found\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", 14, "text/html", "404 not found")
			// a := fmt.Sprintf("HTTP/1.0 404\r\n")
			// fmt.Fprintf(conn, a)
		}

	}
}

func getMessage(message string) (string, string, []string) {
	headers := strings.Split(message, "\n")
	// fmt.Println("headers", headers)
	// if len(headers) == 1 {
	// 	panic("len is 1")
	// }
	method := (strings.Split(headers[0], " "))[0]
	// fmt.Println("len:", len(headers))
	// fmt.Println("headers[0]", headers[0])
	path := (strings.Split(headers[0], " "))[1]
	// path := "path"
	p := strings.Split(path, "/")
	return method, path, p
}

func (s *Server) check(method, path string) (*Route, bool) {
	value, exist := s.routers[method+path]
	return value, exist
}

func (s *Server) Start(port string) {
	s.listen(port)
}

func send(conn net.Conn, d string, c string) {
	fmt.Fprintf(conn, createHeader(d, c))
}

//create header function
func createHeader(d string, contentType string) string {
	contentLength := len(d)
	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\n\n%s", contentLength, contentType, d)
	return headers
}
