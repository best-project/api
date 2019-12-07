package internal

type CourseDTO struct {
	CourseID string `json:"id"`
	UserID   uint   `json:"userId" validate:"required"`

	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Data        []TaskDTO `json:"data,omitempty"`

	Language        string  `json:"language"`
	DifficultyLevel string  `json:"difficultyLevel"`
	Rate            float32 `json:"rate"`
	MaxPoints       int     `json:"maxPoints"`

	Type string `json:"type" validate:"required"`
}

type CourseResultDTO struct {
	UserID   uint   `json:"userId"`
	CourseID string `json:"courseId" validate:"required"`

	Phase  string `json:"phase" validate:"required"`
	Points uint   `json:"points"`
	Passed bool   `json:"passed"`
}

type TaskDTO struct {
	ID        string `json:"id"`
	CourseID  string `json:"courseId"`
	Word      string `json:"word" validate:"required"`
	Translate string `json:"translate" validate:"required"`
	Image     string `json:"image"`
}

type UserDTO struct {
	Email    string `json:"email" validate:"required,email,max=250"`
	Username string `json:"username" validate:"required,max=250"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=250"`
	Avatar   string `json:"avatar"`

	Token string `json:"token"`

	ID int `json:"id"`

	Level  int `json:"level"`
	Points int `json:"points"`
}

type UserStatDTO struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Points   int    `json:"points"`
	Level    int    `json:"level"`
}
