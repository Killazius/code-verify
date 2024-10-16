package main

import (
	"compile-server/internal/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/code", handlers.CodeHandler)
	log.Println("Server start on port 1234")
	http.ListenAndServe(":1234", nil)

}
