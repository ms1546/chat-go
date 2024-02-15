package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server at localhost:8080")
	go readMessage(conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Message to send: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "QUIT" {
			break
		}
		conn.Write([]byte(text + "\n"))
	}
}

func readMessage(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println("Received from server:", text)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from server:", err.Error())
	}
}
