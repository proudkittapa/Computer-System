package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	con, err := net.Dial("tcp", "0.0.0.0:8081")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	// serverReader := bufio.NewReader(con)
	var request = make([]byte, 100)

	for {
		log.Printf("Type something:")
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		default:
			log.Printf("client error: %v\n", err)
			return
		}
		// serverResponse, err := serverReader.ReadString('\n')

		// switch err {
		// case nil:
		// 	log.Println(strings.TrimSpace(serverResponse))
		// default:
		// 	log.Printf("server error: %v\n", err)
		// 	return
		// }
		_, err = con.Read(request)

		if err != nil {
			log.Println("failed to read request contents")
			return
		}
		log.Println(&con, string(request))
		request = make([]byte, 100)

	}
}
