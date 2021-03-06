package main

import (
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	log.Println(&conn)
	defer conn.Close()
	var request = make([]byte, 100)
	for {
		_, err := conn.Read(request)

		if err != nil {
			log.Println("failed to read request contents")
			return
		}
		log.Println(&conn, string(request))
		request = make([]byte, 100)
		if _, err = conn.Write([]byte("Recieved\n")); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}

	}

}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Println("error listening on port 8081")
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		log.Println("received connection")
		if err != nil {
			log.Println("failed to accept connection")
			continue
		}
		go handleConnection(conn)
	}
}
