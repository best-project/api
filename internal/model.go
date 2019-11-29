package internal

import (
	"github.com/jinzhu/gorm"
)

type Course struct {
	gorm.Model

	UserID      uint
	Name        string
	Description string
	Image       string
	Task        []Task
	Language    string
	Difficulty  string
	Rate        float32
	MaxPoints   int
}

type Task struct {
	gorm.Model

	Type      string
	Word      string
	Translate string
	Image     string
}

const (
	ImageType  string = "image"
	PuzzleType string = "puzzle"
)

type User struct {
	gorm.Model

	Email    string
	Username string
	Password []byte
	Avatar   string

	Token string `sql:"size:2000"`

	Level  int
	Points int
}
