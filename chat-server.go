// A simple TCP chat server

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	createServer()
}

func createServer() {
	// Create a TCP Listener interface
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal(err)
	}

	// Accept incoming connections
	connection, err := listener.Accept()

	if err != nil {
		log.Fatal(err)
	}

	// Continously print data from the connection
	for {
		msg, _ := bufio.NewReader(connection).ReadString('\n')
		fmt.Print(string(msg))
	}
}
