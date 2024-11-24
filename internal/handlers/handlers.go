package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/w212w/GoProjectEM/internal/logger"
	"github.com/w212w/GoProjectEM/internal/models"
	"gorm.io/gorm"
)

// GetSongsHandler godoc
// @Summary Получить список песен
// @Description Получить список песен с возможностью фильтрации по артисту и названию
// @Tags songs
// @Accept json
// @Produce json
// @Param artist query string false "Фильтр по артисту"
// @Param title query string false "Фильтр по названию"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество результатов на странице" default(10)
// @Success 200 {array} models.Song "Список песен"
// @Failure 400 {object} models.ErrorResponse "Неверные параметры"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /songs [get]
func GetSongsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("GetSongsHandler: Start processing request")

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

		logger.Log.Debugf("GetSongsHandler: Parameters received - artist: %s, title: %s, page: %d, limit: %d", artist, title, page, limit)

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
			logger.Log.Error("GetSongsHandler: Failed to retrieve songs")
			http.Error(w, "Failed to retrieve songs", http.StatusInternalServerError)
			return
		}

		logger.Log.Debug("GetSongsHandler: Songs retrieved successfully")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(songs); err != nil {
			logger.Log.Error("GetSongsHandler: Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("GetSongsHandler: Successfully responded with songs")
	}
}

// GetSongTextHandler godoc
// @Summary Получить текст песни
// @Description Получить текст песни по ее ID с возможностью пагинации по стихам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество стихов на странице" default(2)
// @Success 200 {object} models.SongTextResponse "Текст песни с пагинацией"
// @Failure 400 {object} models.ErrorResponse "Неверные параметры"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /songs/{id}/text [get]
func GetSongTextHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("GetSongTextHandler: Start processing request")

		vars := mux.Vars(r)
		id := vars["id"]

		logger.Log.Debugf("GetSongTextHandler: Song ID received: %s", id)

		var song models.Song
		if err := db.First(&song, id).Error; err != nil {
			logger.Log.Error("GetSongTextHandler: Song not found")
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
			logger.Log.Error("GetSongTextHandler: Page out of range")
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

		logger.Log.Debug("GetSongTextHandler: Response prepared successfully")

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Log.Error("GetSongTextHandler: Failed to encode response")
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("GetSongTextHandler: Successfully responded with song text")
	}
}

// DeleteSongHandler godoc
// @Summary Удалить песню
// @Description Удалить песню по ее ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Success 200 {string} string "Песня удалена успешно"
// @Failure 400 {object} models.ErrorResponse "ID не указан"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /songs/{id} [delete]
func DeleteSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("DeleteSongHandler: Start processing request")
		vars := mux.Vars(r)
		id := vars["id"]

		logger.Log.Debugf("DeleteSongHandler: Song ID received: %s", id)

		if id == "" {
			logger.Log.Error("DeleteSongHandler: ID is required")
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var song models.Song
		if err := db.First(&song, "id = ?", id).Error; err != nil {
			logger.Log.Error("DeleteSongHandler: Song not found")
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Song not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to find song", http.StatusInternalServerError)
			}
			return
		}

		if err := db.Delete(&song).Error; err != nil {
			logger.Log.Error("DeleteSongHandler: Failed to delete song")
			http.Error(w, "Failed to delete song", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("DeleteSongHandler: Song deleted successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Song deleted successfully"))
	}
}

// UpdateSongHandler godoc
// @Summary Обновить информацию о песне
// @Description Обновить песню по ее ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "ID песни"
// @Param song body models.Song true "Данные для обновления песни"
// @Success 200 {string} string "Песня обновлена успешно"
// @Failure 400 {object} models.ErrorResponse "Неверный формат JSON"
// @Failure 404 {object} models.ErrorResponse "Песня не найдена"
// @Failure 500 {object} models.ErrorResponse "Ошибка сервера"
// @Router /songs/{id} [put]
func UpdateSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debug("UpdateSongHandler: Start processing request")

		vars := mux.Vars(r)
		id := vars["id"]

		logger.Log.Debugf("UpdateSongHandler: Song ID received: %s", id)

		if id == "" {
			logger.Log.Error("UpdateSongHandler: ID is required")
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		var song models.Song
		if err := db.First(&song, "id = ?", id).Error; err != nil {
			logger.Log.Error("UpdateSongHandler: Song not found")
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "Song not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to find song", http.StatusInternalServerError)
			}
			return
		}

		var updatedData models.Song
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			logger.Log.Error("UpdateSongHandler: Invalid JSON format")
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		song.Artist = updatedData.Artist
		song.Title = updatedData.Title
		song.ReleaseDate = updatedData.ReleaseDate
		song.Text = updatedData.Text
		song.Link = updatedData.Link
		song.Group = updatedData.Group

		if err := db.Save(&song).Error; err != nil {
			logger.Log.Error("UpdateSongHandler: Failed to update song")
			http.Error(w, "Failed to update song", http.StatusInternalServerError)
			return
		}

		logger.Log.Info("UpdateSongHandler: Song updated successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Song updated successfully"))
	}
}

// AddSongHandler godoc
// @Summary Добавить песню
// @Description Добавляет песню в базу данных, получая информацию о песне из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.AddSongRequest true "Данные для добавления песни"
// @Success 201 {string} string "Песня успешно добавлена"
// @Failure 400 {object} models.ErrorResponse "Неверный формат данных"
// @Failure 500 {object} models.ErrorResponse "Ошибка при обработке запроса"
// @Router /songs [post]
func AddSongHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Infof("Received request to add song from %s", r.RemoteAddr)

		type Input struct {
			Group string `json:"group"`
			Song  string `json:"song"`
		}

		var input Input
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			logger.Log.Errorf("Invalid input: %v", err)
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		logger.Log.Infof("Parsed input: group=%s, song=%s", input.Group, input.Song)

		baseURL := os.Getenv("EXTERNAL_API_BASE_URL")
		if baseURL == "" {
			logger.Log.Error("External API base URL not configured")
			http.Error(w, "External API base URL not configured", http.StatusInternalServerError)
			return
		}
		url := fmt.Sprintf("%s/info?group=%s&song=%s", baseURL, input.Group, input.Song)
		logger.Log.Infof("Making request to external API: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Errorf("Failed to fetch song info: %v", err)
			http.Error(w, "Failed to fetch song info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logger.Log.Errorf("External API returned an error: %s", resp.Status)
			http.Error(w, "External API returned an error", http.StatusInternalServerError)
			return
		}

		var apiResponse struct {
			Artist      string `json:"artist"`
			ReleaseDate string `json:"releaseDate"`
			Text        string `json:"text"`
			Link        string `json:"link"`
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Errorf("Failed to read API response: %v", err)
			http.Error(w, "Failed to read API response", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(body, &apiResponse); err != nil {
			logger.Log.Errorf("Invalid API response format: %v", err)
			http.Error(w, "Invalid API response format", http.StatusInternalServerError)
			return
		}

		newSong := models.Song{
			Group:       input.Group,
			Title:       input.Song,
			Artist:      apiResponse.Artist,
			ReleaseDate: apiResponse.ReleaseDate,
			Text:        apiResponse.Text,
			Link:        apiResponse.Link,
		}

		if err := db.Create(&newSong).Error; err != nil {
			logger.Log.Errorf("Failed to save song to database: %v", err)
			http.Error(w, "Failed to save song", http.StatusInternalServerError)
			return
		}

		logger.Log.Infof("Song added successfully: %s by %s", input.Song, input.Group)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Song added successfully"))
	}
}
