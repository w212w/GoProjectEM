package models

import "gorm.io/gorm"

// Song модель для песни
// @Description Структура для описания песни
// @Properties:
//   id: integer "ID песни"
//   created_at: string "Дата создания"
//   updated_at: string "Дата последнего обновления"
//   artist: string "Имя артиста"
//   song: string "Название песни"
//   release_date: string "Дата релиза"
//   text: string "Текст песни"
//   link: string "Ссылка на песню"
//   group: string "Группа, к которой принадлежит песня"
type Song struct {
	gorm.Model
	Artist      string `json:"artist"`
	Title       string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	Group       string `json:"group"`
}

// SongTextResponse структура для ответа с текстом песни
// @Description Структура для ответа на запрос получения текста песни
// @Properties:
//   total_verses: int "Общее количество куплетов"
//   page: int "Номер страницы"
//   limit: int "Лимит результатов на странице"
//   verses: array "Массив строк с куплетами песни"
type SongTextResponse struct {
	TotalVerses int      `json:"total_verses"`
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	Verses      []string `json:"verses"`
}
