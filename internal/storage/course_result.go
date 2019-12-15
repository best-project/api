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

func (c *CourseResultDB) ListResultsForCourse(courseID string) ([]internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{CourseID: courseID, Phase: internal.FinishedPhase}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
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

func (c *CourseResultDB) ReplaceIfExist(result *internal.CourseResult) (*internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: result.UserID, CourseID: result.CourseID, Phase: result.Phase}).Find(&results).Error
	if err != nil {
		return nil, err
	}
	if len(results) < 1 {
		return result, nil
	}
	if result.Phase == internal.FinishedPhase {
		if results[0].Points > result.Points {
			return &results[0], nil
		}
		result.Model = results[0].Model
		return result, nil
	}

	return &results[0], nil
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
