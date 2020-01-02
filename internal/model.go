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
	RateCounter     int

	Type      string
	MaxPoints int
}

const (
	NormalType string = "normal"
	PuzzleType string = "puzzle"
)

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

type User struct {
	gorm.Model

	Email     string
	FirstName string
	LastName  string
	Password  string
	Avatar    string

	Level     int
	NextLevel int
	Points    int
}
