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
	Image           string `sql:"size:2000"`
	Task            []Task
	Language        string
	DifficultyLevel string
	Rate            float32
	RateCounter     int

	Type      string
	MaxPoints int
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

	Email     string
	FirstName string
	LastName  string
	Password  string
	Avatar    string

	Token        string `sql:"size:500"`
	RefreshToken string `sql:"size:500"`

	Level     int
	NextLevel int
	Points    int
}
