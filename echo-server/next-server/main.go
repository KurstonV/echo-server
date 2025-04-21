package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	maxMessageLength  = 1024
	inactivityTimeout = 30 * time.Second
)

var port = flag.String("port", "4000", "Port to listen on")

var logMutex sync.Mutex

func logMessage(ip string, message string) {
	logMutex.Lock()
	defer logMutex.Unlock()
	file, err := os.OpenFile(ip+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	fmt.Fprintf(file, "%s: %s\n", time.Now().Format(time.RFC3339), message)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	ip := conn.RemoteAddr().String()
	fmt.Printf("Client connected: %s at %s\n", ip, time.Now().Format(time.RFC3339))
	logMessage(ip, "Connected")

	scanner := bufio.NewScanner(conn)
	buffer := make([]byte, maxMessageLength)
	scanner.Buffer(buffer, maxMessageLength)

	inactive := time.NewTimer(inactivityTimeout)
	msgChan := make(chan string)
	quitChan := make(chan struct{})

	go func() {
		for scanner.Scan() {
			msg := strings.TrimSpace(scanner.Text())
			msgChan <- msg
			inactive.Reset(inactivityTimeout)
		}
		close(quitChan)
	}()

	for {
		select {
		case msg := <-msgChan:
			logMessage(ip, msg)
			if len(msg) > maxMessageLength {
				msg = msg[:maxMessageLength]
			}

			if msg == "" {
				conn.Write([]byte("Say something. . .\n"))
			} else if msg == "hello" {
				conn.Write([]byte("Hi there!\n"))
			} else if msg == "bye" {
				conn.Write([]byte("Goodbye!\n"))
				logMessage(ip, "Disconnected (said bye)")
				fmt.Printf("Client disconnected: %s at %s\n", ip, time.Now().Format(time.RFC3339))
				return
			} else if strings.HasPrefix(msg, "/") {
				switch {
				case msg == "/time":
					conn.Write([]byte(time.Now().Format(time.RFC1123) + "\n"))
				case msg == "/quit":
					conn.Write([]byte("Goodbye!\n"))
					logMessage(ip, "Disconnected (/quit)")
					fmt.Printf("Client disconnected: %s at %s\n", ip, time.Now().Format(time.RFC3339))
					return
				case strings.HasPrefix(msg, "/echo "):
					conn.Write([]byte(msg[6:] + "\n"))
				default:
					conn.Write([]byte("Unknown command\n"))
				}
			} else {
				conn.Write([]byte(msg + "\n"))
			}

		case <-inactive.C:
			conn.Write([]byte("Disconnected due to inactivity.\n"))
			logMessage(ip, "Disconnected (timeout)")
			fmt.Printf("Client disconnected due to timeout: %s at %s\n", ip, time.Now().Format(time.RFC3339))
			return

		case <-quitChan:
			logMessage(ip, "Disconnected")
			fmt.Printf("Client disconnected: %s at %s\n", ip, time.Now().Format(time.RFC3339))
			return
		}
	}
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on :%s\n", *port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConnection(conn)
	}
}
