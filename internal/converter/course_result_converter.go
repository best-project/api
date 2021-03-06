package converter

import (
	"github.com/best-project/api/internal"
	"github.com/pkg/errors"
	"strconv"
)

type CourseResultConverter struct{}

func NewCourseResultConverter() *CourseResultConverter {
	return &CourseResultConverter{}
}

func (c *CourseResultConverter) ToModel(dto *internal.CourseResultDTO) *internal.CourseResult {
	id, _ := strconv.Atoi(dto.CourseID)
	return &internal.CourseResult{
		Points:   dto.Points,
		CourseID: uint(id),
		UserID:   dto.UserID,
		Phase:    dto.Phase,
		Passed:   dto.Passed,
	}
}

func (c *CourseResultConverter) FetchIDs(m []internal.CourseResult) []string {
	result := make([]string, 0)
	for _, item := range m {
		result = append(result, strconv.Itoa(int(item.CourseID)))
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
		CourseID: strconv.Itoa(int(dto.CourseID)),
		Points:   dto.Points,
		UserID:   dto.UserID,
		Phase:    dto.Phase,
		Passed:   dto.Passed,
	}, nil
}

func (c *CourseResultConverter) ManyToDTO(courses []internal.CourseResult) ([]internal.CourseResultDTO, error) {
	result := make([]internal.CourseResultDTO, 0)

	for _, course := range courses {
		dto, err := c.ToDTO(&course)
		if err != nil {
			return nil, errors.Wrapf(err, "while converting course %d", course.ID)
		}
		result = append(result, *dto)
	}

	return result, nil
}
