package storage

import (
	"github.com/best-project/api/internal/model"
	"database/sql"
)

type UserStorage struct {
	db sql.DB
}

func (s *CourseStorage) getUser(id string) *model.User {
	return &model.User{}
}