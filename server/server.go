package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan string
}

var manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan string),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated")
			}
		case message := <-manager.broadcast:
			for connection := range manager.clients {
				select {
				case connection.data <- message:
				default:
					close(connection.data)
					delete(manager.clients, connection)
				}
			}
		}
	}
}

func (client *Client) receive() {
	for {
		message := <-client.data
		client.socket.Write([]byte(message))
	}
}

func handleConnections(conn net.Conn) {
	client := &Client{socket: conn, data: make(chan string)}
	manager.register <- client

	go client.receive()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("Message Received:", text)
		manager.broadcast <- text
	}

	manager.unregister <- client
	if scanner.Err() != nil {
		fmt.Println("Error reading from client:", scanner.Err())
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Server started on localhost:8080")

	go manager.start()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnections(conn)
	}
}
