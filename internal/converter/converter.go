package converter

import (
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Converter struct {
	CourseConverter *CourseConverter
	db              *storage.Database
}

func NewConverter(db *storage.Database) *Converter {
	return &Converter{
		db:              db,
		CourseConverter: NewCourseConverter(db, &TaskConverter{}),
	}
}

func (c *Converter) FetchIDs(m []gorm.Model) []uint {
	result := make([]uint, 0)
	for _, item := range m {
		result = append(result, item.ID)
	}

	return result
}

func (c *Converter) ToModel(dto *internal.UserDTO) *internal.User {
	return &internal.User{
		Username: dto.Username,
		Points:   dto.Points,
		Avatar:   dto.Avatar,
		Level:    dto.Level,
	}
}

func (u *Converter) ToDTO(dto *internal.User) (*internal.UserDTO, error) {
	profileCoursesDTO, err := u.CourseConverter.ManyToDTO(dto.ProfileCourses)
	if err != nil {
		return nil, errors.Wrap(err, "while converting profile courses")
	}

	return &internal.UserDTO{
		Username:       dto.Username,
		Password:       string(dto.Password),
		Level:          dto.Level,
		Avatar:         dto.Avatar,
		Points:         dto.Points,
		ProfileCourses: profileCoursesDTO,
	}, nil
}
