package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WorkOrderStatus string

const (
	WOStatusDraft      WorkOrderStatus = "draft"
	WOStatusInProgress WorkOrderStatus = "in_progress"
	WOStatusCompleted  WorkOrderStatus = "completed"
	WOStatusPaid       WorkOrderStatus = "paid"
	WOStatusCancelled  WorkOrderStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentCash     PaymentMethod = "cash"
	PaymentTransfer PaymentMethod = "transfer"
	PaymentQRIS     PaymentMethod = "qris"
	PaymentDebit    PaymentMethod = "debit"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid PaymentStatus = "unpaid"
	PaymentStatusPaid   PaymentStatus = "paid"
)

type WorkOrder struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	WONumber         string          `gorm:"size:50;not null;uniqueIndex" json:"wo_number"` // Format: WO-AHASS-YYYYMMDD-XXXX
	CustomerID       uint            `gorm:"not null;index" json:"customer_id"`
	Customer         *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	MechanicID       *uint           `gorm:"index" json:"mechanic_id"`
	Mechanic         *Mechanic       `gorm:"foreignKey:MechanicID" json:"mechanic,omitempty"`
	ServiceAdvisorID *uint           `gorm:"index" json:"service_advisor_id"`
	ServiceAdvisor   *User           `gorm:"foreignKey:ServiceAdvisorID" json:"service_advisor,omitempty"`
	EntryKM          int             `json:"entry_km"`
	Complaints       string          `gorm:"type:text" json:"complaints"`         // Keluhan Konsumen
	ServiceDiagnosis string          `gorm:"type:text" json:"service_diagnosis"`   // Analisa & Diagnosa SA
	Status           WorkOrderStatus `gorm:"size:30;not null;default:'draft';index" json:"status"`
	PaymentMethod    PaymentMethod   `gorm:"size:30;default:'cash'" json:"payment_method"`
	PaymentStatus    PaymentStatus   `gorm:"size:20;not null;default:'unpaid';index" json:"payment_status"`

	// Financial Calculations (Decimal Precision)
	PartsSubtotal    decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"parts_subtotal"`
	ServicesSubtotal decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"services_subtotal"`
	TotalDiscount    decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"total_discount"`
	TaxRate          decimal.Decimal `gorm:"type:decimal(5,2);default:11" json:"tax_rate"` // PPN 11%
	TaxAmount        decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"tax_amount"`
	GrandTotal       decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"grand_total"`
	PaidAmount       decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"paid_amount"`
	ChangeAmount     decimal.Decimal `gorm:"type:decimal(12,2);default:0" json:"change_amount"`

	CompletedAt      *time.Time      `json:"completed_at"`
	PaidAt           *time.Time      `json:"paid_at"`
	Notes            string          `gorm:"type:text" json:"notes"`
	Items            []WorkOrderItem `gorm:"foreignKey:WorkOrderID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        gorm.DeletedAt  `gorm:"index" json:"-"`
}

