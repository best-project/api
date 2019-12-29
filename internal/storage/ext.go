package storage

import "github.com/best-project/api/internal"

type Course interface {
	SaveCourse(course *internal.Course, xpForTask int) error
	GetByUserID(id uint) ([]*internal.Course, error)
	GetAll() ([]*internal.Course, error)
	GetByID(id string) (*internal.Course, error)
	Exist(courseID string) bool
	GetManyByID(ids []string) ([]*internal.Course, error)
}
type CourseResult interface {
	SaveResult(course *internal.CourseResult) error
	ReplaceIfExist(result *internal.CourseResult) (*internal.CourseResult, uint, error)
	ListBestResultsForUser(userID uint) ([]internal.CourseResult, error)
	GetBestResultForUserForCourse(userID uint, courseID string) (*internal.CourseResult, error)
	ListAllForUser(userID uint) ([]internal.CourseResult, error)
	ListStartedForUser(userID uint) ([]internal.CourseResult, error)
	ListFinishedForUser(userID uint) ([]internal.CourseResult, error)
	ListResultsForCourse(courseID string) ([]internal.CourseResult, error)
}
type Task interface {
	SaveTask(task *internal.Task) error
	GetTasksForCourse(course *internal.Course) []internal.Task
	MapTasksForCourses(courses []*internal.Course) map[string][]internal.Task
	GetByID(id uint) (*internal.Task, error)
	GetManyByID(ids []uint) ([]internal.Task, error)
	RemoveByID(task *internal.Task) error
}
type User interface {
	SaveUser(user *internal.User) error
	Exist(user *internal.User) bool
	GetByID(id uint) (*internal.User, error)
	GetManyByID(ids []uint) ([]internal.User, error)
	GetByMail(email string) ([]internal.User, error)
	GetAll() ([]internal.User, error)
}
