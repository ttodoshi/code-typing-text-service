package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
)

type codeExampleRepository struct {
	DB *gorm.DB
}

func NewCodeExampleRepository() ports.CodeExampleRepository {
	return &codeExampleRepository{DB: database.ConnectToDb()}
}

func (r *codeExampleRepository) GetProgrammingLanguages() (programmingLanguages []domain.ProgrammingLanguage) {
	r.DB.Find(&programmingLanguages)
	return
}

func (r *codeExampleRepository) GetCodeExamples() (codeExamples []domain.CodeExample) {
	r.DB.Find(&codeExamples)
	return
}

func (r *codeExampleRepository) GetCodeExamplesByProgrammingLanguageUUID(programmingLanguageUUID string) (codeExamples []domain.CodeExample, err error) {
	var programmingLanguage domain.ProgrammingLanguage
	r.DB.First(&programmingLanguage, "uuid = ?", programmingLanguageUUID)

	if programmingLanguage.UUID == "" {
		return codeExamples, &errors.NotFoundError{
			Message: fmt.Sprintf(`programming programmingLanguage by uuid '%s' not found`, programmingLanguageUUID),
		}
	}

	r.DB.Find(&codeExamples, "programming_language_uuid = ?", programmingLanguage.UUID)
	return codeExamples, nil
}
