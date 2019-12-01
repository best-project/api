package internal

import (
	"github.com/jinzhu/gorm"
)

type Course struct {
	gorm.Model

	UserID          uint
	Name            string
	Description     string
	Image           string
	Task            []Task
	Language        string
	DifficultyLevel string
	Rate            float32
	MaxPoints       int
}

type CourseResult struct {
	gorm.Model

	UserID   uint
	CourseID uint

	Points uint
	Phase  string
	Passed bool
}

const (
	StartedPhase  string = "started"
	FinishedPhase string = "finished"
)

type Task struct {
	gorm.Model

	CourseID  uint
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
	Password string
	Avatar   string

	Token string `sql:"size:2000"`

	Level  int
	Points int
}
