package service

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/pkg/errors"
)

type CourseLogic struct {
	db *storage.Database

	passPercent float32
}

var firstLevel = 100

func NewCourseLogic(db *storage.Database, passPercent float32) *CourseLogic {
	return &CourseLogic{
		db:          db,
		passPercent: passPercent,
	}
}

func (c *CourseLogic) CheckResult(result *internal.CourseResult) (bool, error) {
	course, err := c.db.Course.GetByID(result.CourseID)
	if err != nil {
		return false, errors.Wrap(err, "while getting course")
	}
	if !c.isWon(result.Points, uint(course.MaxPoints)) {
		return false, nil
	}
	result.Passed = true

	fmt.Println(*result)
	if err := c.db.CourseResult.SaveResult(result); err != nil {
		return true, err
	}

	user, err := c.db.User.GetByID(result.UserID)
	if err != nil {
		return true, err
	}
	user.Points += int(result.Points)
	user.Level = c.calculateLevel(user.Points)
	if err = c.db.User.SaveUser(user); err != nil {
		return true, err
	}

	return true, nil

}

func (c *CourseLogic) isWon(points, maxPoints uint) bool {
	result := float32(points) / float32(maxPoints)
	return result > c.passPercent
}

func (c *CourseLogic) calculateLevel(points int) int {
	xp := 0
	for i := 0; true; i++ {
		xp += int(float64(firstLevel+10*i) * float64(1.2))
		if xp > points {
			return i
		}
	}
	return 1
}
