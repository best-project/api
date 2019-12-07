package converter

import (
	"github.com/best-project/api/internal"
	"strconv"
)

type TaskConverter struct{}

func (*TaskConverter) ConvertToModel(dto internal.TaskDTO) internal.Task {
	return internal.Task{
		Image:     dto.Image,
		Translate: dto.Translate,
		Word:      dto.Word,
		CourseID:  dto.CourseID,
	}
}

func (t *TaskConverter) ManyToModel(dto []internal.TaskDTO) []internal.Task {
	model := make([]internal.Task, len(dto))

	for _, d := range dto {
		model = append(model, t.ConvertToModel(d))
	}

	return model
}

func (*TaskConverter) ConvertToDTO(dto internal.Task) internal.TaskDTO {
	return internal.TaskDTO{
		ID:        strconv.Itoa(int(dto.ID)),
		Image:     dto.Image,
		Translate: dto.Translate,
		Word:      dto.Word,
		CourseID:  dto.CourseID,
	}
}

func (t *TaskConverter) ManyToDTO(dto []internal.Task) []internal.TaskDTO {
	model := make([]internal.TaskDTO, 0)

	for _, d := range dto {
		if d.Word == "" {
			continue
		}
		model = append(model, t.ConvertToDTO(d))
	}

	return model
}
