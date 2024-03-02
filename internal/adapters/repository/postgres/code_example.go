package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
)

type codeExampleRepository struct {
	DB *gorm.DB
}

func NewCodeExampleRepository() ports.CodeExampleRepository {
	DB := database.ConnectToDb()
	err := DB.AutoMigrate(&domain.ProgrammingLanguage{}, &domain.CodeExample{})
	if err != nil {
		log.Fatal("error while migrating")
	}
	return &codeExampleRepository{DB: DB}
}

func (r *codeExampleRepository) GetProgrammingLanguages() (programmingLanguages []domain.ProgrammingLanguage) {
	r.DB.Find(&programmingLanguages)
	return
}

func (r *codeExampleRepository) GetCodeExampleByUUID(UUID string) (codeExample domain.CodeExample, err error) {
	r.DB.Find(&codeExample, "uuid = ?", UUID)
	if codeExample.UUID == "" {
		return codeExample, &errors.NotFoundError{
			Message: fmt.Sprintf("code example by uuid '%s' not found", UUID),
		}
	}
	return codeExample, nil
}

func (r *codeExampleRepository) GetCodeExamples() (codeExamples []domain.CodeExample) {
	r.DB.Find(&codeExamples)
	return
}

func (r *codeExampleRepository) GetCodeExamplesByProgrammingLanguageName(programmingLanguageName string) (codeExamples []domain.CodeExample, err error) {
	var programmingLanguage domain.ProgrammingLanguage
	r.DB.First(&programmingLanguage, "name = ?", programmingLanguageName)

	if programmingLanguage.UUID == "" {
		return codeExamples, &errors.NotFoundError{
			Message: fmt.Sprintf(`programming language '%s' not found`, programmingLanguageName),
		}
	}

	r.DB.Find(&codeExamples, "programming_language_uuid = ?", programmingLanguage.UUID)
	return codeExamples, nil
}
