package converter

import (
	"github.com/best-project/api/internal"
	"github.com/pkg/errors"
)

type CourseConverter struct {
	TaskConverter *TaskConverter
}

func NewCourseConverter(taskConverter *TaskConverter) *CourseConverter {
	return &CourseConverter{
		TaskConverter: taskConverter,
	}
}

func (c *CourseConverter) ToModel(dto *internal.CourseDTO) *internal.Course {
	return &internal.Course{
		Name:            dto.Name,
		Description:     dto.Description,
		DifficultyLevel: dto.DifficultyLevel,
		Image:           dto.Image,
		Language:        dto.Language,
		MaxPoints:       dto.MaxPoints,
		Rate:            dto.Rate,
		Type:            dto.Type,
		UserID:          dto.UserID,
		Task:            c.TaskConverter.ManyToModel(dto.Data),
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
		CourseID:        dto.CourseID,
		Type:            dto.Type,
		Name:            dto.Name,
		Rate:            dto.Rate,
		Image:           dto.Image,
		Language:        dto.Language,
		MaxPoints:       dto.MaxPoints,
		DifficultyLevel: dto.DifficultyLevel,
		Description:     dto.Description,
		Data:            c.TaskConverter.ManyToDTO(dto.Task),
		UserID:          dto.UserID,
	}, nil
}

func (c *CourseConverter) ManyToDTO(courses []*internal.Course) ([]internal.CourseDTO, error) {
	result := make([]internal.CourseDTO, 0)

	for _, course := range courses {
		if course.UserID == 0 {
			continue
		}
		dto, err := c.ToDTO(course)
		if err != nil {
			return nil, errors.Wrapf(err, "while converting course %d", course.ID)
		}
		result = append(result, *dto)
	}

	return result, nil
}
