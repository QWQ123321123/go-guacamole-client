package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/QWQ123321123/go-guacamole-client/client"
)

func main() {
	// Create tunnel with default config (localhost:4822)
	config := client.DefaultConfig()
	tunnel := client.NewTunnel(config)

	// Connect to target server via SSH
	fmt.Println("Connecting to SSH server...")
	err := tunnel.Connect(
		"your-server.com", // hostname
		"username",        // username
		"password",        // password
		22,                // port
		"ssh",             // protocol
		1024,              // width
		768,               // height
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer tunnel.Close()

	fmt.Println("Connected successfully!")

	// Simple read/write loop
	go func() {
		for {
			data, err := tunnel.ReadSome()
			if err != nil {
				if err == io.EOF {
					fmt.Println("Connection closed")
					return
				}
				log.Printf("Read error: %v", err)
				return
			}
			if len(data) > 0 {
				fmt.Printf("Received: %s", string(data))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Write a command
	fmt.Println("Sending 'ls' command...")
	_, err = tunnel.Write([]byte("ls\n"))
	if err != nil {
		log.Fatalf("Failed to write: %v", err)
	}

	// Wait a bit for response
	time.Sleep(2 * time.Second)
}
