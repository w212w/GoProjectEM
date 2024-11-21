package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/w212w/GoProjectEM/internal/handlers"
	"github.com/w212w/GoProjectEM/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env")
	}
}

func setupDatabase() *gorm.DB {
	loadEnv()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	return db
}

func main() {

	db := setupDatabase()
	fmt.Println("Connected to database")

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		log.Fatal("Migartion failed:", err)
	}
	fmt.Println("Table created successfully")

	router := mux.NewRouter()

	router.HandleFunc("/api/songs", handlers.GetSongsHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}/text", handlers.GetSongTextHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}", handlers.DeleteSongHandler).Methods("DELETE")
	router.HandleFunc("/api/songs/{id}", handlers.UpdateSongHandler).Methods("PUT")
	router.HandleFunc("/api/songs", handlers.AddSongHandler).Methods("POST")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
