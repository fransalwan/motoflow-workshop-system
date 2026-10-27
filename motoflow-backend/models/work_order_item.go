package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WorkOrderItem struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	WorkOrderID        uint            `gorm:"not null;index" json:"work_order_id"`
	WorkOrder          *WorkOrder      `gorm:"foreignKey:WorkOrderID" json:"-"`
	SparepartID        *uint           `gorm:"index" json:"sparepart_id"`
	Sparepart          *Sparepart      `gorm:"foreignKey:SparepartID" json:"sparepart,omitempty"`
	ItemType           ItemType        `gorm:"size:20;not null;default:'part'" json:"item_type"` // "part" or "service"
	PartNumberSnapshot string          `gorm:"size:50;not null" json:"part_number_snapshot"`
	ItemNameSnapshot   string          `gorm:"size:150;not null" json:"item_name_snapshot"`
	Quantity           int             `gorm:"not null;default:1" json:"quantity"`
	PriceSnapshot      decimal.Decimal `gorm:"type:decimal(12,2);not null;default:0" json:"price_snapshot"`
	DiscountPercent    decimal.Decimal `gorm:"type:decimal(5,2);default:0" json:"discount_percent"`
	Subtotal           decimal.Decimal `gorm:"type:decimal(12,2);not null;default:0" json:"subtotal"`
	MechanicNotes      string          `gorm:"size:255" json:"mechanic_notes"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`
}

