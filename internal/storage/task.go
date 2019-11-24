package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type Task struct {
	db *gorm.DB
}

func (c *Task) SaveTask(task *internal.Task) error {
	return c.db.Save(task).Error
}

func (c *Task) GetByID(id uint) (*internal.Task, error) {
	c.db.Lock()
	defer c.db.Unlock()

	tasks := make([]internal.Task, 0)
	c.db.Where(int64(id)).Find(&tasks)

	if len(tasks) == 0 {
		return nil, fmt.Errorf("task with id: %d not exit", id)
	}

	return &tasks[0], nil
}

func (c *Task) GetManyByID(ids []uint) ([]internal.Task, error) {
	c.db.Lock()
	defer c.db.Unlock()

	tasks := make([]internal.Task, 0)
	c.db.Where(ids).Find(&tasks)

	return tasks, nil
}
