package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Initialize
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Create the websocket variable
var clients []websocket.Conn

func main() {
	// Create endpoint for websocket connection
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		// Initialize configs
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Error in [upgrader.Upgrade]: %v", err)
			return
		}

		clients = append(clients, *conn)

		// Loop if client send to server
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error in [ReadMessage]: %v", err)
				return
			}

			fmt.Printf("%s send: %s\n", conn.RemoteAddr(), string(msg))

			// Loop if message found and send again to client for
			// write in your browser
			for _, client := range clients {
				err = client.WriteMessage(msgType, msg)
				if err != nil {

					return
				}

			}
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		// w,r is write and delete
	})

	fmt.Println("Your server run 8182")
	err := http.ListenAndServe(":8182", nil)
	if err != nil {
		log.Fatalf("Error in [ListenAndServe]: %v", err)
		return
	}
}
