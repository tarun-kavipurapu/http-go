package main

import (
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()

	responseSent := []byte("HTTP/1.1 200 OK\r\n\r\n")

	_, err := conn.Write(responseSent)

	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to Accept the connection")
		}
		handleConnection(conn)

	}

}
