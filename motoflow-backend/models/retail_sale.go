package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type RetailSale struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	InvoiceNumber  string          `gorm:"size:50;not null;uniqueIndex" json:"invoice_number"` // INV-NSS-YYYYMMDD-XXXX
	CustomerID     *uint           `gorm:"index" json:"customer_id"`
	Customer       *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	CustomerName   string          `gorm:"size:100;default:'Konsumen Walk-in'" json:"customer_name"`
	CustomerPhone  string          `gorm:"size:30" json:"customer_phone"`
	CashierID      uint            `gorm:"not null;index" json:"cashier_id"`
	Cashier        *User           `gorm:"foreignKey:CashierID" json:"cashier,omitempty"`
	Subtotal       decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"subtotal"`
	DiscountAmount decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"discount_amount"`
	TaxPPN         decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"tax_ppn"`
	GrandTotal     decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"grand_total"`
	PaidAmount     decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"paid_amount"`
	ChangeAmount   decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"change_amount"`
	PaymentMethod  PaymentMethod   `gorm:"size:30;not null;default:'cash'" json:"payment_method"`
	PaymentStatus  PaymentStatus   `gorm:"size:20;not null;default:'paid'" json:"payment_status"`
	PromoCode      string          `gorm:"size:50" json:"promo_code"`
	Notes          string          `gorm:"type:text" json:"notes"`
	Items          []RetailSaleItem `gorm:"foreignKey:RetailSaleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      gorm.DeletedAt  `gorm:"index" json:"-"`
}

type RetailSaleItem struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	RetailSaleID       uint            `gorm:"not null;index" json:"retail_sale_id"`
	SparepartID        uint            `gorm:"not null;index" json:"sparepart_id"`
	Sparepart          *Sparepart      `gorm:"foreignKey:SparepartID" json:"sparepart,omitempty"`
	PartNumberSnapshot string          `gorm:"size:50;not null" json:"part_number_snapshot"`
	PartNameSnapshot   string          `gorm:"size:150;not null" json:"part_name_snapshot"`
	Quantity           int             `gorm:"not null;default:1" json:"quantity"`
	PriceSnapshot      decimal.Decimal `gorm:"type:numeric(15,2);not null" json:"price_snapshot"`
	DiscountPercent    decimal.Decimal `gorm:"type:numeric(5,2);not null;default:0" json:"discount_percent"`
	Subtotal           decimal.Decimal `gorm:"type:numeric(15,2);not null" json:"subtotal"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}
