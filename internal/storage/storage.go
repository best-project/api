package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/best-project/api/internal/config"
	"github.com/pkg/errors"
	"fmt"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s/%s", cfg.DbName, cfg.DbPass, cfg.DbHost))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}
	return &Database{db}, nil
}