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

	// A map to keep track of users and sockets
	users := make(map[string]net.Conn)

	// A map to keep track of a user's connected room
	user_rooms := make(map[string]string)

	// A map to keep track of rooms and connected users
	rooms := make(map[string][]string)

	// Accept all incoming connections
	for {
		connection, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		// Promt for username
		go func() {
			var username string
			for {
				connection.Write([]byte("Enter username: "))
				username, _ = bufio.NewReader(connection).ReadString('\n')

				_, prs := users[username]
				if prs {
					connection.Write([]byte("Username not available, try something else :( \n"))
					continue
				}
				break
			}

			// Store username and relevant socket in our map
			users[username] = connection

			// Add user to lobby upon connection
			user_rooms[username] = "lobby"
			rooms["lobby"] = append(rooms["lobby"], username)
			fmt.Println(strings.TrimSpace(username) + " joined lobby")

			// Handle connections concurrently
			handleConnection(connection, username, users, user_rooms, rooms)
		}()
		time.Sleep(time.Second)

	}
}

// Handle client upon connection
func handleConnection(connection net.Conn, username string, users map[string]net.Conn, user_rooms map[string]string, rooms map[string][]string) {
	defer connection.Close()
	for {
		message, err := bufio.NewReader(connection).ReadString('\n')

		// Handle client exits gracefully
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Parse client commands
		if message[0] == '/' {
			parseCommands(message, username, user_rooms, rooms)
		}

		// Broadcast a client's message to all connected clients
		for _, val := range users {
			// Skip sender
			if val == connection {
				continue
			}
			val.Write([]byte(message + "\n"))
		}
	}
}

// Parse and respond to client commands
func parseCommands(message string, username string, user_rooms map[string]string, rooms map[string][]string) {

	words := strings.Fields(message)
	command := words[0]
	new_room := words[1]
	old_room := user_rooms[username]

	if command == "/join" {
		// Remove user from old room
		for _, val := range rooms[old_room] {
			if val == old_room {
				val = ""
				break
			}
		}
		// Add user to new room
		user_rooms[username] = new_room
		rooms[new_room] = append(rooms[new_room], username)
		fmt.Println(strings.TrimSpace(username) + " joined " + new_room)
	}
}
