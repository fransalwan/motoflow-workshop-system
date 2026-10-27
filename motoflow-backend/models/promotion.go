package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Promotion struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	PromoCode     string          `gorm:"size:50;not null;uniqueIndex" json:"promo_code"`
	Title         string          `gorm:"size:150;not null" json:"title"`
	Description   string          `gorm:"type:text" json:"description"`
	PromoType     string          `gorm:"size:30;not null;default:'percentage'" json:"promo_type"` // percentage, fixed_nominal, bundle
	DiscountValue decimal.Decimal `gorm:"type:numeric(15,2);not null" json:"discount_value"`
	MinSpend      decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"min_spend"`
	IsActive      bool            `gorm:"default:true;index" json:"is_active"`
	ValidUntil    time.Time       `json:"valid_until"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

