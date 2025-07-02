package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"reservio/config"
	"reservio/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDatabase()
	config.InitSessionStore()

	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
