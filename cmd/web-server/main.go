package main

import (
	"compile-server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handlers.HandleConnection)
	log.Println("Server start on port 1235")
	err := http.ListenAndServe(":1235", nil)
	if err != nil {
		return
	}

}
