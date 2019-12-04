package internal

import (
	"github.com/jinzhu/gorm"
)

type Course struct {
	gorm.Model

	CourseID        string
	UserID          uint
	Name            string
	Description     string
	Image           string
	Task            []Task
	Language        string
	DifficultyLevel string
	Rate            float32
	Type            string
	MaxPoints       int
}

const (
	ImageType  string = "image"
	PuzzleType string = "puzzle"
)

type CourseResult struct {
	gorm.Model

	UserID   uint
	CourseID string

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

	CourseID  string
	Word      string
	Translate string
	Image     string
}

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
