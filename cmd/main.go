package main

import (
	"log"
	"net/http"
	"os"

	"reservio/config"
	"reservio/routes"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	config.ConnectDatabase()
	config.InitSessionStore()

	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("zap sync error: %v", err)
		}
	}()

	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}
