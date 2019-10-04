package storage

import (
	"database/sql"
)

type CourseStorage struct {
	db sql.DB
}

