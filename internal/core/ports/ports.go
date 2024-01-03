package ports

import (
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/domain"
)

type CodeExampleService interface {
	GetProgrammingLanguages() ([]dto.GetProgrammingLanguageDto, error)
	GetCodeExamples() ([]dto.GetCodeExampleDto, error)
	GetCodeExamplesByProgrammingLanguageName(programmingLanguageName string) ([]dto.GetCodeExampleDto, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=CodeExampleRepository
type CodeExampleRepository interface {
	GetProgrammingLanguages() []domain.ProgrammingLanguage
	GetCodeExamples() []domain.CodeExample
	GetCodeExamplesByProgrammingLanguageName(programmingLanguageName string) ([]domain.CodeExample, error)
}

type RegularTextService interface {
	GetRegularTexts() ([]dto.GetRegularTextDto, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=RegularTextRepository
type RegularTextRepository interface {
	GetRegularTexts() []domain.RegularText
}
