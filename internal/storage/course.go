package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type Course struct {
	db *gorm.DB
}

func (c *Course) SaveCourse(course *internal.Course) error {
	return c.db.Save(course).Error
}

func (c *Course) GetByID(id uint) (*internal.Course, error) {
	c.db.Lock()
	defer c.db.Unlock()

	courses := make([]internal.Course, 0)
	c.db.Where(int64(id)).Find(&courses)

	if len(courses) == 0 {
		return nil, fmt.Errorf("course with id: %d not exit", id)
	}

	return &courses[0], nil
}

func (c *Course) GetManyByID(ids []uint) ([]internal.Course, error) {
	c.db.Lock()
	defer c.db.Unlock()

	courses := make([]internal.Course, 0)
	c.db.Where(ids).Find(&courses)

	return courses, nil
}
