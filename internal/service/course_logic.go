package service

import (
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/pkg/errors"
	"strconv"
)

type CourseLogic struct {
	db          *storage.Database
	passPercent float32
}

var firstLevel = 100

type CourseResultHandler interface {
	CheckResult(result *internal.CourseResult) (bool, error)
	PointForNextLevel(points int) int
}

func NewCourseLogic(db *storage.Database, passPercent float32) *CourseLogic {
	return &CourseLogic{
		db:          db,
		passPercent: passPercent,
	}
}

func (c *CourseLogic) CheckResult(result *internal.CourseResult) (bool, error) {
	course, err := c.db.Course.GetByID(strconv.Itoa(int(result.CourseID)))
	if err != nil {
		return false, errors.Wrap(err, "while getting course")
	}
	if !c.isWon(result.Points, uint(course.MaxPoints)) {
		return false, nil
	}
	result.Passed = true

	if err := c.db.CourseResult.SaveResult(result); err != nil {
		return true, err
	}

	user, err := c.db.User.GetByID(result.UserID)
	if err != nil {
		return true, err
	}
	user.Points += int(result.Points)
	user.Level, user.NextLevel = c.calculateLevel(user.Points)
	if err = c.db.User.SaveUser(user); err != nil {
		return true, err
	}

	return true, nil
}

func (c *CourseLogic) PointForNextLevel(points int) int {
	_, nextLvl := c.calculateLevel(points)
	return nextLvl
}

func (c *CourseLogic) isWon(points, maxPoints uint) bool {
	result := float32(points) / float32(maxPoints)
	return result > c.passPercent
}

func (c *CourseLogic) calculateLevel(points int) (int, int) {
	nextLvl := 0
	for i := 1; true; i++ {
		xpNeeded := int(float64(firstLevel+10*i) * 1.2)
		nextLvl += xpNeeded
		if nextLvl > points {
			return i, xpNeeded
		}
	}
	return 1, nextLvl
}
