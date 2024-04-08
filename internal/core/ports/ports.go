package ports

import (
	"code-typing-text-service/internal/core/domain"
	"code-typing-text-service/internal/core/ports/dto"
)

type CodeExampleService interface {
	GetProgrammingLanguages() ([]dto.GetProgrammingLanguageDto, error)
	GetCodeExampleByUUID(userID, UUID string) (dto.GetCodeExampleDto, error)
	GetCodeExamples(userID string) ([]dto.GetCodeExampleDto, error)
	GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName string) ([]dto.GetCodeExampleDto, error)
	CreateCodeExample(userID string, createCodeExampleDto dto.CreateCodeExampleDto) (string, error)
	DeleteCodeExample(userID, UUID string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.39.1 --name=CodeExampleRepository
type CodeExampleRepository interface {
	GetProgrammingLanguages() []domain.ProgrammingLanguage
	GetCodeExampleByUUID(UUID string) (domain.CodeExample, error)
	GetCodeExamples(userID string) []domain.CodeExample
	GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName string) ([]domain.CodeExample, error)
	SaveCodeExample(codeExample domain.CodeExample) (string, error)
	DeleteCodeExample(UUID string) error
}
