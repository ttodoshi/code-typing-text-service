package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CodeExample struct {
	UUID                    string `gorm:"primaryKey"`
	Content                 string `gorm:"not null"`
	ProgrammingLanguageUUID string `gorm:"not null;index"`
	ProgrammingLanguage     ProgrammingLanguage
}

func (e *CodeExample) BeforeCreate(_ *gorm.DB) (err error) {
	e.UUID = uuid.NewString()
	return
}

type ProgrammingLanguage struct {
	UUID string `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
	Logo string `gorm:"unique;not null"`
}

func (l *ProgrammingLanguage) BeforeCreate(_ *gorm.DB) (err error) {
	l.UUID = uuid.NewString()
	return
}
