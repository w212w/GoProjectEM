package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/w212w/GoProjectEM/internal/models"
	"gorm.io/gorm"
)

func GetSongsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		artist := r.URL.Query().Get("artist")
		title := r.URL.Query().Get("title")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page := 1
		limit := 10

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		var songs []models.Song
		query := db.Model(&models.Song{})

		if artist != "" {
			query = query.Where("artist ILIKE ?", "%"+artist+"%")
		}
		if title != "" {
			query = query.Where("title ILIKE ?", "%"+title+"%")
		}

		offset := (page - 1) * limit
		if err := query.Offset(offset).Limit(limit).Find(&songs).Error; err != nil {
			http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(songs); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func GetSongTextHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]

		var song models.Song
		if err := db.First(&song, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Song not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to retrieve song", http.StatusInternalServerError)
			}
			return
		}

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			limit = 2
		}

		verses := strings.Split(song.Text, "\n\n")
		totalVerses := len(verses)

		start := (page - 1) * limit
		end := start + limit
		if start > totalVerses {
			http.Error(w, "Page out of range", http.StatusBadRequest)
			return
		}
		if end > totalVerses {
			end = totalVerses
		}

		response := models.SongTextResponse{
			TotalVerses: totalVerses,
			Page:        page,
			Limit:       limit,
			Verses:      verses[start:end],
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var song models.Song
		if err := db.First(&song, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Song not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to find song", http.StatusInternalServerError)
			}
			return
		}

		if err := db.Delete(&song).Error; err != nil {
			http.Error(w, "Failed to delete song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Song deleted successfully"))
	}
}

func UpdateSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var song models.Song
		if err := db.First(&song, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Song not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to find song", http.StatusInternalServerError)
			}
			return
		}

		var updatedData models.Song
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		song.Artist = updatedData.Artist
		song.Title = updatedData.Title
		song.ReleaseDate = updatedData.ReleaseDate
		song.Text = updatedData.Text
		song.Link = updatedData.Link

		if err := db.Save(&song).Error; err != nil {
			http.Error(w, "Failed to update song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Song updated successfully"))
	}
}

func AddSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newSong models.Song
		if err := json.NewDecoder(r.Body).Decode(&newSong); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if newSong.Artist == "" || newSong.Title == "" {
			http.Error(w, "Artist and Title are required", http.StatusBadRequest)
			return
		}

		if err := db.Create(&newSong).Error; err != nil {
			http.Error(w, "Failed to add song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Song added successfully"))
	}
}
