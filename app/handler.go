package main

import (
	"fmt"
	"os"
	"strings"
)

type Request struct {
	Path    string
	Method  string
	Headers map[string]string
	Body    string
}
type Response struct {
	Header        string
	StatusCode    int
	ContentType   string
	ContentLength int
	Body          string
}

var CRLF = "\r\n"

func handleError(message string, err error) {
	if err != nil {
		fmt.Println(message)
		os.Exit(1)
	}
}

func createRequest(req string) Request {
	reqRows := strings.Split(req, CRLF)
	path := strings.Split(reqRows[0], " ")[1]
	method := strings.Split(reqRows[0], " ")[0]
	headers := extractHeaders(method, reqRows)
	var body string
	if method != "GET" {
		body = reqRows[len(reqRows)-1]
	}
	return Request{Path: path, Method: method, Headers: headers, Body: body}

}

func extractHeaders(method string, requestRows []string) map[string]string {
	headers := make(map[string]string)
	var lastHeader int
	if method == "GET" {
		lastHeader = len(requestRows)
	} else {
		lastHeader = len(requestRows) - 1
	}
	for _, str := range requestRows[1:lastHeader] {
		parts := strings.SplitN(str, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}
	return headers
}

func HandleRequest(req string) string {

	request := createRequest(req)
	var responseContent []string
	if request.Method == "GET" && request.Path == "/" {
		response := handleGet()
		responseContent = []string{response.Header, CRLF, CRLF}
	} else if request.Method == "GET" && strings.Split(request.Path, "/")[1] == "echo" {
		response := handleGetEcho(request)
		contentLength := fmt.Sprintf("Content-Length: %d", response.ContentLength)
		responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, CRLF, CRLF, response.Body, CRLF}
	} else if request.Method == "GET" && strings.Split(request.Path, "/")[1] == "user-agent" {
		response := handleGetUserAgent(request)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", response.ContentLength)
		responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, response.Body, CRLF}
	} else if request.Method == "GET" && strings.Split(request.Path, "/")[1] == "files" {
		response := handleGetFiles(request)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response.Body))
		if response.StatusCode == 200 {
			responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, response.Body, CRLF}
		} else {
			responseContent = []string{response.Header, CRLF, CRLF}
		}
	} else if request.Method == "POST" && strings.Split(request.Path, "/")[1] == "files" {
		response := handlePostFiles(request)
		contentLength := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(response.Body))
		if response.StatusCode == 200 {
			responseContent = []string{response.Header, CRLF, response.ContentType, CRLF, contentLength, response.Body, CRLF}
		} else {
			responseContent = []string{response.Header, CRLF, CRLF}
		}

	} else {
		responseContent = notFound404()
	}
	return strings.Join(responseContent, "")

}
func handlePostFiles(request Request) Response {
	uriParts := strings.Split(request.Path, "/")
	header := "HTTP/1.1 201 OK"
	contentType := "Content-Type: text/plain"
	writeFile(request.Body, directory+uriParts[2])
	return Response{
		Header:      header,
		ContentType: contentType,
	}
}

func notFound404() []string {
	header := "HTTP/1.1 404 Not Found "
	responseContent := []string{header, CRLF, CRLF}

	return responseContent
}
func handleGet() Response {
	header := "HTTP/1.1 200 OK"
	return Response{
		Header:     header,
		StatusCode: 200,
	}
}
func handleGetEcho(request Request) Response {
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	body := strings.TrimPrefix(request.Path, "/echo/")
	return Response{
		Header:        header,
		StatusCode:    200,
		ContentType:   contentType,
		ContentLength: len(body),
		Body:          body,
	}
}

func handleGetUserAgent(request Request) Response {
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	body := request.Headers["User-Agent"]
	return Response{
		header,
		200,
		contentType,
		len(body),
		body,
	}
}

func handleGetFiles(request Request) Response {
	uriParts := strings.Split(request.Path, "/")
	header := "HTTP/1.1 200 OK"
	contentType := "Content-Type: text/plain"
	var body string
	var statusCode int
	if _, statErr := os.Stat(directory + uriParts[2]); statErr == nil {
		contentType = "Content-Type: application/octet-stream"
		body = readFileContent(uriParts[2])
		statusCode = 200
	} else {
		header = "HTTP/1.1 404 Not Found "
		body = ""
		statusCode = 404
	}
	return Response{
		Header:        header,
		StatusCode:    statusCode,
		ContentType:   contentType,
		ContentLength: len(body),
		Body:          body,
	}
}
func readFileContent(filename string) string {
	fileContent, err := os.ReadFile(directory + filename)
	if err != nil {
		fmt.Println("Error while reading file content: ", err.Error())
		panic(err)
	}
	body := string(fileContent)
	return body
}
func writeFile(body string, filePath string) {
	fileContent := []byte(body)
	err := os.WriteFile(filePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error while writing file content: ", err.Error())
		panic(err)
	}
}
