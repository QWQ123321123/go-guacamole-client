package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"

	"github.com/QWQ123321123/go-guacamole-client/client"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (for development only)
	},
}

func main() {
	http.HandleFunc("/terminal", handleTerminal)
	fmt.Println("WebSocket server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTerminal(w http.ResponseWriter, r *http.Request) {
	// Upgrade to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer ws.Close()

	// Get connection parameters from query string
	hostname := r.URL.Query().Get("hostname")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	port := 22
	if p := r.URL.Query().Get("port"); p != "" {
		fmt.Sscanf(p, "%d", &port)
	}
	protocol := r.URL.Query().Get("protocol")
	if protocol == "" {
		protocol = "ssh"
	}

	// Create tunnel
	config := client.DefaultConfig()
	tunnel := client.NewTunnel(config)

	// Connect to target server
	if err := tunnel.Connect(hostname, username, password, port, protocol, 1024, 768); err != nil {
		log.Printf("Failed to connect: %v", err)
		ws.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	defer tunnel.Close()

	// Use errgroup to manage bidirectional data forwarding
	eg, ctx := errgroup.WithContext(context.Background())

	// Forward from tunnel to WebSocket
	eg.Go(func() error {
		buf := bytes.NewBuffer(make([]byte, 0, 8192))
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				data, err := tunnel.ReadSome()
				if err != nil {
					if err == io.EOF {
						return nil
					}
					return err
				}
				if len(data) > 0 {
					buf.Write(data)
					// Ensure complete instructions
					lastComplete := client.FindLastCompleteInstruction(buf.Bytes())
					if lastComplete > 0 {
						toSend := buf.Bytes()[:lastComplete]
						if err := ws.WriteMessage(websocket.TextMessage, toSend); err != nil {
							return err
						}
						remaining := buf.Bytes()[lastComplete:]
						buf.Reset()
						buf.Write(remaining)
					}
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	})

	// Forward from WebSocket to tunnel
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				_, message, err := ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						return err
					}
					return nil
				}
				if _, err := tunnel.Write(message); err != nil {
					return err
				}
			}
		}
	})

	// Wait for either goroutine to finish
	if err := eg.Wait(); err != nil {
		log.Printf("Forwarding error: %v", err)
	}
}
