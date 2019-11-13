package internal

import "github.com/jinzhu/gorm"

type Course struct {
	gorm.Model

	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Data        []Task  `json:"data"`
	Language    string  `json:"language"`
	Difficulty  string  `json:"difficulty"`
	Rate        float32 `json:"rate"`
	MaxPoints   int     `json:"maxPoints"`
}

type Task struct {
	gorm.Model

	Type      string `json:"type"`
	Word      string `json:"word"`
	Translate string `json:"translate"`
	Image     string `json:"image"`
}

const (
	ImageType  string = "image"
	PuzzleType string = "puzzle"
)

type User struct {
	gorm.Model

	Username string `json:"username"`
	Password []byte `json:"password,omitempty"`
	Avatar   string `json:"avatar"`

	Level  int `json:"level"`
	Points int `json:"points"`

	ProfileCourses []Course `json:"profileCourses"`
}
