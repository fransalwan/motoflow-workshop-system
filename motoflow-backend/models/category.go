package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Code        string         `gorm:"size:30;not null;uniqueIndex" json:"code"`
	Description string         `gorm:"size:255" json:"description"`
	Spareparts  []Sparepart    `gorm:"foreignKey:CategoryID" json:"spareparts,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

