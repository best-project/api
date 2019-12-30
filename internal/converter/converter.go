package converter

import (
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type Converter struct {
	CourseConverter       *CourseConverter
	CourseResultConverter *CourseResultConverter
}

func NewConverter() *Converter {
	return &Converter{
		CourseConverter:       NewCourseConverter(&TaskConverter{}),
		CourseResultConverter: NewCourseResultConverter(),
	}
}

func (c *Converter) FetchIDs(m []gorm.Model) []uint {
	result := make([]uint, 0)
	for _, item := range m {
		result = append(result, item.ID)
	}

	return result
}

func (c *Converter) ToUserStat(dto internal.User) *internal.UserStatDTO {
	return &internal.UserStatDTO{
		Points: dto.Points,
		UserID: dto.ID,
		Level:  dto.Level,
		Email:  dto.Email,
		Avatar: dto.Avatar,
	}
}

func (c *Converter) ManyToUserStat(dtos []internal.User) []*internal.UserStatDTO {
	results := make([]*internal.UserStatDTO, 0)

	for _, dto := range dtos {
		results = append(results, c.ToUserStat(dto))
	}
	return results
}

func (c *Converter) ToModel(dto *internal.UserDTO) *internal.User {
	return &internal.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
		NextLevel: dto.NextLevel,
		Points:    dto.Points,
		Avatar:    dto.Avatar,
		Level:     dto.Level,
	}
}

func (u *Converter) ToDTO(dto internal.User) (*internal.UserDTO, error) {
	return &internal.UserDTO{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  "",
		Level:     dto.Level,
		Avatar:    dto.Avatar,
		NextLevel: dto.NextLevel,
		Points:    dto.Points,
		Email:     dto.Email,
		ID:        int(dto.ID),
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
