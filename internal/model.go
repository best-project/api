package internal

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Course struct {
	gorm.Model

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

func NewUser(user *User) *User {
	now := time.Now()
	return &User{
		Username: user.Username,
		Password: user.Password,
		Model: gorm.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
		ProfileCourses: []Course{},
		MyCourses:      []Course{},
	}
}

type User struct {
	gorm.Model

	Username string
	Password []byte
	Avatar   string

	Token string `sql:"size:2000"`

	Level  int
	Points int

	ProfileCourses []Course
	MyCourses      []Course
}
