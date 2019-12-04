package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

type CourseDB struct {
	db *gorm.DB
}

func (c *CourseDB) SaveCourse(course *internal.Course) error {
	guid := xid.New()
	uid := guid.String()
	course.CourseID = uid
	tasks := course.Task
	course.Task = []internal.Task{}

	c.db.Save(course)
	for _, task := range tasks {
		task.CourseID = uid
		c.db.Save(&task)
	}
	return nil
}

func (c *CourseDB) GetByUserID(id uint) ([]*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]*internal.Course, 0)

	if err := c.db.Where(&internal.Course{UserID: id}).Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses for user id %d", id)
	}

	return courses, nil
}

func (c *CourseDB) GetAll() ([]*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]*internal.Course, 0)

	if err := c.db.Find(&courses).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses")
	}

	return courses, nil
}

func (c *CourseDB) GetByID(id string) (*internal.Course, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	courses := make([]internal.Course, 0)
	c.db.Where(&internal.Course{CourseID: id}).Find(&courses)

	if len(courses) == 0 {
		return nil, fmt.Errorf("course with id: %s not exit", id)
	}

	return &courses[0], nil
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

	c.db.Where(ids).Find(&courses)

	return courses, nil
}
