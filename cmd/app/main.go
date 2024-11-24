package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/w212w/GoProjectEM/internal/handlers"
	"github.com/w212w/GoProjectEM/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Log = logrus.New()

func setupLogger() {

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		Log.Fatalf("Неверный уровень логирования: %v", err)
	}
	Log.SetLevel(level)

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.Debug("Logrus успешно инициализирован")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		Log.Fatal("Ошибка загрузки .env файла")
	} else {
		Log.Debug("Файл .env успешно загружен")
	}
}

func setupDatabase() *gorm.DB {

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
		Log.Fatal("Failed to connect to the database:", err)
	}

	return db
}

func main() {

	loadEnv()
	setupLogger()
	db := setupDatabase()
	Log.Info("Connected to database")

	if err := db.AutoMigrate(&models.Song{}); err != nil {
		Log.Fatal("Migartion failed:", err)
	}
	Log.Debug("Table created successfully")

	router := mux.NewRouter()

	router.HandleFunc("/api/songs", handlers.GetSongsHandler(db)).Methods("GET")
	router.HandleFunc("/api/songs/{id}/text", handlers.GetSongTextHandler(db)).Methods("GET")
	router.HandleFunc("/api/songs/{id}", handlers.DeleteSongHandler(db)).Methods("DELETE")
	router.HandleFunc("/api/songs/{id}", handlers.UpdateSongHandler(db)).Methods("PUT")
	router.HandleFunc("/api/songs", handlers.AddSongHandler(db)).Methods("POST")

	Log.Info("Server is running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		Log.Fatalf("Error running server: %v", err)
	}
}
