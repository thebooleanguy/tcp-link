// A simple TCP chat server

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
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

	// A map to keep track of users
	users := make(map[string]net.Conn)

	// Accept all incoming connections
	for {
		connection, err := listener.Accept()

		// Promt for username
		go func() {
			connection.Write([]byte("Enter username: "))
			username, _ := bufio.NewReader(connection).ReadString('\n')
			fmt.Println(strings.TrimSpace(username) + " connected")

			// Store username and relevant socket in our map
			users[username] = connection
		}()
		time.Sleep(time.Second)

		if err != nil {
			log.Fatal(err)
		}

		// Handle connections concurrently
		go broadcastMessage(connection, users)
	}
}

// Send a client's message to all connected clients
func broadcastMessage(connection net.Conn, users map[string]net.Conn) {
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
		for _, val := range users {
			// Skip sender
			if val == connection {
				continue
			}
			val.Write([]byte(message + "\n"))
		}
	}
}
