package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w212w/GoProjectEM/internal/handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/songs", handlers.GetSongsHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}/text", handlers.GetSongTextHandler).Methods("GET")
	router.HandleFunc("/api/songs/{id}", handlers.DeleteSongHandler).Methods("DELETE")
	router.HandleFunc("api/songs/{id}", handlers.UpdateSongHandler).Methods("PUT")
	router.HandleFunc("api/songs", handlers.AddSongHandler).Methods("POST")

	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
