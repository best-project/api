package converter

import "github.com/best-project/api/internal"

type Converter struct {
	CourseConverter *CourseConverter
}

func NewConverter() *Converter {
	return &Converter{
		CourseConverter: NewCourseConverter(&TaskConverter{}),
	}
}

func (c *Converter) ToModel(dto *internal.UserDTO) *internal.User {
	return &internal.User{
		Username:       dto.Username,
		Points:         dto.Points,
		Avatar:         dto.Avatar,
		Level:          dto.Level,
		ProfileCourses: c.CourseConverter.ManyToModel(dto.ProfileCourses),
	}
}

func (u *Converter) ToDTO(dto *internal.User) *internal.UserDTO {
	return &internal.UserDTO{
		Username:       dto.Username,
		Password:       string(dto.Password),
		Level:          dto.Level,
		Avatar:         dto.Avatar,
		Points:         dto.Points,
		ProfileCourses: u.CourseConverter.ManyToDTO(dto.ProfileCourses),
	}
}
