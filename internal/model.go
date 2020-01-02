package internal

import (
	"time"
)

type Course struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

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
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

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
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	CourseID  uint
	Word      string
	Translate string
	Image     string
}

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Email     string
	FirstName string
	LastName  string
	Password  string
	Avatar    string

	Level     int
	NextLevel string
	Points    int
}
