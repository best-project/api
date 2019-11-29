package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Course struct {
	db *gorm.DB
}

func (c *Course) SaveCourse(course *internal.Course) error {
	return c.db.Save(course).Error
}

func (c *Course) GetByUserID(id uint) ([]internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]internal.Course, 0)

	if err := c.db.Where(&internal.Course{UserID: id}).Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses for user id %d", id)
	}

	return courses, nil
}

func (c *Course) GetAll() ([]internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]internal.Course, 0)

	if err := c.db.Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses")
	}

	return courses, nil
}

func (c *Course) GetByID(id uint) (*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]internal.Course, 0)
	c.db.Where(int64(id)).Find(&courses)

	if len(courses) == 0 {
		return nil, fmt.Errorf("course with id: %d not exit", id)
	}

	return &courses[0], nil
}

func (c *Course) GetManyByID(ids []uint) ([]internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]internal.Course, 0)
	c.db.Where(ids).Find(&courses)

	return courses, nil
}
