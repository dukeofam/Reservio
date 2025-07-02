package config

import (
	"log"
	"net/http"
	"os"

	"github.com/boj/redistore"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"reservio/models"
)

var DB *gorm.DB
var Store *redistore.RediStore

func init() {
	_ = godotenv.Load()
}

func GetEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key)
	}
	return val
}

var DBUri = GetEnvOrFatal("DATABASE_URL")

func ConnectDatabase() {
	dsn := DBUri
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(&models.User{}, &models.Child{}, &models.Reservation{}, &models.Slot{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	DB = database
}

func InitSessionStore() {
	var err error
	Store, err = redistore.NewRediStore(10, "tcp", "localhost:6379", "", os.Getenv("SESSION_SECRET"))
	if err != nil {
		log.Fatalf("Failed to connect to Redis session store: %v", err)
	}
	Store.SetMaxAge(3600) // seconds
	Store.Options.HttpOnly = true
	Store.Options.Secure = true
	Store.Options.SameSite = http.SameSiteStrictMode
}
