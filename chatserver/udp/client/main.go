package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:9001")
	if err != nil {
		log.Fatalf("Failed to resolve server address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatalf("Failed to dial UDP server: %v", err)
	}
	defer conn.Close()

	go func() {
		buf := make([]byte, 2048)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading from server: %v", err)
				return
			}
			fmt.Println(strings.TrimSpace(string(buf[:n])))
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		_, err := conn.Write([]byte(text))
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
		if text == "/quit" || text == "bye" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
