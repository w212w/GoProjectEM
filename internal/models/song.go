package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Artist      string `json:"artist"`
	Title       string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type SongTextResponse struct {
	TotalVerses int      `json:"total_verses"`
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	Verses      []string `json:"verses"`
}
