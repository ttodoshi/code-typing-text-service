package ports

import (
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/domain"
)

type CodeExampleService interface {
	GetProgrammingLanguages() ([]dto.GetProgrammingLanguageDto, error)
	GetCodeExamples() ([]dto.GetCodeExampleDto, error)
	GetCodeExamplesByProgrammingLanguageUUID(programmingLanguageUUID string) ([]dto.GetCodeExampleDto, error)
}

type CodeExampleRepository interface {
	GetProgrammingLanguages() []domain.ProgrammingLanguage
	GetCodeExamples() []domain.CodeExample
	GetCodeExamplesByProgrammingLanguageUUID(programmingLanguageUUID string) ([]domain.CodeExample, error)
}

type RegularTextService interface {
	GetRegularTexts() ([]dto.GetRegularTextDto, error)
}

type RegularTextRepository interface {
	GetRegularTexts() []domain.RegularText
}
