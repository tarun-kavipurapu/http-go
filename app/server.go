package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func handleConnection(conn net.Conn) {

	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}
	request := string(buf[:n])
	requestLines := strings.Split(request, "\r\n")

	path := strings.Split(requestLines[0], " ")[1]

	var responseSent []byte
	if path == "/" {
		responseSent = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else {
		responseSent = []byte("HTTP/1.1 404 NOT FOUND\r\n\r\n")
	}

	_, errWrite := conn.Write(responseSent)

	if errWrite != nil {
		fmt.Println("Error writing:", errWrite)
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

		// fmt.Println("connection", conn.)
		handleConnection(conn)

	}

}
