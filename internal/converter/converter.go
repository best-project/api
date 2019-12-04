package converter

import (
	"github.com/best-project/api/internal"
	"github.com/best-project/api/internal/storage"
	"github.com/jinzhu/gorm"
)

type Converter struct {
	CourseConverter       *CourseConverter
	CourseResultConverter *CourseResultConverter
	db                    *storage.Database
}

func NewConverter(db *storage.Database) *Converter {
	return &Converter{
		db:                    db,
		CourseConverter:       NewCourseConverter(db, &TaskConverter{}),
		CourseResultConverter: NewCourseResultConverter(db),
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
		Email:    dto.Email,
		Password: dto.Password,
		//ProfileCourses: c.CourseConverter.ManyToModel(dto.ProfileCourses),
		Points: dto.Points,
		Avatar: dto.Avatar,
		Level:  dto.Level,
	}
}

func (u *Converter) ToDTO(dto internal.User) (*internal.UserDTO, error) {
	//profileCoursesDTO, err := u.CourseConverter.ManyToDTO(dto.ProfileCourses)
	//if err != nil {
	//	return nil, errors.Wrap(err, "while converting profile courses")
	//}

	return &internal.UserDTO{
		Username: dto.Username,
		Password: "",
		Level:    dto.Level,
		Avatar:   dto.Avatar,
		Points:   dto.Points,
		Token:    dto.Token,
		//ProfileCourses: profileCoursesDTO,
		Email: dto.Email,
		ID:    int(dto.ID),
	}, nil
}

func unique(strings []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strings {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
