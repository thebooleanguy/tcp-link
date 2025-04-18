// A simple TCP chat server

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net"
	"slices"
	"strings"
	"time"
)

func main() {
	handleDatabase()
	createServer()
}

func createServer() {
	// Create a TCP Listener interface
	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		log.Fatal(err)
	}

	// Note: Can probably reduce the number of maps used by using a custom user struct
	// A map to keep track of usernames and sockets
	users := make(map[string]net.Conn)

	// A map to keep track of a user's connected room
	user_rooms := make(map[string]string)

	// A map to keep track of rooms and connected clients
	rooms := make(map[string][]net.Conn)

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
			rooms["lobby"] = append(rooms["lobby"], connection)
			fmt.Println(strings.TrimSpace(username) + " joined lobby")

			// Handle connections concurrently
			handleConnection(connection, username, user_rooms, rooms)
		}()
		time.Sleep(time.Second)

	}
}

// Handle client upon connection
func handleConnection(connection net.Conn, username string, user_rooms map[string]string, rooms map[string][]net.Conn) {
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
			parseCommands(message, connection, username, user_rooms, rooms)
		}

		// Broadcast a client's message to all clients in our room
		for _, val := range rooms[user_rooms[username]] {
			// Skip commands
			if message[0] == '/' {
				break
			}
			// Skip sender
			if val == connection {
				continue
			}
			val.Write([]byte(message + "\n"))
		}
	}
}

// Parse and respond to client commands
func parseCommands(message string, connection net.Conn, username string, user_rooms map[string]string, rooms map[string][]net.Conn) {

	words := strings.Fields(message)
	command := words[0]
	new_room := words[1]
	old_room := user_rooms[username]

	if command == "/join" {
		// Remove user from old room
		for idx, val := range rooms[old_room] {
			if val == connection {
				rooms[old_room] = slices.Delete(rooms[old_room], idx, idx+1)
				break
			}
		}
		// Add user to new room
		user_rooms[username] = new_room
		rooms[new_room] = append(rooms[new_room], connection)
		fmt.Println(strings.TrimSpace(username) + " joined " + new_room)
	} else {
		connection.Write([]byte("Unknown command: " + command + " :( \n"))
	}
}

func handleDatabase() {

	db, err := sql.Open("sqlite3", "chat.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create table if it does not exist already
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, room TEXT NOT NULL, username TEXT NOT NULL, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")

	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
}
