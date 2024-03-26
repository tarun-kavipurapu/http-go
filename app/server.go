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
	// fmt.Println(request)
	handleContent(request, conn)

}

func handleContent(request string, conn net.Conn) {
	requestLines := strings.Split(request, "\r\n")
	path := strings.Split(requestLines[0], " ")[1]
	var responseSent []byte

	if path == "/" {
		responseSent = []byte("HTTP/1.1 200 OK\r\n\r\n")
	} else if strings.HasPrefix(path, "/echo") {
		// Here you can implement logic to echo back the request, if needed
		echoString := path[6:]
		// fmt.Println("echo", echoString)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(echoString))

		responseSent = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + contentLength + echoString)
	} else if path == "/user-agent" {
		requestLines := strings.Split(request, "\r\n")[2]
		userAgent := strings.Split(requestLines, ": ")[1]

		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(userAgent))
		// userAgentString:=fmt.Sprintf("User-Agent: %s",userAgent)
		responseSent = []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n" + contentLength + userAgent)
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
