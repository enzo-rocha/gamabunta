package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error in [conn.ReadMessage]: %v", err)
			return
		}

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Error in [WriteMessage]: %v", err)
			return
		}

	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error in [upgrader.Upgrade]: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("===================================")
	log.Println("== CLIENT SUCCESSFULLY CONNECTED ==")
	log.Println("===================================")

	reader(conn)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", ws)
}

func main() {
	setupRoutes()

	log.Println("=======================")
	log.Println("== RUNNING SERVER... ==")
	log.Println("=======================")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error in [ListenAndServe]: %v", err)
		return
	}
}
