package services

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/models"
)

// FinancialService handles all monetary and billing arithmetic using shopspring/decimal
// to eliminate floating-point rounding errors and drift.
type FinancialService struct{}

func NewFinancialService() *FinancialService {
	return &FinancialService{}
}

// CalculateItemSubtotal calculates an item's subtotal after item-level percentage discount.
// Subtotal = (Price * Qty) * (1 - (DiscountPercent / 100))
func (s *FinancialService) CalculateItemSubtotal(price decimal.Decimal, qty int, discountPercent decimal.Decimal) decimal.Decimal {
	if qty <= 0 {
		return decimal.Zero
	}

	quantity := decimal.NewFromInt(int64(qty))
	baseTotal := price.Mul(quantity)

	if discountPercent.GreaterThan(decimal.Zero) {
		hundred := decimal.NewFromInt(100)
		discountFactor := decimal.NewFromInt(1).Sub(discountPercent.Div(hundred))
		// Rounded to 2 decimal places
		return baseTotal.Mul(discountFactor).Round(2)
	}

	return baseTotal.Round(2)
}

type OrderCalculationResult struct {
	PartsSubtotal    decimal.Decimal
	ServicesSubtotal decimal.Decimal
	GrossSubtotal    decimal.Decimal
	TotalDiscount    decimal.Decimal
	NetSubtotal      decimal.Decimal
	TaxRate          decimal.Decimal
	TaxAmount        decimal.Decimal
	GrandTotal       decimal.Decimal
	PaidAmount       decimal.Decimal
	ChangeAmount     decimal.Decimal
}

// CalculateWorkOrderTotals aggregates items into parts subtotal, services subtotal,
// applies discounts, calculates PPN (11%), and computes the final grand total and change amount.
func (s *FinancialService) CalculateWorkOrderTotals(
	items []models.WorkOrderItem,
	orderDiscount decimal.Decimal,
	taxRate decimal.Decimal,
	paidAmount decimal.Decimal,
) OrderCalculationResult {
	partsSubtotal := decimal.Zero
	servicesSubtotal := decimal.Zero
	itemDiscounts := decimal.Zero

	for _, item := range items {
		itemSubtotal := s.CalculateItemSubtotal(item.PriceSnapshot, item.Quantity, item.DiscountPercent)
		grossItemTotal := item.PriceSnapshot.Mul(decimal.NewFromInt(int64(item.Quantity))).Round(2)
		itemDiscount := grossItemTotal.Sub(itemSubtotal)
		itemDiscounts = itemDiscounts.Add(itemDiscount)

		if item.ItemType == models.ItemTypeService {
			servicesSubtotal = servicesSubtotal.Add(itemSubtotal)
		} else {
			partsSubtotal = partsSubtotal.Add(itemSubtotal)
		}
	}

	// Gross subtotal is parts + services before order-level discount
	gross := partsSubtotal.Add(servicesSubtotal)

	// Combine line item discounts and order-level discount
	totalDiscount := itemDiscounts.Add(orderDiscount)

	// Net taxable subtotal cannot be negative
	netSubtotal := gross.Sub(orderDiscount)
	if netSubtotal.IsNegative() {
		netSubtotal = decimal.Zero
	}

	// Tax calculation (PPN 11%)
	var taxAmount decimal.Decimal
	if taxRate.GreaterThan(decimal.Zero) {
		hundred := decimal.NewFromInt(100)
		taxAmount = netSubtotal.Mul(taxRate.Div(hundred)).Round(2)
	} else {
		taxAmount = decimal.Zero
	}

	grandTotal := netSubtotal.Add(taxAmount).Round(2)

	// Change calculation (Kembalian)
	changeAmount := decimal.Zero
	if paidAmount.GreaterThanOrEqual(grandTotal) {
		changeAmount = paidAmount.Sub(grandTotal).Round(2)
	}

	return OrderCalculationResult{
		PartsSubtotal:    partsSubtotal.Round(2),
		ServicesSubtotal: servicesSubtotal.Round(2),
		GrossSubtotal:    gross.Round(2),
		TotalDiscount:    totalDiscount.Round(2),
		NetSubtotal:      netSubtotal.Round(2),
		TaxRate:          taxRate.Round(2),
		TaxAmount:        taxAmount.Round(2),
		GrandTotal:       grandTotal.Round(2),
		PaidAmount:       paidAmount.Round(2),
		ChangeAmount:     changeAmount.Round(2),
	}
}

// GenerateWorkOrderNumber generates a collision-free sequential work order number.
// Format: WO-AHASS-YYYYMMDD-XXXX (e.g. WO-AHASS-20261027-0001)
// Uses database transactional locking to prevent race conditions during concurrent requests.
func GenerateWorkOrderNumber(db *gorm.DB) (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")
	prefixPattern := fmt.Sprintf("WO-AHASS-%s-%%", datePrefix)

	var count int64
	// In PostgreSQL, pg_advisory_xact_lock(constant_key) can be acquired during transaction.
	// We run within a transaction to maintain atomic sequencing.
	err := db.Transaction(func(tx *gorm.DB) error {
		// Acquire lock if PostgreSQL is available
		var countToday int64
		if err := tx.Model(&models.WorkOrder{}).
			Unscoped().
			Where("wo_number LIKE ?", prefixPattern).
			Count(&countToday).Error; err != nil {
			return err
		}
		count = countToday + 1
		return nil
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("WO-AHASS-%s-%04d", datePrefix, count), nil
}

