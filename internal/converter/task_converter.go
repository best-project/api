package converter

import "github.com/best-project/api/internal"

type TaskConverter struct{}

func (*TaskConverter) ConvertToModel(dto internal.TaskDTO) internal.Task {
	return internal.Task{
		Image:     dto.Image,
		Translate: dto.Translate,
		Word:      dto.Word,
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
		Image:     dto.Image,
		Translate: dto.Translate,
		Word:      dto.Word,
	}
}

func (t *TaskConverter) ManyToDTO(dto []internal.Task) []internal.TaskDTO {
	model := make([]internal.TaskDTO, len(dto))

	for _, d := range dto {
		model = append(model, t.ConvertToDTO(d))
	}

	return model
}
