package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	http.HandleFunc("/foobar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(".")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			p = append([]byte("SERVER GOT: "), p...)
			if err = conn.WriteMessage(messageType, p); err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	fmt.Println("Starting server")
	err := http.ListenAndServe("0.0.0.0:4000", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
