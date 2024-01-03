package servises

import (
	"github.com/jinzhu/copier"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/logging"
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
		return getProgrammingLanguagesDto, &errors.MappingError{Message: `struct mapping error`}
	}
	return getProgrammingLanguagesDto, nil
}

func (s *CodeExampleService) GetCodeExamples() (getCodeExamplesDto []dto.GetCodeExampleDto, err error) {
	codeExamples := s.repo.GetCodeExamples()
	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		return getCodeExamplesDto, &errors.MappingError{Message: `struct mapping error`}
	}
	return getCodeExamplesDto, nil
}

func (s *CodeExampleService) GetCodeExamplesByProgrammingLanguageUUID(programmingLanguageUUID string) (getCodeExamplesDto []dto.GetCodeExampleDto, err error) {
	codeExamples, err := s.repo.GetCodeExamplesByProgrammingLanguageUUID(programmingLanguageUUID)
	if err != nil {
		s.log.Infof(`error getting code code examples by programming language uuid: '%s' due to error: %v`, programmingLanguageUUID, err)
		return getCodeExamplesDto, err
	}

	err = copier.Copy(&getCodeExamplesDto, &codeExamples)
	if err != nil {
		return getCodeExamplesDto, &errors.MappingError{Message: `struct mapping error`}
	}
	return getCodeExamplesDto, nil
}
