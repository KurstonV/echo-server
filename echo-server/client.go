// Filename: main.go
// Filename: client.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the server on localhost:4000
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type messages and press Enter.")

	// Handle incoming messages from the server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("Server:", scanner.Text())
		}
	}()

	// Send user input to the server
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		_, err := fmt.Fprintln(conn, input.Text())
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
