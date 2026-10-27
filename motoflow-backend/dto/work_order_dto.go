package dto

import (
	"github.com/shopspring/decimal"
	"motoflow-backend/models"
)

type WorkOrderItemInput struct {
	SparepartID     *uint           `json:"sparepart_id"`
	ItemType        models.ItemType `json:"item_type"` // "part" or "service"
	PartNumber      string          `json:"part_number"`
	ItemName        string          `json:"item_name" binding:"required"`
	Quantity        int             `json:"quantity" binding:"required,min=1"`
	Price           decimal.Decimal `json:"price" binding:"required"`
	DiscountPercent decimal.Decimal `json:"discount_percent"`
	MechanicNotes   string          `json:"mechanic_notes"`
}

type CreateWorkOrderRequest struct {
	CustomerID       uint                 `json:"customer_id" binding:"required"`
	MechanicID       *uint                `json:"mechanic_id"`
	EntryKM          int                  `json:"entry_km"`
	Complaints       string               `json:"complaints"`
	ServiceDiagnosis string               `json:"service_diagnosis"`
	TaxRate          *decimal.Decimal     `json:"tax_rate"` // e.g., 11.00
	TotalDiscount    *decimal.Decimal     `json:"total_discount"`
	Notes            string               `json:"notes"`
	Items            []WorkOrderItemInput `json:"items" binding:"required,min=1"`
}

type UpdateWorkOrderStatusRequest struct {
	Status models.WorkOrderStatus `json:"status" binding:"required"`
	Notes  string                 `json:"notes"`
}

type ProcessPaymentRequest struct {
	PaymentMethod models.PaymentMethod `json:"payment_method" binding:"required"`
	PaidAmount    decimal.Decimal      `json:"paid_amount" binding:"required"`
	Notes         string               `json:"notes"`
}

