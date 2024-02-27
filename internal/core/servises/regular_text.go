package servises

import (
	"github.com/jinzhu/copier"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/errors"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/logging"
)

type RegularTextService struct {
	repo ports.RegularTextRepository
	log  logging.Logger
}

func NewRegularTextService(repo ports.RegularTextRepository, log logging.Logger) ports.RegularTextService {
	return &RegularTextService{
		repo: repo,
		log:  log,
	}
}

func (s *RegularTextService) GetRegularTexts() (getTextsDto []dto.GetRegularTextDto, err error) {
	texts := s.repo.GetRegularTexts()

	err = copier.Copy(&getTextsDto, &texts)
	if err != nil {
		err = &errors.MappingError{Message: `struct mapping error`}
	}
	return
}
