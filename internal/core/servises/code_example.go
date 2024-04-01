package servises

import (
	"code-typing-text-service/internal/adapters/dto"
	"code-typing-text-service/internal/core/domain"
	"code-typing-text-service/internal/core/errors"
	"code-typing-text-service/internal/core/ports"
	"code-typing-text-service/pkg/logging"
	"github.com/jinzhu/copier"
)

type CodeExampleService struct {
	repo ports.CodeExampleRepository
	log  logging.Logger
}

func NewCodeExampleService(repo ports.CodeExampleRepository, log logging.Logger) ports.CodeExampleService {
	return &CodeExampleService{
		repo: repo,
		log:  log,
	}
}

func (s *CodeExampleService) GetProgrammingLanguages() (getProgrammingLanguagesDto []dto.GetProgrammingLanguageDto, err error) {
	programmingLanguages := s.repo.GetProgrammingLanguages()

	err = copier.Copy(&getProgrammingLanguagesDto, &programmingLanguages)
	if err != nil {
		err = &errors.MappingError{Message: `struct mapping error`}
	}
	return
}

func (s *CodeExampleService) GetCodeExampleByUUID(userID, UUID string) (getCodeExampleDto dto.GetCodeExampleDto, err error) {
	var codeExample domain.CodeExample
	codeExample, err = s.repo.GetCodeExampleByUUID(UUID)
	if err != nil {
		return
	}
	if codeExample.UserID != userID && codeExample.UserID != "" {
		err = &errors.NotFoundError{Message: `code example not found`}
		return
	}

	err = copier.Copy(&getCodeExampleDto, &codeExample)
	if err != nil {
		err = &errors.MappingError{Message: `struct mapping error`}
	}
	return
}

func (s *CodeExampleService) GetCodeExamples(userID string) (getCodeExamplesDto []dto.GetCodeExampleDto, err error) {
	codeExamples := s.repo.GetCodeExamples(userID)

	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		err = &errors.MappingError{Message: `struct mapping error`}
	}
	return
}

func (s *CodeExampleService) GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName string) (getCodeExamplesDto []dto.GetCodeExampleDto, err error) {
	codeExamples, err := s.repo.GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName)
	if err != nil {
		s.log.Infof(`error getting code examples by programming language '%s' due to error: %v`, programmingLanguageName, err)
		return
	}

	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		err = &errors.MappingError{Message: `struct mapping error`}
	}
	return
}

func (s *CodeExampleService) CreateCodeExample(userID string, createCodeExampleDto dto.CreateCodeExampleDto) (string, error) {
	UUID, err := s.repo.SaveCodeExample(
		domain.CodeExample{
			UserID:                  userID,
			ProgrammingLanguageUUID: createCodeExampleDto.ProgrammingLanguageUUID,
			Content:                 createCodeExampleDto.Content,
		},
	)
	if err != nil {
		s.log.Infof(`error creating code example due to error: %v`, err)
		return "", err
	}

	return UUID, nil
}

func (s *CodeExampleService) DeleteCodeExample(userID, UUID string) (err error) {
	codeExample, err := s.repo.GetCodeExampleByUUID(UUID)
	if err != nil {
		return
	}
	if codeExample.UserID != userID {
		err = &errors.NoAccessError{Message: `no access to code example`}
		return
	}
	err = s.repo.DeleteCodeExample(UUID)
	if err != nil {
		s.log.Infof(`error deleting code example due to error: %v`, err)
	}
	return err
}
