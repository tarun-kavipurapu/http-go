package main

import (
	"flag"
	"fmt"
	"net"
)

var directory string

func main() {
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&directory, "directory", ".", "directory")
	flag.Parse()

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	handleError("Failed to bind to port 4221", err)

	defer l.Close()
	for {
		conn, err := l.Accept()
		handleError("Failed to accept the client", err)
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	requestBytes, err := conn.Read(buf)
	handleError("Error Reading", err)
	request := string(buf[:requestBytes])

	serverResponse := HandleRequest(request)

	_, err = conn.Write([]byte(serverResponse))
	handleError("Error	 writing response: ", err)
}
