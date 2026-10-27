package services

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/models"
)

var (
	PPNRate = decimal.NewFromFloat(0.11) // 11% PPN Indonesia
)

// CalculateRetailItemSubtotal menghitung subtotal item dengan presisi desimal
func CalculateRetailItemSubtotal(price decimal.Decimal, qty int, discountPercent decimal.Decimal) decimal.Decimal {
	if qty <= 0 {
		return decimal.Zero
	}

	q := decimal.NewFromInt(int64(qty))
	gross := price.Mul(q)

	if discountPercent.GreaterThan(decimal.Zero) {
		pct := discountPercent.Div(decimal.NewFromInt(100))
		discountVal := gross.Mul(pct)
		net := gross.Sub(discountVal)
		return net.Round(2)
	}

	return gross.Round(2)
}

// CalculateRetailTotals menghitung total belanja kasir OTC, PPN 11%, dan kembalian
func CalculateRetailTotals(items []models.RetailSaleItem, promoDiscount decimal.Decimal, paidAmount decimal.Decimal) (subtotal, totalDiscount, taxPPN, grandTotal, change decimal.Decimal) {
	subtotal = decimal.Zero
	for _, it := range items {
		subtotal = subtotal.Add(it.Subtotal)
	}

	totalDiscount = promoDiscount
	if totalDiscount.GreaterThan(subtotal) {
		totalDiscount = subtotal
	}

	taxable := subtotal.Sub(totalDiscount)
	taxPPN = taxable.Mul(PPNRate).Round(2)
	grandTotal = taxable.Add(taxPPN).Round(2)

	if paidAmount.GreaterThanOrEqual(grandTotal) {
		change = paidAmount.Sub(grandTotal).Round(2)
	} else {
		change = decimal.Zero
	}

	return subtotal, totalDiscount, taxPPN, grandTotal, change
}

// GenerateRetailInvoiceNumber menghasilkan nomor faktur kasir OTC anti-race condition (INV-NSS-YYYYMMDD-XXXX)
func GenerateRetailInvoiceNumber(tx *gorm.DB) (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")

	// Advisory Lock ID hash untuk retail invoice
	lockKey := int64(88991122)
	isPostgres := tx.Dialector.Name() == "postgres"

	if isPostgres {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			return "", fmt.Errorf("failed to obtain advisory lock for retail invoice: %w", err)
		}
	}

	var count int64
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if err := tx.Model(&models.RetailSale{}).
		Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
		Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to count daily retail sales: %w", err)
	}

	sequence := count + 1
	invoiceNumber := fmt.Sprintf("INV-NSS-%s-%04d", datePrefix, sequence)
	return invoiceNumber, nil
}

