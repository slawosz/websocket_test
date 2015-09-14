package main

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		msg := append([]byte("MSG RCVD: "), p...)
		if err = conn.WriteMessage(messageType, msg); err != nil {
			return err
		}
	}
}

func main() {
	http.HandleFunc("/bar", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
