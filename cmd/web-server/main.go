package main

import (
	"compile-server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleConnection)
	http.HandleFunc("/code", handlers.CodeHandler)
	log.Println("Server start on port 1234")
	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		return
	}

}
