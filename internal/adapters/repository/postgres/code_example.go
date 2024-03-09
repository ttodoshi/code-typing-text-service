package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
	"speed-typing-text-service/pkg/logging"
)

type codeExampleRepository struct {
	log logging.Logger
	db  *gorm.DB
}

func NewCodeExampleRepository(log logging.Logger) ports.CodeExampleRepository {
	db := database.ConnectToDb()
	err := db.AutoMigrate(&domain.ProgrammingLanguage{}, &domain.CodeExample{})
	if err != nil {
		log.Fatal("error while migrating")
	}
	return &codeExampleRepository{db: db}
}

func (r *codeExampleRepository) GetProgrammingLanguages() (programmingLanguages []domain.ProgrammingLanguage) {
	r.db.Find(&programmingLanguages)
	return
}

func (r *codeExampleRepository) GetCodeExampleByUUID(UUID string) (codeExample domain.CodeExample, err error) {
	r.db.Find(&codeExample, "uuid = ?", UUID)
	if codeExample.UUID == "" {
		return codeExample, &errors.NotFoundError{
			Message: fmt.Sprintf("code example by uuid '%s' not found", UUID),
		}
	}
	return codeExample, nil
}

func (r *codeExampleRepository) GetCodeExamples() (codeExamples []domain.CodeExample) {
	r.db.Find(&codeExamples)
	return
}

func (r *codeExampleRepository) GetCodeExamplesByProgrammingLanguageName(programmingLanguageName string) (codeExamples []domain.CodeExample, err error) {
	var programmingLanguage domain.ProgrammingLanguage
	r.db.First(&programmingLanguage, "name = ?", programmingLanguageName)

	if programmingLanguage.UUID == "" {
		return codeExamples, &errors.NotFoundError{
			Message: fmt.Sprintf(`programming language '%s' not found`, programmingLanguageName),
		}
	}

	r.db.Find(&codeExamples, "programming_language_uuid = ?", programmingLanguage.UUID)
	return codeExamples, nil
}
