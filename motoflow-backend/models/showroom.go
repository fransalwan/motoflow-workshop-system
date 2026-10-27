package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MotorUnit struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	ModelName         string          `gorm:"size:100;not null;index" json:"model_name"` // Honda Vario 160 ABS, etc.
	Category          string          `gorm:"size:50;not null" json:"category"`          // Matic, Sport, Cub, EV
	FrameNumber       string          `gorm:"size:50;not null;uniqueIndex" json:"frame_number"` // VIN / No. Rangka
	EngineNumber      string          `gorm:"size:50;not null;uniqueIndex" json:"engine_number"` // No. Mesin
	Color             string          `gorm:"size:50;not null" json:"color"`
	Year              int             `gorm:"not null" json:"year"`
	OTRPrice          decimal.Decimal `gorm:"type:numeric(15,2);not null" json:"otr_price"`
	MinDP             decimal.Decimal `gorm:"type:numeric(15,2);not null" json:"min_dp"`
	Status            string          `gorm:"size:30;not null;default:'ready';index" json:"status"` // ready, booked, delivered
	ImageURL          string          `gorm:"size:255" json:"image_url"`
	SpecsSummary      string          `gorm:"type:text" json:"specs_summary"`
	WarehouseLocation string          `gorm:"size:100;default:'Showroom NSS Ciledug'" json:"warehouse_location"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`
}

type SPKUnit struct {
	ID                 uint            `gorm:"primaryKey" json:"id"`
	SPKNumber          string          `gorm:"size:50;not null;uniqueIndex" json:"spk_number"` // SPK-NSS-YYYYMMDD-XXXX
	CustomerID         uint            `gorm:"not null;index" json:"customer_id"`
	Customer           *Customer       `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	SalesID            uint            `gorm:"not null;index" json:"sales_id"`
	Sales              *User           `gorm:"foreignKey:SalesID" json:"sales,omitempty"`
	MotorUnitID        uint            `gorm:"not null;index" json:"motor_unit_id"`
	MotorUnit          *MotorUnit      `gorm:"foreignKey:MotorUnitID" json:"motor_unit,omitempty"`
	PurchaseType       string          `gorm:"size:30;not null;default:'cash'" json:"purchase_type"` // cash, credit_nsc
	BookingFee         decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"booking_fee"`
	DownPayment        decimal.Decimal `gorm:"type:numeric(15,2);not null;default:0" json:"down_payment"`
	TenorMonths        int             `gorm:"default:0" json:"tenor_months"` // 11, 17, 23, 29, 35
	MonthlyInstallment decimal.Decimal `gorm:"type:numeric(15,2);default:0" json:"monthly_installment"`
	LeasingStatus      string          `gorm:"size:30;default:'not_applicable'" json:"leasing_status"` // not_applicable, survey, approved, rejected
	Status             string          `gorm:"size:30;not null;default:'draft';index" json:"status"`   // draft, verified, delivered, cancelled
	Notes              string          `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`
}

