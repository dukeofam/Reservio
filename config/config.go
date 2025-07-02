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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	storeKey := os.Getenv("SESSION_SECRET")
	if storeKey == "" {
		log.Fatal("SESSION_SECRET must be set")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	var store *redistore.RediStore

	if redisPassword == "" {
		store, err = redistore.NewRediStoreWithDB(10, "tcp", redisAddr, "", "0", storeKey)
	} else {
		store, err = redistore.NewRediStoreWithDB(10, "tcp", redisAddr, redisPassword, "0", storeKey)
	}

	if err != nil {
		log.Fatalf("Failed to connect to Redis session store: %v", err)
	}

	store.SetMaxAge(3600)
	store.Options.HttpOnly = true
	store.Options.Secure = true
	store.Options.SameSite = http.SameSiteStrictMode

	Store = store
}
