package services

import (
	"testing"

	"github.com/shopspring/decimal"

	"motoflow-backend/models"
)

func TestCalculateRetailItemSubtotal(t *testing.T) {
	tests := []struct {
		name            string
		price           decimal.Decimal
		qty             int
		discountPercent decimal.Decimal
		expectedSubtotal decimal.Decimal
	}{
		{
			name:            "Direct purchase Oli MPX2 (1x @ 56.000)",
			price:           decimal.NewFromInt(56000),
			qty:             1,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.NewFromInt(56000),
		},
		{
			name:            "Multiple quantities (3x Busi @ 25.000)",
			price:           decimal.NewFromInt(25000),
			qty:             3,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.NewFromInt(75000),
		},
		{
			name:            "With member discount (10% off V-Belt Kit 150.000)",
			price:           decimal.NewFromInt(150000),
			qty:             1,
			discountPercent: decimal.NewFromInt(10),
			expectedSubtotal: decimal.NewFromInt(135000),
		},
		{
			name:            "Zero quantity returns zero",
			price:           decimal.NewFromInt(50000),
			qty:             0,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.Zero,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRetailItemSubtotal(tt.price, tt.qty, tt.discountPercent)
			if !result.Equal(tt.expectedSubtotal) {
				t.Errorf("got %s, expected %s", result.String(), tt.expectedSubtotal.String())
			}
		})
	}
}

func TestCalculateRetailTotals(t *testing.T) {
	// Item 1: Oli MPX2 (Rp 56.000)
	// Item 2: Busi (Rp 25.000)
	// Subtotal = 81.000
	// Promo Voucher = Rp 10.000
	// Taxable = 71.000
	// PPN 11% = 71.000 * 0.11 = 7.810
	// Grand Total = 78.810
	// Paid = Rp 100.000
	// Change = Rp 21.190
	items := []models.RetailSaleItem{
		{Subtotal: decimal.NewFromInt(56000)},
		{Subtotal: decimal.NewFromInt(25000)},
	}
	promoDiscount := decimal.NewFromInt(10000)
	paidAmount := decimal.NewFromInt(100000)

	subtotal, discount, tax, grandTotal, change := CalculateRetailTotals(items, promoDiscount, paidAmount)

	if !subtotal.Equal(decimal.NewFromInt(81000)) {
		t.Errorf("Subtotal mismatch: got %s, want 81000", subtotal.String())
	}
	if !discount.Equal(decimal.NewFromInt(10000)) {
		t.Errorf("Discount mismatch: got %s, want 10000", discount.String())
	}
	if !tax.Equal(decimal.NewFromInt(7810)) {
		t.Errorf("Tax PPN 11%% mismatch: got %s, want 7810", tax.String())
	}
	if !grandTotal.Equal(decimal.NewFromInt(78810)) {
		t.Errorf("Grand total mismatch: got %s, want 78810", grandTotal.String())
	}
	if !change.Equal(decimal.NewFromInt(21190)) {
		t.Errorf("Change amount mismatch: got %s, want 21190", change.String())
	}
}
