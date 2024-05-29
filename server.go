package main

import (
	"net"
	"fmt"
	"flag"
	"log"
	"bufio"
	"errors"
	"strings"
)

var port int 


func getPathFromRequest(request string) (string, error) {
    parts := strings.Split(request, " ")
	if len(parts) < 3 {
		return "", errors.New("invalid request format")
	}
	return parts[1], nil
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
		_, e := c.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nRequested path: %s\r\n", path)))
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