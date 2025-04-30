package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from server: %v", err)
		}
		os.Exit(0)
	}()

	inputScanner := bufio.NewScanner(os.Stdin)
	for inputScanner.Scan() {
		text := inputScanner.Text()
		if text == "" {
			continue
		}
		_, err := fmt.Fprintln(conn, text)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
	if err := inputScanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
