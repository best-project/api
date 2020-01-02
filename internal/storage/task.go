package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type TaskDB struct {
	db *gorm.DB
}

func (c *TaskDB) SaveTask(task *internal.Task) error {
	return c.db.Save(task).Error
}

func (c *TaskDB) GetTasksForCourse(course *internal.Course) []internal.Task {
	c.db.RLock()
	defer c.db.RUnlock()

	tasks := make([]internal.Task, 0)
	c.db.Where(&internal.Task{CourseID: course.ID}).Find(&tasks)

	return tasks
}

func (c *TaskDB) MapTasksForCourses(courses []*internal.Course) map[uint][]internal.Task {
	c.db.RLock()
	defer c.db.RUnlock()

	courseTasks := make(map[uint][]internal.Task, 0)
	tasks := make([]internal.Task, 0)

	for _, course := range courses {
		c.db.Where(&internal.Task{CourseID: course.ID}).Find(&tasks)
		courseTasks[course.ID] = tasks
	}

	return courseTasks
}

func (c *TaskDB) GetByID(id uint) (*internal.Task, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	tasks := make([]internal.Task, 0)
	c.db.Where(int64(id)).Find(&tasks)

	if len(tasks) == 0 {
		return nil, fmt.Errorf("task with id: %d not exit", id)
	}

	return &tasks[0], nil
}

func (c *TaskDB) GetManyByID(ids []uint) ([]internal.Task, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	tasks := make([]internal.Task, 0)
	c.db.Where(ids).Find(&tasks)

	return tasks, nil
}

func (c *TaskDB) RemoveByID(task *internal.Task) error {
	c.db.RLock()
	defer c.db.RUnlock()

	if task.ID == 0 {
		return errors.New("cannot remove task with id 0")
	}

	return c.db.Delete(task).Error
}
