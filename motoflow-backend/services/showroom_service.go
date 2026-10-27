package services

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/models"
)

type CreditSimulationResult struct {
	OTRPrice           decimal.Decimal `json:"otr_price"`
	DPPercent          decimal.Decimal `json:"dp_percent"`
	DPAmount           decimal.Decimal `json:"dp_amount"`
	PrincipalLoan      decimal.Decimal `json:"principal_loan"` // Pokok Pinjaman (OTR - DP)
	TenorMonths        int             `json:"tenor_months"`
	AnnualInterestRate decimal.Decimal `json:"annual_interest_rate"` // e.g. 18.00%
	TotalInterest      decimal.Decimal `json:"total_interest"`
	TotalLoanPayable   decimal.Decimal `json:"total_loan_payable"`
	MonthlyInstallment decimal.Decimal `json:"monthly_installment"`
	AdminFee           decimal.Decimal `json:"admin_fee"`
	InsuranceTLO       decimal.Decimal `json:"insurance_tlo"`
	TotalFirstPayment  decimal.Decimal `json:"total_first_payment"` // DP + Admin + Asuransi
}

// CalculateCreditSimulation menghitung simulasi kredit pembiayaan motor Honda NSC Finance
func CalculateCreditSimulation(otrPrice decimal.Decimal, dpPercent decimal.Decimal, tenorMonths int) CreditSimulationResult {
	if dpPercent.LessThan(decimal.NewFromInt(10)) {
		dpPercent = decimal.NewFromInt(10) // Minimal DP 10% regulasi OJK
	}
	if tenorMonths <= 0 {
		tenorMonths = 35 // Default tenor terpopuler
	}

	dpAmount := otrPrice.Mul(dpPercent.Div(decimal.NewFromInt(100))).Round(2)
	principal := otrPrice.Sub(dpAmount).Round(2)

	// Bunga flat tahunan standar pembiayaan motor NSC: 18.00% per tahun (1.5% per bulan)
	annualRate := decimal.NewFromFloat(18.00)
	monthlyRate := annualRate.Div(decimal.NewFromInt(12)).Div(decimal.NewFromInt(100))

	// Total Bunga = Pokok Pinjaman * (Monthly Rate * Tenor Bulan)
	totalInterest := principal.Mul(monthlyRate.Mul(decimal.NewFromInt(int64(tenorMonths)))).Round(2)
	totalPayable := principal.Add(totalInterest).Round(2)

	installment := totalPayable.Div(decimal.NewFromInt(int64(tenorMonths))).Round(0)

	adminFee := decimal.NewFromInt(350000) // Biaya fidusia & administrasi NSS/NSC
	insuranceTLO := otrPrice.Mul(decimal.NewFromFloat(0.015)).Round(2) // 1.5% TLO Asuransi Kehilangan

	firstPayment := dpAmount.Add(adminFee).Add(insuranceTLO).Round(2)

	return CreditSimulationResult{
		OTRPrice:           otrPrice,
		DPPercent:          dpPercent,
		DPAmount:           dpAmount,
		PrincipalLoan:      principal,
		TenorMonths:        tenorMonths,
		AnnualInterestRate: annualRate,
		TotalInterest:      totalInterest,
		TotalLoanPayable:   totalPayable,
		MonthlyInstallment: installment,
		AdminFee:           adminFee,
		InsuranceTLO:       insuranceTLO,
		TotalFirstPayment:  firstPayment,
	}
}

// GenerateSPKNumber menghasilkan nomor Surat Pemesanan Kendaraan (SPK-NSS-YYYYMMDD-XXXX)
func GenerateSPKNumber(tx *gorm.DB) (string, error) {
	now := time.Now()
	datePrefix := now.Format("20060102")

	lockKey := int64(88992233)
	isPostgres := tx.Dialector.Name() == "postgres"

	if isPostgres {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(?)", lockKey).Error; err != nil {
			return "", fmt.Errorf("failed to obtain advisory lock for SPK: %w", err)
		}
	}

	var count int64
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if err := tx.Model(&models.SPKUnit{}).
		Where("created_at >= ? AND created_at < ?", todayStart, todayEnd).
		Count(&count).Error; err != nil {
		return "", fmt.Errorf("failed to count daily SPKs: %w", err)
	}

	sequence := count + 1
	spkNumber := fmt.Sprintf("SPK-NSS-%s-%04d", datePrefix, sequence)
	return spkNumber, nil
}

