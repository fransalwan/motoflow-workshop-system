package services

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestCalculateCreditSimulation(t *testing.T) {
	// Case: Honda Vario 160 ABS OTR Rp 30.000.000
	// DP 20% = Rp 6.000.000
	// Principal = Rp 24.000.000
	// Tenor 35 Bulan, Annual interest 18% (1.5% per month)
	// Total Interest = 24.000.000 * (0.015 * 35) = 24.000.000 * 0.525 = Rp 12.600.000
	// Total Payable = 24.000.000 + 12.600.000 = Rp 36.600.000
	// Monthly Installment = 36.600.000 / 35 = 1.045.714
	otr := decimal.NewFromInt(30000000)
	dpPct := decimal.NewFromInt(20)
	tenor := 35

	res := CalculateCreditSimulation(otr, dpPct, tenor)

	if !res.DPAmount.Equal(decimal.NewFromInt(6000000)) {
		t.Errorf("DP amount: got %s, want 6000000", res.DPAmount.String())
	}
	if !res.PrincipalLoan.Equal(decimal.NewFromInt(24000000)) {
		t.Errorf("Principal loan: got %s, want 24000000", res.PrincipalLoan.String())
	}
	if !res.TotalInterest.Equal(decimal.NewFromInt(12600000)) {
		t.Errorf("Total interest: got %s, want 12600000", res.TotalInterest.String())
	}
	if !res.TotalLoanPayable.Equal(decimal.NewFromInt(36600000)) {
		t.Errorf("Total loan payable: got %s, want 36600000", res.TotalLoanPayable.String())
	}
	if !res.MonthlyInstallment.Equal(decimal.NewFromInt(1045714)) {
		t.Errorf("Monthly installment: got %s, want 1045714", res.MonthlyInstallment.String())
	}
}
