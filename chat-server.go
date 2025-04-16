// A simple TCP chat server

package main

import (
	"bufio"
	"io"
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

	// A list to keep track of all incoming connections
	var connections []net.Conn

	// Accept all incoming connections
	for {
		connection, err := listener.Accept()

		// Promt for username
		connection.Write([]byte("Enter username: "))

		connections = append(connections, connection)

		if err != nil {
			log.Fatal(err)
		}

		// Handle connections concurrently
		go broadcastMessage(connection, &connections)
	}
}

// Send a client's message to all connected clients
func broadcastMessage(connection net.Conn, connections *[]net.Conn) {
	defer connection.Close()
	for {
		message, err := bufio.NewReader(connection).ReadString('\n')

		// Handle client exits gracefully
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Loop through our connections list and forward message
		for _, element := range *connections {
			// Skip sender
			if element == connection {
				continue
			}
			element.Write([]byte(message + "\n"))
		}
	}
}
