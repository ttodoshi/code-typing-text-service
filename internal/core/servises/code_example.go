package servises

import (
	"code-typing-text-service/internal/core/domain"
	"code-typing-text-service/internal/core/ports"
	"code-typing-text-service/internal/core/ports/dto"
	"code-typing-text-service/pkg/logging"
	"fmt"
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
		err = fmt.Errorf(`struct mapping error: %w`, ports.InternalServerError)
	}
	return
}

func (s *CodeExampleService) GetCodeExampleByUUID(userID, UUID string) (dto.GetCodeExampleDto, error) {
	codeExample, err := s.repo.GetCodeExampleByUUID(UUID)
	if err != nil || codeExample.UserID != userID && codeExample.UserID != "" {
		return dto.GetCodeExampleDto{}, fmt.Errorf(`code example not found: %w`, ports.NotFoundError)
	}

	var getCodeExampleDto dto.GetCodeExampleDto
	err = copier.Copy(&getCodeExampleDto, &codeExample)
	if err != nil {
		err = fmt.Errorf(`struct mapping error: %w`, ports.InternalServerError)
	}
	return getCodeExampleDto, nil
}

func (s *CodeExampleService) GetCodeExamples(userID string) (getCodeExamplesDto []dto.GetCodeExampleDto, err error) {
	codeExamples := s.repo.GetCodeExamples(userID)

	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		err = fmt.Errorf(`struct mapping error: %w`, ports.InternalServerError)
	}
	return
}

func (s *CodeExampleService) GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName string) ([]dto.GetCodeExampleDto, error) {
	codeExamples, err := s.repo.GetCodeExamplesByProgrammingLanguageName(userID, programmingLanguageName)
	if err != nil {
		s.log.Infof(`error getting code examples by programming language '%s' due to error: %v`, programmingLanguageName, err)
		err = fmt.Errorf(`code examples not found: %w`, ports.NotFoundError)
		return nil, err
	}

	var getCodeExamplesDto []dto.GetCodeExampleDto
	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		return nil, fmt.Errorf(`struct mapping error: %w`, ports.InternalServerError)
	}
	return getCodeExamplesDto, nil
}

func (s *CodeExampleService) CreateCodeExample(userID string, createCodeExampleDto dto.CreateCodeExampleDto) (string, error) {
	UUID, err := s.repo.CreateCodeExample(
		domain.CodeExample{
			UserID:                  userID,
			ProgrammingLanguageUUID: createCodeExampleDto.ProgrammingLanguageUUID,
			Content:                 createCodeExampleDto.Content,
		},
	)
	if err != nil {
		s.log.Infof(`error creating code example due to error: %v`, err)
		return "", fmt.Errorf(`error creating code example: %w`, ports.BadRequestError)
	}

	return UUID, nil
}

func (s *CodeExampleService) DeleteCodeExample(userID, UUID string) error {
	codeExample, err := s.repo.GetCodeExampleByUUID(UUID)
	if err != nil || codeExample.UserID != userID {
		return fmt.Errorf(`code example not found: %w`, ports.NotFoundError)
	}
	s.repo.DeleteCodeExample(UUID)
	return nil
}
