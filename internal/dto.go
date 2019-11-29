package internal

type CourseDTO struct {
	UserID      uint      `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Data        []TaskDTO `json:"data"`
	Language    string    `json:"language"`
	Difficulty  string    `json:"difficulty"`
	Rate        float32   `json:"rate"`
	MaxPoints   int       `json:"maxPoints"`
}

type TaskDTO struct {
	Type      string `json:"type"`
	Word      string `json:"word"`
	Translate string `json:"translate"`
	Image     string `json:"image"`
}

type UserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar"`

	Token string `json:"token"`

	ID int `json:"id"`

	Level  int `json:"level"`
	Points int `json:"points"`

	ProfileCourses []CourseDTO `json:"profileCourses"`
}
