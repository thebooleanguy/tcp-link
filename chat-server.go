// A simple TCP chat server

package main

import (
	"bufio"
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

	// Accept all incoming connections
	for {
		connection, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		// Handle connections concurrently
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	for {
		message, _ := bufio.NewReader(connection).ReadString('\n')
		connection.Write([]byte(message + "\n"))
	}
}
