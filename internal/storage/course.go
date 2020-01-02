package storage

import (
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strconv"
)

type CourseDB struct {
	db *gorm.DB
}

func (c *CourseDB) SaveCourse(course *internal.Course, xpForTask int) error {
	course.MaxPoints = len(course.Task) * xpForTask

	return c.db.Save(course).Error
}

func (c *CourseDB) GetByUserID(id uint) ([]*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]*internal.Course, 0)

	if err := c.db.Preload("Task").Where(&internal.Course{UserID: id}).Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses for user id %d", id)
	}

	return courses, nil
}

func (c *CourseDB) GetAll() ([]*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]*internal.Course, 0)

	if err := c.db.Preload("Task").Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses")
	}

	return courses, nil
}

func (c *CourseDB) GetByID(id string) (*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	course := &internal.Course{}
	if err := c.db.Preload("Task").First(course, id).Error; err != nil {
		return nil, err
	}

	return course, nil
}

func (u *CourseDB) Exist(courseID string) bool {
	u.db.RLock()
	defer u.db.RUnlock()
	courses := make([]internal.Course, 0)
	id, _ := strconv.Atoi(courseID)

	u.db.Preload("Task").Where(id).Find(&courses)
	if len(courses) > 0 {
		return true
	}
	return false
}

func (c *CourseDB) GetManyByID(ids []string) ([]*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]*internal.Course, 0)
	for _, id := range ids {
		course, err := c.GetByID(id)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}
