package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	Group       string `json:"group"`
	SongName    string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}
