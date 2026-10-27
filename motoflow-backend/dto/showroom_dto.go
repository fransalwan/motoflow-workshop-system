package dto

import (
	"github.com/shopspring/decimal"
)

type CreditSimulationRequest struct {
	OTRPrice    decimal.Decimal `json:"otr_price" binding:"required"`
	DPPercent   decimal.Decimal `json:"dp_percent" binding:"required"`
	TenorMonths int             `json:"tenor_months" binding:"required"`
}

type CreateSPKRequest struct {
	CustomerID         uint            `json:"customer_id" binding:"required"`
	MotorUnitID        uint            `json:"motor_unit_id" binding:"required"`
	PurchaseType       string          `json:"purchase_type" binding:"required"` // cash, credit_nsc
	BookingFee         decimal.Decimal `json:"booking_fee"`
	DownPayment        decimal.Decimal `json:"down_payment"`
	TenorMonths        int             `json:"tenor_months"`
	MonthlyInstallment decimal.Decimal `json:"monthly_installment"`
	Notes              string          `json:"notes"`
}
