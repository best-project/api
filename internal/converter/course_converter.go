package converter

import (
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/pkg/errors"
)

type CourseConverter struct {
	TaskConverter *TaskConverter
	db            *storage.Database
}

func NewCourseConverter(db *storage.Database, taskConverter *TaskConverter) *CourseConverter {
	return &CourseConverter{
		TaskConverter: taskConverter,
		db:            db,
	}
}

func (c *CourseConverter) ToModel(dto *internal.CourseDTO) *internal.Course {
	return &internal.Course{
		Name:        dto.Name,
		Description: dto.Description,
		Difficulty:  dto.Difficulty,
		Image:       dto.Image,
		Language:    dto.Language,
		MaxPoints:   dto.MaxPoints,
		Rate:        dto.Rate,
	}
}

func (c *CourseConverter) ManyToModel(dtos []internal.CourseDTO) []internal.Course {
	result := make([]internal.Course, len(dtos))

	for _, dto := range dtos {
		result = append(result, *c.ToModel(&dto))
	}

	return result
}

func (c *CourseConverter) ToDTO(dto *internal.Course) (*internal.CourseDTO, error) {
	return &internal.CourseDTO{
		Name:        dto.Name,
		Rate:        dto.Rate,
		Image:       dto.Image,
		Language:    dto.Language,
		MaxPoints:   dto.MaxPoints,
		Difficulty:  dto.Difficulty,
		Description: dto.Description,
		Data:        c.TaskConverter.ManyToDTO(dto.Task),
		UserID:      dto.UserID,
	}, nil
}

func (c *CourseConverter) ManyToDTO(courses []internal.Course) ([]internal.CourseDTO, error) {
	result := make([]internal.CourseDTO, len(courses))

	for _, course := range courses {
		dto, err := c.ToDTO(&course)
		if err != nil {
			return nil, errors.Wrapf(err, "while converting course %d", course.ID)
		}
		result = append(result, *dto)
	}

	return result, nil
}
