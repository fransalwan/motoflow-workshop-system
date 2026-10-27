package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `gorm:"size:100;not null;index" json:"name"`
	Phone        string         `gorm:"size:25;not null;index" json:"phone"`
	Email        string         `gorm:"size:100" json:"email"`
	Address      string         `gorm:"type:text" json:"address"`
	PlateNumber  string         `gorm:"size:20;not null;index" json:"plate_number"` // Contoh: B 3456 KRF
	MotorModel   string         `gorm:"size:60;not null;index" json:"motor_model"`  // Contoh: Vario 160, Beat CBS, PCX 160
	Year         int            `json:"year"`
	FrameNumber  string         `gorm:"size:50" json:"frame_number"`  // No Rangka
	EngineNumber string         `gorm:"size:50" json:"engine_number"` // No Mesin
	CurrentKM    int            `json:"current_km"`                   // Odometer terakhir
	WorkOrders   []WorkOrder    `gorm:"foreignKey:CustomerID" json:"work_orders,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

