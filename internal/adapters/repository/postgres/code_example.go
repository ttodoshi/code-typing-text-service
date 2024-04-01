package postgres

import (
	"code-typing-text-service/internal/core/domain"
	"code-typing-text-service/internal/core/errors"
	"code-typing-text-service/internal/core/ports"
	"code-typing-text-service/pkg/database"
	"code-typing-text-service/pkg/logging"
	"fmt"
	"gorm.io/gorm"
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

func (r *codeExampleRepository) GetCodeExamples(userID string) (codeExamples []domain.CodeExample) {
	r.db.Find(&codeExamples, "user_id is null or user_id = ?", userID)
	return
}

func (r *codeExampleRepository) GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName string) (codeExamples []domain.CodeExample, err error) {
	var programmingLanguage domain.ProgrammingLanguage
	r.db.First(&programmingLanguage, "name = ?", programmingLanguageName)

	if programmingLanguage.UUID == "" {
		return codeExamples, &errors.NotFoundError{
			Message: fmt.Sprintf(`programming language '%s' not found`, programmingLanguageName),
		}
	}

	r.db.Where("user_id is null or user_id = ?", userID).Find(&codeExamples, "programming_language_uuid = ?", programmingLanguage.UUID)
	return codeExamples, nil
}

func (r *codeExampleRepository) SaveCodeExample(codeExample domain.CodeExample) (string, error) {
	var programmingLanguage domain.ProgrammingLanguage
	r.db.First(&programmingLanguage, "uuid = ?", codeExample.ProgrammingLanguageUUID)

	if programmingLanguage.UUID == "" {
		return "", &errors.NotFoundError{
			Message: fmt.Sprintf(`programming language not found`),
		}
	}

	r.db.Create(&codeExample)
	return codeExample.UUID, nil
}

func (r *codeExampleRepository) DeleteCodeExample(UUID string) error {
	r.db.Delete(&domain.CodeExample{}, "uuid = ?", UUID)
	return nil
}
