package dto

import (
	"github.com/shopspring/decimal"
	"motoflow-backend/models"
)

type CreateSparepartRequest struct {
	PartNumber   string          `json:"part_number" binding:"required"`
	Name         string          `json:"name" binding:"required"`
	CategoryID   uint            `json:"category_id" binding:"required"`
	Type         models.ItemType `json:"type" binding:"required"`
	CostPrice    decimal.Decimal `json:"cost_price"`
	SellingPrice decimal.Decimal `json:"selling_price" binding:"required"`
	Stock        int             `json:"stock"`
	MinStock     int             `json:"min_stock"`
	Unit         string          `json:"unit"`
	Location     string          `json:"location"`
	IsActive     bool            `json:"is_active"`
}

type UpdateSparepartRequest struct {
	PartNumber   string          `json:"part_number"`
	Name         string          `json:"name"`
	CategoryID   uint            `json:"category_id"`
	Type         models.ItemType `json:"type"`
	CostPrice    decimal.Decimal `json:"cost_price"`
	SellingPrice decimal.Decimal `json:"selling_price"`
	Stock        int             `json:"stock"`
	MinStock     int             `json:"min_stock"`
	Unit         string          `json:"unit"`
	Location     string          `json:"location"`
	IsActive     *bool           `json:"is_active"`
}

type UpdateStockRequest struct {
	Quantity int    `json:"quantity" binding:"required"` // positive to add, negative to deduct
	Reason   string `json:"reason"`
}

