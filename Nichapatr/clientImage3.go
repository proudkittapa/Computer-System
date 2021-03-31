package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

func main() {
	connection, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	// defer connection.Close()
	sendFileToServer(connection)
	receive(connection)
}

const BUFFERSIZE = 1024

func sendFileToServer(connection net.Conn) {
	file, err := os.Open("IMG_4.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	// fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileSize := strconv.FormatInt(fileInfo.Size(), 10)
	// fileName := fillString(fileInfo.Name(), 64)
	// var size int64 = fileInfo.Size()
	// fileSize := make([]byte, size)
	fmt.Println("Sending filesize!")
	connection.Write([]byte(fileSize))
	// connection.Write([]byte(fileName))
	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent")
	return
}

// func fillString(retunString string, toLength int) string {
// 	for {
// 		lengtString := len(retunString)
// 		if lengtString < toLength {
// 			retunString = retunString + ":"
// 			continue
// 		}
// 		break
// 	}
// 	return retunString
// }

func receive(connection net.Conn) {
	defer connection.Close()
}
