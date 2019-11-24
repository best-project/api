package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Database struct {
	User   *User
	Course *Course
	Task   *Task
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)

	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	// for development
	tables := []interface{}{&internal.Course{}, &internal.User{}, &internal.Task{}}

	db.DropTableIfExists(tables...)
	db.CreateTable(tables...)

	return &Database{&User{db}, &Course{db}, &Task{db}}, nil
}
