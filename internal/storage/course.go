package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type Course struct {
	db *gorm.DB
}

func (c *Course) SaveCourse(course *internal.Course) {
	c.db.Lock()
	defer c.db.Unlock()
	c.db.Save(course)
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
