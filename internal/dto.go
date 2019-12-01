package internal

type CourseDTO struct {
	UserID uint `json:"userId"`

	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Data        []TaskDTO `json:"data"`

	InProgress      bool    `json:"inProgress"`
	Language        string  `json:"language"`
	DifficultyLevel string  `json:"difficultyLevel"`
	Rate            float32 `json:"rate"`
	MaxPoints       int     `json:"maxPoints"`

	Type string `json:"type"`
}

type CourseResultDTO struct {
	UserID   uint `json:"userId" validate:"required"`
	CourseID uint `json:"courseId" validate:"required"`

	Phase  string `json:"phase" validate:"required"`
	Points uint   `json:"points"`
	Passed bool   `json:"passed"`
}

type TaskDTO struct {
	Word      string `json:"word" validate:"required"`
	Translate string `json:"translate" validate:"required"`
	Image     string `json:"image"`
}

type UserDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Avatar   string `json:"avatar"`

	Token string `json:"token"`

	ID int `json:"id"`

	Level  int `json:"level"`
	Points int `json:"points"`

	ProfileCourses []CourseDTO `json:"profileCourses"`
}
