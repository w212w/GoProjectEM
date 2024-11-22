package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func GetSongTextHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteSongHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateSongHandler(w http.ResponseWriter, r *http.Request) {

}

func AddSongHandler(w http.ResponseWriter, r *http.Request) {

}
