package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const myMachineIP = "YOUR_IP_ADDRESS" // Replace with your actual IP address

func main() {
	http.HandleFunc("/ws", echo)
	http.HandleFunc("/", homePage)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error in [ListenAndServe]: %v", err)
		return
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error in [Upgrade]: %v", err)
		return
	}

	clients[conn] = true
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Error in [SplitHostPort]: %v", err)
		return
	}
	userName := "Guigas"
	if ip == myMachineIP {
		userName = "Enzo"
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			log.Printf("Error in [ReadMessage]: %v", err)
			return
		}

		fullMessage := fmt.Sprintf("%s: %s", userName, string(message))
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(fullMessage))
			if err != nil {
				log.Printf("Error in [WriteMessage]: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
