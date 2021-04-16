package main

import (
	"log"
	"net/http"

	"github.com/ogustavobelo/go-study-chat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener...")
	go handlers.ListenToWebSocketChannel()

	log.Println("Webserver is starting on port 8080...")
	_ = http.ListenAndServe(":8080", mux)
}
