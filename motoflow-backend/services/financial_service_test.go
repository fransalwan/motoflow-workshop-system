package services

import (
	"testing"

	"github.com/shopspring/decimal"

	"motoflow-backend/models"
)

func TestCalculateItemSubtotal(t *testing.T) {
	service := NewFinancialService()

	tests := []struct {
		name            string
		price           decimal.Decimal
		qty             int
		discountPercent decimal.Decimal
		expectedSubtotal decimal.Decimal
	}{
		{
			name:            "Regular sparepart without discount (Oli MPX2)",
			price:           decimal.NewFromInt(58000),
			qty:             1,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.NewFromInt(58000),
		},
		{
			name:            "Multiple quantities (2x Kampas Rem @ Rp 45.000)",
			price:           decimal.NewFromInt(45000),
			qty:             2,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.NewFromInt(90000),
		},
		{
			name:            "With 10% discount (Paket Servis Rp 150.000)",
			price:           decimal.NewFromInt(150000),
			qty:             1,
			discountPercent: decimal.NewFromInt(10),
			expectedSubtotal: decimal.NewFromInt(135000),
		},
		{
			name:            "Free promotional item (100% discount)",
			price:           decimal.NewFromInt(25000),
			qty:             1,
			discountPercent: decimal.NewFromInt(100),
			expectedSubtotal: decimal.Zero,
		},
		{
			name:            "Zero quantity returns zero",
			price:           decimal.NewFromInt(50000),
			qty:             0,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.Zero,
		},
		{
			name:            "Negative quantity returns zero",
			price:           decimal.NewFromInt(50000),
			qty:             -2,
			discountPercent: decimal.Zero,
			expectedSubtotal: decimal.Zero,
		},
		{
			name:            "Decimal price with fractional discount rounded to 2 decimals",
			price:           decimal.NewFromFloat(54500.50),
			qty:             3,
			discountPercent: decimal.NewFromFloat(7.5), // 163501.50 * 0.925 = 151238.8875 -> 151238.89
			expectedSubtotal: decimal.NewFromFloat(151238.89),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateItemSubtotal(tt.price, tt.qty, tt.discountPercent)
			if !result.Equal(tt.expectedSubtotal) {
				t.Errorf("CalculateItemSubtotal() = %s, want %s", result.String(), tt.expectedSubtotal.String())
			}
		})
	}
}

func TestCalculateWorkOrderTotals(t *testing.T) {
	service := NewFinancialService()

	t.Run("Standard AHASS Service with Parts, Labor, PPN 11%, and Cash Payment", func(t *testing.T) {
		items := []models.WorkOrderItem{
			{
				ItemType:        models.ItemTypeService,
				ItemNameSnapshot: "Jasa Servis Lengkap Vario 160",
				PriceSnapshot:   decimal.NewFromInt(120000),
				Quantity:        1,
				DiscountPercent: decimal.Zero,
			},
			{
				ItemType:        models.ItemTypePart,
				ItemNameSnapshot: "AHM Oil MPX2 0.8L",
				PriceSnapshot:   decimal.NewFromInt(58000),
				Quantity:        1,
				DiscountPercent: decimal.Zero,
			},
			{
				ItemType:        models.ItemTypePart,
				ItemNameSnapshot: "Busi NGK CPR9EA-9",
				PriceSnapshot:   decimal.NewFromInt(25000),
				Quantity:        1,
				DiscountPercent: decimal.Zero,
			},
		}

		orderDiscount := decimal.Zero
		taxRate := decimal.NewFromInt(11) // PPN 11%
		paidAmount := decimal.NewFromInt(250000)

		res := service.CalculateWorkOrderTotals(items, orderDiscount, taxRate, paidAmount)

		// Parts subtotal: 58,000 + 25,000 = 83,000
		expectedParts := decimal.NewFromInt(83000)
		if !res.PartsSubtotal.Equal(expectedParts) {
			t.Errorf("PartsSubtotal = %s, want %s", res.PartsSubtotal, expectedParts)
		}

		// Services subtotal: 120,000
		expectedServices := decimal.NewFromInt(120000)
		if !res.ServicesSubtotal.Equal(expectedServices) {
			t.Errorf("ServicesSubtotal = %s, want %s", res.ServicesSubtotal, expectedServices)
		}

		// Net subtotal: 203,000
		expectedNet := decimal.NewFromInt(203000)
		if !res.NetSubtotal.Equal(expectedNet) {
			t.Errorf("NetSubtotal = %s, want %s", res.NetSubtotal, expectedNet)
		}

		// PPN 11%: 203,000 * 0.11 = 22,330
		expectedTax := decimal.NewFromInt(22330)
		if !res.TaxAmount.Equal(expectedTax) {
			t.Errorf("TaxAmount = %s, want %s", res.TaxAmount, expectedTax)
		}

		// Grand Total: 203,000 + 22,330 = 225,330
		expectedGrand := decimal.NewFromInt(225330)
		if !res.GrandTotal.Equal(expectedGrand) {
			t.Errorf("GrandTotal = %s, want %s", res.GrandTotal, expectedGrand)
		}

		// Change (Kembalian): 250,000 - 225,330 = 24,670
		expectedChange := decimal.NewFromInt(24670)
		if !res.ChangeAmount.Equal(expectedChange) {
			t.Errorf("ChangeAmount = %s, want %s", res.ChangeAmount, expectedChange)
		}
	})

	t.Run("Verification of Floating-Point Drift Resistance", func(t *testing.T) {
		// In IEEE 754 float64: 0.1 + 0.2 = 0.30000000000000004
		// In shopspring/decimal: 0.1 + 0.2 == 0.3 EXACTLY
		d1 := decimal.NewFromFloat(0.1)
		d2 := decimal.NewFromFloat(0.2)
		expected := decimal.NewFromFloat(0.3)

		if !d1.Add(d2).Equal(expected) {
			t.Errorf("Decimal failed floating-point drift test: %s + %s != %s", d1, d2, expected)
		}
	})

	t.Run("Exact Payment and Underpayment Handling", func(t *testing.T) {
		items := []models.WorkOrderItem{
			{
				ItemType:        models.ItemTypeService,
				ItemNameSnapshot: "Jasa Ganti Oli",
				PriceSnapshot:   decimal.NewFromInt(15000),
				Quantity:        1,
				DiscountPercent: decimal.Zero,
			},
		}

		// Zero tax
		resExact := service.CalculateWorkOrderTotals(items, decimal.Zero, decimal.Zero, decimal.NewFromInt(15000))
		if !resExact.ChangeAmount.Equal(decimal.Zero) {
			t.Errorf("Exact payment should give zero change, got %s", resExact.ChangeAmount)
		}

		// Underpaid: paid 10,000 for 15,000 order
		resUnder := service.CalculateWorkOrderTotals(items, decimal.Zero, decimal.Zero, decimal.NewFromInt(10000))
		if !resUnder.ChangeAmount.Equal(decimal.Zero) {
			t.Errorf("Underpayment should give zero change, got %s", resUnder.ChangeAmount)
		}
	})
}

