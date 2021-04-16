package main

import (
	"log"
	"net/http"

	"os"

	"github.com/ogustavobelo/go-study-chat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener...")
	go handlers.ListenToWebSocketChannel()

	log.Println("Webserver is starting on port 8080...")
	port := os.Getenv("PORT")
	parsedPort := ":" + port
	_ = http.ListenAndServe(parsedPort, mux)
}
