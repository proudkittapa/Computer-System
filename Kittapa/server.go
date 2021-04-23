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
type HandlerFunc func()

var count = 0

func (s *Server) GET(path string, h HandlerFunc) *Route {
	m := "GET"
	h()
	return s.Add(m, path, h)
}

func (s *Server) POST(path string, h HandlerFunc) *Route {
	m := "POST"
	h()
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

func listen(port string) {
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
		go handle(conn)
	}
}

func handle(conn net.Conn) {
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
		method, path, _ := getMessage(string(buffer[:n]))
		fmt.Println(method, path)

	}
}

func getMessage(message string) (string, string, []string) {
	headers := strings.Split(message, "\n")
	method := (strings.Split(headers[0], " "))[0]
	path := (strings.Split(headers[0], " "))[1]
	p := strings.Split(path, "/")
	return method, path, p
}

func (s *Server) check(method, path string) (*Route, bool) {
	value, exist := s.routers[method+path]
	return value, exist
}

func (s *Server) Start(port string) {
	listen(port)
}
