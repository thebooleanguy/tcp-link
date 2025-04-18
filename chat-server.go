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
	sql_driver := "sqlite3"
	db_name := "chat.db"

	initDatabase(sql_driver, db_name)
	createServer(sql_driver, db_name)
}

func createServer(sql_driver string, db_name string) {
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
			joinRoom("lobby", connection, username, user_rooms, rooms, sql_driver, db_name)

			// Handle connections concurrently
			handleConnection(connection, username, user_rooms, rooms, sql_driver, db_name)
		}()
		time.Sleep(time.Second)

	}
}

// Handle client upon connection
func handleConnection(connection net.Conn, username string, user_rooms map[string]string, rooms map[string][]net.Conn, sql_driver string, db_name string) {
	// Handle database connection
	db, _ := sql.Open(sql_driver, db_name)
	defer db.Close()

	defer connection.Close()
	for {
		message, err := bufio.NewReader(connection).ReadString('\n')

		// Handle client exits gracefully
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Parse and respond to client commands
		if message[0] == '/' {
			words := strings.Fields(message)
			command := words[0]
			suffix := words[1]

			if command == "/join" {
				joinRoom(suffix, connection, username, user_rooms, rooms, sql_driver, db_name)
			} else {
				connection.Write([]byte("Unknown command: " + command + " :( \n"))
			}
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
		if message[0] != '/' {
			// Insert to database
			statement, _ := db.Prepare("INSERT INTO messages (room, username, content) VALUES (?, ?, ?)")
			statement.Exec(user_rooms[username], username, message)
		}
	}
}

// Remove client from its current room and add to another room
func joinRoom(new_room string, connection net.Conn, username string, user_rooms map[string]string, rooms map[string][]net.Conn, sql_driver string, db_name string) {

	// Handle database connection
	db, _ := sql.Open(sql_driver, db_name)
	defer db.Close()

	old_room := user_rooms[username]

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

	// Read room's last 10 messages from database
	rows, _ := db.Query("SELECT content FROM messages WHERE room = ? ORDER BY timestamp DESC LIMIT 10", new_room)
	var temp_msg string
	for rows.Next() {
		rows.Scan(&temp_msg)
		connection.Write([]byte(temp_msg + " \n"))
	}
}

// Create database and tables if they do not already exist
func initDatabase(sql_driver string, db_name string) {

	db, err := sql.Open(sql_driver, db_name)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Prepared statemnt to create table if it does not exist already
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY AUTOINCREMENT, room TEXT NOT NULL, username TEXT NOT NULL, content TEXT NOT NULL, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")

	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
}
