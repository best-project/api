package converter

import "github.com/best-project/api/internal"

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
		Name:        dto.Name,
		Description: dto.Description,
		Data:        c.TaskConverter.ManyToModel(dto.Data),
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

func (c *CourseConverter) ToDTO(dto internal.Course) internal.CourseDTO {
	return internal.CourseDTO{
		Name:        dto.Name,
		Rate:        dto.Rate,
		Image:       dto.Image,
		Language:    dto.Language,
		MaxPoints:   dto.MaxPoints,
		Difficulty:  dto.Difficulty,
		Description: dto.Description,
		Data:        c.TaskConverter.ManyToDTO(dto.Data),
	}
}

func (c *CourseConverter) ManyToDTO(dtos []internal.Course) []internal.CourseDTO {
	result := make([]internal.CourseDTO, len(dtos))

	for _, dto := range dtos {
		result = append(result, c.ToDTO(dto))
	}

	return result
}
