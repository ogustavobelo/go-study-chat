package main

import (
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/joho/godotenv"
	"github.com/ogustavobelo/go-study-chat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener...")
	go handlers.ListenToWebSocketChannel()
	log.Println("Webserver is starting on port 8080...")
	port := os.Getenv("PORT")
	parsedPort := ":" + port
	fmt.Println(parsedPort)
	_ = http.ListenAndServe(parsedPort, mux)
}

func init() {
	//Check ENV variables.
	envChecks()
}

func envChecks() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error to load .env")
	}
	port, portExist := os.LookupEnv("PORT")

	if !portExist || port == "" {
		log.Fatal("PORT must be set in .env and not empty")
	}
}
