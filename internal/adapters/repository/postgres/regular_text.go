package postgres

import (
	"gorm.io/gorm"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
)

type regularTextRepository struct {
	DB *gorm.DB
}

func NewRegularTextRepository() ports.RegularTextRepository {
	return &regularTextRepository{DB: database.ConnectToDb()}
}

func (r *regularTextRepository) GetRegularTexts() (texts []domain.RegularText) {
	r.DB.Find(&texts)
	return
}
