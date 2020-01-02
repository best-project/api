package storage

import (
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"sort"
	"strconv"
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
	id, _ := strconv.Atoi(courseID)
	err := c.db.Where(internal.CourseResult{CourseID: uint(id), Phase: internal.FinishedPhase}).Find(&results).Error
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

func (c *CourseResultDB) ListBestResultsForUser(userID uint) ([]internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: userID, Phase: internal.FinishedPhase}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	topResults := make(map[string]internal.CourseResult, 0)
	for _, result := range results {
		uid := strconv.Itoa(int(result.CourseID))
		if topResult, ok := topResults[uid]; !ok {
			topResults[uid] = result
		} else {
			if result.Points > topResult.Points {
				topResults[uid] = result
			}
		}
	}

	results = make([]internal.CourseResult, 0)
	for _, result := range topResults {
		results = append(results, result)
	}

	return results, nil
}

func (c *CourseResultDB) GetBestResultForUserForCourse(userID uint, courseID string) (*internal.CourseResult, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	id, _ := strconv.Atoi(courseID)

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: userID, CourseID: uint(id), Phase: internal.FinishedPhase}).Find(&results).Error
	if err != nil {
		return nil, err
	}

	mem := 0
	r := &internal.CourseResult{}
	for i, result := range results {
		if r.Points < result.Points {
			r.Points = result.Points
			mem = i
		}
	}

	return &results[mem], nil
}

func (c *CourseResultDB) ReplaceIfExist(result *internal.CourseResult) (*internal.CourseResult, uint, error) {
	c.db.RLock()
	defer c.db.RUnlock()

	results := make([]internal.CourseResult, 0)
	err := c.db.Where(internal.CourseResult{UserID: result.UserID, CourseID: result.CourseID, Phase: internal.FinishedPhase}).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}
	if len(results) < 1 {
		return result, result.Points, nil
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Points < results[j].Points
	})
	if results[0].Points > result.Points {
		return &results[0], results[0].Points, nil
	}

	return result, result.Points, nil
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
