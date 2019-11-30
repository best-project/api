package storage

import (
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type CourseResultDB struct {
	db *gorm.DB
}

func (c *CourseResultDB) SaveResult(course *internal.CourseResult) error {
	return c.db.Save(course).Error
}

func (c *CourseResultDB) ListFinishedForUser(userID uint) ([]internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: userID, Phase: internal.FinishedPhase}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *CourseResultDB) ListStartedForUser(userID uint) ([]internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: userID, Phase: internal.StartedPhase}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *CourseResultDB) ListAllForUser(userID uint) ([]internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: userID}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
