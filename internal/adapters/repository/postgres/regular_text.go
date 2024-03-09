package postgres

import (
	"gorm.io/gorm"
	"speed-typing-text-service/internal/core/domain"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/pkg/database"
	"speed-typing-text-service/pkg/logging"
)

type regularTextRepository struct {
	log logging.Logger
	db  *gorm.DB
}

func NewRegularTextRepository(log logging.Logger) ports.RegularTextRepository {
	db := database.ConnectToDb()
	err := db.AutoMigrate(&domain.RegularText{})
	if err != nil {
		log.Fatal("error while migrating")
	}
	return &regularTextRepository{db: db}
}

func (r *regularTextRepository) GetRegularTexts() (texts []domain.RegularText) {
	r.db.Find(&texts)
	return
}
