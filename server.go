package main

import (
	"net"
	"fmt"
	"flag"
	"log"
	"bufio"
	"errors"
	"strings"
	"os"
)

var port int 


func getPathFromRequest(request string) (string, error) {
    parts := strings.Split(request, " ")
	if len(parts) < 3 {
		return "", errors.New("invalid request format")
	}
	return parts[1], nil
}

func getResponse(path string) (string, []byte) {
	data, err := os.ReadFile("www" + path)
	if err != nil {
		fmt.Println("Failed to read file", err)
		return "404 Not Found", []byte("")
	}
	return "200 OK", data
}

func handleConnection(c net.Conn) {
	fmt.Println("Handling connection", c)
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Failed to read from", c)
		}
		fmt.Println("Received message", msg)
		path, err := getPathFromRequest(msg)
		if err != nil {
			log.Fatalln("Failed to get path from request", c, err)
		}
		code, text := getResponse(path)
		_, e := c.Write([]byte(fmt.Sprintf("HTTP/1.1 %s\r\n\r\n%s\r\n", code, string(text))))
		if e != nil {
			log.Fatalln("Failed to response to the client", c, err)
		}
		break
	}
}

func main() {
	flag.IntVar(&port, "p", 9589, "Port on which server will listen")
	flag.Parse()

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalln("Error on connecting to listening port", port, err)
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Fatalln("Failed to accept connection", conn)
		}
		go handleConnection(conn)
	}
}