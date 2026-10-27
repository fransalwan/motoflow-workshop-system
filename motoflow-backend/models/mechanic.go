package models

import (
	"time"

	"gorm.io/gorm"
)

type Mechanic struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"size:30;not null;uniqueIndex" json:"code"` // Contoh: MEK-01
	Name           string         `gorm:"size:100;not null;index" json:"name"`
	Phone          string         `gorm:"size:25" json:"phone"`
	Specialization string         `gorm:"size:100" json:"specialization"` // Injeksi PGM-FI, CVT Matic, Mesin Bebek/Sport
	PitNumber      int            `gorm:"default:1" json:"pit_number"`    // Nomor Pit Servis AHASS (Pit 1, Pit 2, dst)
	IsAvailable    bool           `gorm:"default:true" json:"is_available"`
	WorkOrders     []WorkOrder    `gorm:"foreignKey:MechanicID" json:"work_orders,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

