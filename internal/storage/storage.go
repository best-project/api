package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	User   *User
	Course *Course
	Task   *Task
}

func NewDatabase(cfg *config.Config, entry *logrus.Logger) (*Database, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)

	entry.Info("Starting database connection")
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	// for development
	tables := []interface{}{&internal.Course{}, &internal.User{}, &internal.Task{}}

	entry.Info("Clearing database")
	db.DropTableIfExists(tables...)
	db.CreateTable(tables...)

	userDB := &User{db}
	courseDB := &Course{db}
	taskDB := &Task{db}

	pass, err := bcrypt.GenerateFromPassword([]byte("root123"), bcrypt.MinCost)
	if err != nil {
		return nil, errors.Wrap(err, "while hashing pass")
	}

	userDB.SaveUser(&internal.User{Model: gorm.Model{ID: uint(1)}, Username: "root", Email: "root", Password: string(pass)})
	courseDB.SaveCourse(&internal.Course{Name: "XD", UserID: 1, Difficulty: "HARD"})
	courseDB.SaveCourse(&internal.Course{Name: "2", UserID: 1, MaxPoints: 213, Description: "D"})
	courseDB.SaveCourse(&internal.Course{Name: "3 XD", UserID: 1, Description: "a tu description"})

	return &Database{userDB, courseDB, taskDB}, nil
}
