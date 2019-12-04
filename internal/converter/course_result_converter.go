package converter

import (
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/pkg/errors"
)

type CourseResultConverter struct {
	db *storage.Database
}

func NewCourseResultConverter(db *storage.Database) *CourseResultConverter {
	return &CourseResultConverter{
		db: db,
	}
}

func (c *CourseResultConverter) ToModel(dto *internal.CourseResultDTO) *internal.CourseResult {
	return &internal.CourseResult{
		Points:   dto.Points,
		CourseID: dto.CourseID,
		UserID:   dto.UserID,
		Phase:    dto.Phase,
		Passed:   dto.Passed,
	}
}

func (c *CourseResultConverter) FetchIDs(m []internal.CourseResult) []string {
	result := make([]string, 0)
	for _, item := range m {
		result = append(result, item.CourseID)
	}

	return unique(result)
}

func (c *CourseResultConverter) ManyToModel(dtos []internal.CourseResultDTO) []internal.CourseResult {
	result := make([]internal.CourseResult, len(dtos))

	for _, dto := range dtos {
		result = append(result, *c.ToModel(&dto))
	}

	return result
}

func (c *CourseResultConverter) ToDTO(dto *internal.CourseResult) (*internal.CourseResultDTO, error) {
	return &internal.CourseResultDTO{
		CourseID: dto.CourseID,
		Points:   dto.Points,
		UserID:   dto.UserID,
		Phase:    dto.Phase,
		Passed:   dto.Passed,
	}, nil
}

func (c *CourseResultConverter) ManyToDTO(courses []internal.CourseResult) ([]internal.CourseResultDTO, error) {
	result := make([]internal.CourseResultDTO, len(courses))

	for _, course := range courses {
		dto, err := c.ToDTO(&course)
		if err != nil {
			return nil, errors.Wrapf(err, "while converting course %d", course.ID)
		}
		result = append(result, *dto)
	}

	return result, nil
}
