package postgres

import (
	"gorm.io/gorm"
	"log"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
)

type regularTextRepository struct {
	DB *gorm.DB
}

func NewRegularTextRepository() ports.RegularTextRepository {
	DB := database.ConnectToDb()
	err := DB.AutoMigrate(&domain.RegularText{})
	if err != nil {
		log.Fatal("error while migrating")
	}
	return &regularTextRepository{DB: DB}
}

func (r *regularTextRepository) GetRegularTexts() (texts []domain.RegularText) {
	r.DB.Find(&texts)
	return
}
