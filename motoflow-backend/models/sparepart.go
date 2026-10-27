package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ItemType string

const (
	ItemTypePart    ItemType = "part"
	ItemTypeService ItemType = "service"
)

type Sparepart struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	PartNumber   string          `gorm:"size:50;not null;uniqueIndex" json:"part_number"` // Contoh: 08232-2MB-K0NZ1
	Name         string          `gorm:"size:150;not null;index" json:"name"`
	CategoryID   uint            `gorm:"not null;index" json:"category_id"`
	Category     *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Type         ItemType        `gorm:"size:20;not null;default:'part'" json:"type"` // "part" atau "service"
	CostPrice    decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"cost_price"`
	SellingPrice decimal.Decimal `gorm:"type:decimal(12,2);not null;default:0" json:"selling_price"`
	Stock        int             `gorm:"not null;default:0" json:"stock"`
	MinStock     int             `gorm:"not null;default:5" json:"min_stock"` // Alert threshold
	Unit         string          `gorm:"size:20;default:'Pcs'" json:"unit"`   // Pcs, Botol, Set, Paket
	Location     string          `gorm:"size:60" json:"location"`             // Rak A-01, Area Bengkel
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}

