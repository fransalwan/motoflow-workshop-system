package dto

import (
	"github.com/shopspring/decimal"
)

type RetailItemRequest struct {
	SparepartID     uint            `json:"sparepart_id" binding:"required"`
	Quantity        int             `json:"quantity" binding:"required,min=1"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
}

type RetailCheckoutRequest struct {
	CustomerID    *uint               `json:"customer_id"`
	CustomerName  string              `json:"customer_name"`
	CustomerPhone string              `json:"customer_phone"`
	PaymentMethod string              `json:"payment_method" binding:"required"` // cash, qris, debit, transfer
	PaidAmount    decimal.Decimal     `json:"paid_amount" binding:"required"`
	PromoCode     string              `json:"promo_code"`
	Notes         string              `json:"notes"`
	Items         []RetailItemRequest `json:"items" binding:"required,min=1"`
}

