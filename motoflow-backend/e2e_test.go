package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/config"
	"motoflow-backend/dto"
	"motoflow-backend/models"
	"motoflow-backend/routes"
	"motoflow-backend/seeders"
)

var (
	testRouter *gin.Engine
	testDB     *gorm.DB
)

func setupTestEnvironment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.LoadConfig()

	// Use sqlite file or in-memory
	cfg.DBDriver = "sqlite"
	cfg.DBPath = fmt.Sprintf("file:test_%d?mode=memory&cache=shared", time.Now().UnixNano())

	var err error
	testDB, err = config.InitDB(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize test DB: %v", err)
	}

	if err := seeders.SeedAll(testDB); err != nil {
		t.Fatalf("Failed to seed initial test data: %v", err)
	}
	seeders.SeedRetailAndShowroom(testDB)

	testRouter = routes.SetupRouter(testDB)
}

func loginHelper(t *testing.T, username, password string) string {
	payload := dto.LoginRequest{
		Username: username,
		Password: password,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Login failed for %s: code %d, body: %s", username, w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}
	return resp.Data.Token
}

func performRequest(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestEndToEndAHASSSystem(t *testing.T) {
	setupTestEnvironment(t)

	// =========================================================================
	// 1. UJI OTENTIKASI & ROLE-BASED ACCESS CONTROL (RBAC)
	// =========================================================================
	t.Run("Auth: Multi-Role Login & Security", func(t *testing.T) {
		// Valid logins for all NSS AHASS roles
		adminToken := loginHelper(t, "admin", "password123")
		advisorToken := loginHelper(t, "advisor", "password123")
		kasirToken := loginHelper(t, "kasir", "password123")
		mekanikToken := loginHelper(t, "mekanik", "password123")

		if adminToken == "" || advisorToken == "" || kasirToken == "" || mekanikToken == "" {
			t.Fatal("Expected non-empty JWT tokens for all AHASS roles")
		}

		// Invalid credentials test
		badReq := performRequest(testRouter, "POST", "/api/v1/auth/login", dto.LoginRequest{
			Username: "admin",
			Password: "wrongpassword",
		}, "")
		if badReq.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP 401 for bad password, got %d", badReq.Code)
		}

		// Unauthorized access to protected route
		unauthReq := performRequest(testRouter, "GET", "/api/v1/dashboard/summary", nil, "")
		if unauthReq.Code != http.StatusUnauthorized {
			t.Errorf("Expected HTTP 401 for missing token, got %d", unauthReq.Code)
		}

		// Verify profile for advisor
		profReq := performRequest(testRouter, "GET", "/api/v1/auth/profile", nil, advisorToken)
		if profReq.Code != http.StatusOK {
			t.Errorf("Expected HTTP 200 for profile, got %d", profReq.Code)
		}
	})

	advisorToken := loginHelper(t, "advisor", "password123")
	kasirToken := loginHelper(t, "kasir", "password123")
	mekanikToken := loginHelper(t, "mekanik", "password123")

	// =========================================================================
	// 2. UJI MASTER DATA PELANGGAN & SUKU CADANG HGP
	// =========================================================================
	var newCustomerID uint
	t.Run("Master Data: Register Customer & Verify HGP Inventory", func(t *testing.T) {
		custPayload := dto.CreateCustomerRequest{
			Name:         "Bambang Soedirman",
			Phone:        "081299887766",
			PlateNumber:  "B 8888 NSS",
			MotorModel:   "Honda PCX 160 CBS",
			Year:         2025,
			CurrentKM:    8500,
			Address:      "Jl. Daan Mogot No. 12, Jakarta Barat",
		}

		w := performRequest(testRouter, "POST", "/api/v1/customers", custPayload, advisorToken)
		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create customer: code %d, body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool            `json:"success"`
			Data    models.Customer `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Data.ID == 0 || resp.Data.PlateNumber != "B 8888 NSS" {
			t.Fatalf("Customer data mismatch: %+v", resp.Data)
		}
		newCustomerID = resp.Data.ID

		// Verify Spareparts list
		spW := performRequest(testRouter, "GET", "/api/v1/spareparts", nil, advisorToken)
		if spW.Code != http.StatusOK {
			t.Fatalf("Failed to list spareparts: code %d", spW.Code)
		}
	})

	// =========================================================================
	// 3. UJI WORK ORDER SERVIS MOTOR (ALUR 5 TAHAPAN AHASS)
	// =========================================================================
	var createdWO models.WorkOrder
	t.Run("Service WO: 5-Stage Life-Cycle & Cashier Settlement", func(t *testing.T) {
		// Dapatkan ID sparepart Oli MPX2 & Jasa Servis Lengkap
		var oliPart models.Sparepart
		var jasaServis models.Sparepart
		testDB.Where("part_number = ?", "08232-2MB-K0NZ1").First(&oliPart) // AHM Oil MPX2 (58.000)
		testDB.Where("part_number = ?", "SRV-LENGKAP-MATIC").First(&jasaServis)  // Servis Lengkap Matic (120.000)
		initialOliStock := oliPart.Stock

		targetCustomerID := newCustomerID
		if targetCustomerID == 0 {
			var c models.Customer
			testDB.First(&c)
			targetCustomerID = c.ID
		}

		// Tahap 1: Pendaftaran Servis oleh Service Advisor (Draft)
		woReq := dto.CreateWorkOrderRequest{
			CustomerID:       targetCustomerID,
			EntryKM:          8500,
			Complaints:       "Tarikan motor berat, jadwal ganti oli rutin dan servis berkala 8.000 KM.",
			ServiceDiagnosis: "Servis lengkap sistem injeksi PGM-FI, ganti oli mesin AHM Oil MPX2.",
			Items: []dto.WorkOrderItemInput{
				{
					SparepartID:     &oliPart.ID,
					ItemType:        models.ItemTypePart,
					PartNumber:      oliPart.PartNumber,
					ItemName:        oliPart.Name,
					Price:           oliPart.SellingPrice,
					Quantity:        1,
					DiscountPercent: decimal.Zero,
				},
				{
					SparepartID:     &jasaServis.ID,
					ItemType:        models.ItemTypeService,
					PartNumber:      jasaServis.PartNumber,
					ItemName:        jasaServis.Name,
					Price:           jasaServis.SellingPrice,
					Quantity:        1,
					DiscountPercent: decimal.Zero,
				},
			},
		}

		wCreate := performRequest(testRouter, "POST", "/api/v1/work-orders", woReq, advisorToken)
		if wCreate.Code != http.StatusCreated {
			t.Fatalf("Failed to create WO: code %d, body: %s", wCreate.Code, wCreate.Body.String())
		}

		var createResp struct {
			Success bool             `json:"success"`
			Data    models.WorkOrder `json:"data"`
		}
		json.Unmarshal(wCreate.Body.Bytes(), &createResp)
		createdWO = createResp.Data

		if !strings.HasPrefix(createdWO.WONumber, "WO-AHASS-") {
			t.Errorf("Expected WO number format WO-AHASS-..., got %s", createdWO.WONumber)
		}
		if createdWO.Status != models.WOStatusDraft {
			t.Errorf("Expected initial status draft, got %s", createdWO.Status)
		}

		// Verifikasi kalkulasi desimal: Parts (58.000) + Service (120.000) = 178.000. PPN 11% = 19.580. Grand Total = 197.580
		expectedGrand := decimal.NewFromInt(197580)
		if !createdWO.GrandTotal.Equal(expectedGrand) {
			t.Errorf("Grand total mismatch: got %s, want %s", createdWO.GrandTotal.String(), expectedGrand.String())
		}

		// Tahap 2: Mekanik Mulai Kerjakan di Pit (In Progress)
		wProg := performRequest(testRouter, "PATCH", fmt.Sprintf("/api/v1/work-orders/%d/status", createdWO.ID), dto.UpdateWorkOrderStatusRequest{
			Status: models.WOStatusInProgress,
		}, mekanikToken)
		if wProg.Code != http.StatusOK {
			t.Fatalf("Failed to set WO in_progress: code %d, body: %s", wProg.Code, wProg.Body.String())
		}

		// Tahap 3: Mekanik Selesai & Inspeksi Akhir (Completed)
		wComp := performRequest(testRouter, "PATCH", fmt.Sprintf("/api/v1/work-orders/%d/status", createdWO.ID), dto.UpdateWorkOrderStatusRequest{
			Status: models.WOStatusCompleted,
		}, mekanikToken)
		if wComp.Code != http.StatusOK {
			t.Fatalf("Failed to set WO completed: code %d, body: %s", wComp.Code, wComp.Body.String())
		}

		// Tahap 4: Kasir POS Pembayaran Billing Lunas & Potong Stok Suku Cadang
		payReq := dto.ProcessPaymentRequest{
			PaymentMethod: models.PaymentCash,
			PaidAmount:    decimal.NewFromInt(200000), // Bayar 200rb
		}
		wPay := performRequest(testRouter, "POST", fmt.Sprintf("/api/v1/work-orders/%d/pay", createdWO.ID), payReq, kasirToken)
		if wPay.Code != http.StatusOK {
			t.Fatalf("Failed to pay WO: code %d, body: %s", wPay.Code, wPay.Body.String())
		}

		var payResp struct {
			Success bool             `json:"success"`
			Data    models.WorkOrder `json:"data"`
		}
		json.Unmarshal(wPay.Body.Bytes(), &payResp)
		if payResp.Data.Status != models.WOStatusPaid || payResp.Data.PaymentStatus != models.PaymentStatusPaid {
			t.Errorf("Expected paid status, got %s / %s", payResp.Data.Status, payResp.Data.PaymentStatus)
		}
		// Kembalian: 200.000 - 197.580 = 2.420
		expectedChange := decimal.NewFromInt(2420)
		if !payResp.Data.ChangeAmount.Equal(expectedChange) {
			t.Errorf("Change amount mismatch: got %s, want %s", payResp.Data.ChangeAmount.String(), expectedChange.String())
		}

		// Verifikasi stok oli berkurang tepat 1 botol
		var updatedOli models.Sparepart
		testDB.First(&updatedOli, oliPart.ID)
		if updatedOli.Stock != initialOliStock-1 {
			t.Errorf("Stock deduction failure: initial %d, current %d, expected %d", initialOliStock, updatedOli.Stock, initialOliStock-1)
		}

		// Tahap 5: Cetak Struk Kasir Thermal (58mm/80mm ESC/POS)
		wRcpt := performRequest(testRouter, "GET", fmt.Sprintf("/api/v1/work-orders/%d/receipt", createdWO.ID), nil, kasirToken)
		if wRcpt.Code != http.StatusOK {
			t.Errorf("Failed to generate thermal receipt: code %d", wRcpt.Code)
		}
		var receiptResp struct {
			Success bool `json:"success"`
			Data    struct {
				RawText string `json:"raw_text"`
			} `json:"data"`
		}
		json.Unmarshal(wRcpt.Body.Bytes(), &receiptResp)
		if !strings.Contains(receiptResp.Data.RawText, "NUSANTARA SAKTI") {
			t.Errorf("Receipt text missing AHASS NSS header: %s", receiptResp.Data.RawText)
		}
		if !strings.Contains(receiptResp.Data.RawText, "GARANSI SERVIS AHASS 500 KM") {
			t.Errorf("Receipt text missing warranty guarantee: %s", receiptResp.Data.RawText)
		}

		// Tahap 6: Cetak Dokumen Faktur Resmi A4 (GoFPDF Engine)
		wPDF := performRequest(testRouter, "GET", fmt.Sprintf("/api/v1/work-orders/%d/invoice-pdf", createdWO.ID), nil, kasirToken)
		if wPDF.Code != http.StatusOK {
			t.Errorf("Failed to generate A4 PDF: code %d", wPDF.Code)
		}
		if wPDF.Header().Get("Content-Type") != "application/pdf" {
			t.Errorf("Expected application/pdf, got %s", wPDF.Header().Get("Content-Type"))
		}
		pdfBytes := wPDF.Body.Bytes()
		if len(pdfBytes) < 1000 || !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
			t.Errorf("Invalid PDF binary output stream: size %d bytes", len(pdfBytes))
		}
	})

	// =========================================================================
	// 4. UJI KASIR RETAIL OVER-THE-COUNTER (OTC) DENGAN PROMO VOUCHER
	// =========================================================================
	t.Run("Retail OTC: Counter Checkout with Promo Voucher & Barcode Simulation", func(t *testing.T) {
		// Cek daftar promo aktif
		wPromo := performRequest(testRouter, "GET", "/api/v1/retail/promotions", nil, kasirToken)
		if wPromo.Code != http.StatusOK {
			t.Fatalf("Failed to fetch promotions: code %d", wPromo.Code)
		}

		// Busi NGK Part
		var busiPart models.Sparepart
		testDB.Where("part_number = ?", "31919-K25-601").First(&busiPart)
		initialBusiStock := busiPart.Stock

		// Kasir OTC Checkout: 2x Busi @ 25.000 = 50.000, Gunakan promo NSSHEMAT10 (10% = -5.000)
		// Taxable: 45.000. PPN 11% = 4.950. Grand Total = 49.950. Bayar 50.000, Kembalian = 50
		otcReq := dto.RetailCheckoutRequest{
			CustomerName:  "Pak Hendra (Walk-in)",
			CustomerPhone: "081912345678",
			PaymentMethod: "cash",
			PaidAmount:    decimal.NewFromInt(50000),
			PromoCode:     "NSSHEMAT10",
			Notes:         "Pembelian busi langsung di counter HGP",
			Items: []dto.RetailItemRequest{
				{
					SparepartID:     busiPart.ID,
					Quantity:        2,
					DiscountPercent: decimal.Zero,
				},
			},
		}

		wCheckout := performRequest(testRouter, "POST", "/api/v1/retail/checkout", otcReq, kasirToken)
		if wCheckout.Code != http.StatusCreated {
			t.Fatalf("Failed to execute retail checkout: code %d, body: %s", wCheckout.Code, wCheckout.Body.String())
		}

		var retailResp struct {
			Success bool              `json:"success"`
			Data    models.RetailSale `json:"data"`
		}
		json.Unmarshal(wCheckout.Body.Bytes(), &retailResp)
		sale := retailResp.Data

		if !strings.HasPrefix(sale.InvoiceNumber, "INV-NSS-") {
			t.Errorf("Expected retail invoice format INV-NSS-..., got %s", sale.InvoiceNumber)
		}
		if !sale.DiscountAmount.Equal(decimal.NewFromInt(5000)) {
			t.Errorf("Promo discount mismatch: got %s, want 5000", sale.DiscountAmount.String())
		}
		expectedGrand := decimal.NewFromInt(49950)
		if !sale.GrandTotal.Equal(expectedGrand) {
			t.Errorf("Grand total mismatch: got %s, want %s", sale.GrandTotal.String(), expectedGrand.String())
		}

		// Verifikasi stok busi berkurang tepat 2 pcs
		var updatedBusi models.Sparepart
		testDB.First(&updatedBusi, busiPart.ID)
		if updatedBusi.Stock != initialBusiStock-2 {
			t.Errorf("Busi stock deduction failure: initial %d, current %d, expected %d", initialBusiStock, updatedBusi.Stock, initialBusiStock-2)
		}

		// Verifikasi struk kasir retail OTC
		wRcpt := performRequest(testRouter, "GET", fmt.Sprintf("/api/v1/retail/sales/%d/receipt", sale.ID), nil, kasirToken)
		if wRcpt.Code != http.StatusOK {
			t.Errorf("Failed to get retail receipt: code %d", wRcpt.Code)
		}
		if !strings.Contains(wRcpt.Body.String(), "COUNTER RETAIL HGP RESMI") {
			t.Errorf("Missing retail header in receipt: %s", wRcpt.Body.String())
		}
	})

	// =========================================================================
	// 5. UJI SHOWROOM MOTOR HONDA, SIMULASI KREDIT NSC, & TERBITKAN SPK
	// =========================================================================
	t.Run("Showroom & NSC Finance: Credit Simulation & Unit VIN Booking", func(t *testing.T) {
		// 1. Ambil katalog unit motor
		wUnits := performRequest(testRouter, "GET", "/api/v1/showroom/units?category=Matic", nil, advisorToken)
		if wUnits.Code != http.StatusOK {
			t.Fatalf("Failed to fetch motor units: code %d", wUnits.Code)
		}
		var unitsResp struct {
			Data []models.MotorUnit `json:"data"`
		}
		json.Unmarshal(wUnits.Body.Bytes(), &unitsResp)
		if len(unitsResp.Data) == 0 {
			t.Fatal("Showroom motor units catalog is empty")
		}

		selectedUnit := unitsResp.Data[0]
		if selectedUnit.Status != "ready" {
			t.Fatalf("Expected unit to be ready initially, got %s", selectedUnit.Status)
		}

		// 2. Simulasi Kredit Pembiayaan NSC Finance
		simReq := dto.CreditSimulationRequest{
			OTRPrice:    selectedUnit.OTRPrice,
			DPPercent:   decimal.NewFromInt(20),
			TenorMonths: 35,
		}
		wSim := performRequest(testRouter, "POST", "/api/v1/showroom/simulate-credit", simReq, advisorToken)
		if wSim.Code != http.StatusOK {
			t.Fatalf("Failed to simulate credit: code %d", wSim.Code)
		}
		var simResp struct {
			Data struct {
				PrincipalLoan      decimal.Decimal `json:"principal_loan"`
				MonthlyInstallment decimal.Decimal `json:"monthly_installment"`
				TotalFirstPayment  decimal.Decimal `json:"total_first_payment"`
			} `json:"data"`
		}
		json.Unmarshal(wSim.Body.Bytes(), &simResp)
		if simResp.Data.MonthlyInstallment.LessThanOrEqual(decimal.Zero) {
			t.Errorf("Expected positive monthly installment, got %s", simResp.Data.MonthlyInstallment.String())
		}

		// 3. Terbitkan Surat Pemesanan Kendaraan (SPK) dengan Alokasi VIN
		targetCustomerID := newCustomerID
		if targetCustomerID == 0 {
			var c models.Customer
			testDB.First(&c)
			targetCustomerID = c.ID
		}

		spkReq := dto.CreateSPKRequest{
			CustomerID:         targetCustomerID,
			MotorUnitID:        selectedUnit.ID,
			PurchaseType:       "credit_nsc",
			BookingFee:         decimal.NewFromInt(500000),
			DownPayment:        selectedUnit.OTRPrice.Mul(decimal.NewFromFloat(0.20)).Round(2),
			TenorMonths:        35,
			MonthlyInstallment: simResp.Data.MonthlyInstallment,
			Notes:              "Pemesanan unit Vario 160 ABS warna Grande Matte Black via NSC Finance.",
		}
		wSPK := performRequest(testRouter, "POST", "/api/v1/showroom/spk", spkReq, advisorToken)
		if wSPK.Code != http.StatusCreated {
			t.Fatalf("Failed to create SPK: code %d, body: %s", wSPK.Code, wSPK.Body.String())
		}

		var spkResp struct {
			Success bool           `json:"success"`
			Data    models.SPKUnit `json:"data"`
		}
		json.Unmarshal(wSPK.Body.Bytes(), &spkResp)
		spk := spkResp.Data

		if !strings.HasPrefix(spk.SPKNumber, "SPK-NSS-") {
			t.Errorf("Expected SPK number format SPK-NSS-..., got %s", spk.SPKNumber)
		}

		// 4. Verifikasi status unit motor telah berubah menjadi 'booked' (anti-double booking)
		var checkUnit models.MotorUnit
		testDB.First(&checkUnit, selectedUnit.ID)
		if checkUnit.Status != "booked" {
			t.Errorf("Expected motor unit status booked, got %s", checkUnit.Status)
		}

		// Coba booking ulang unit yang sama -> harus ditolak HTTP 400
		wDuplicate := performRequest(testRouter, "POST", "/api/v1/showroom/spk", spkReq, advisorToken)
		if wDuplicate.Code != http.StatusBadRequest {
			t.Errorf("Expected HTTP 400 for duplicate booking of already booked unit, got %d", wDuplicate.Code)
		}

		// 5. Cetak Dokumen Konfirmasi SPK
		wSPKRcpt := performRequest(testRouter, "GET", fmt.Sprintf("/api/v1/showroom/spk/%d/receipt", spk.ID), nil, advisorToken)
		if wSPKRcpt.Code != http.StatusOK {
			t.Errorf("Failed to get SPK receipt: code %d", wSPKRcpt.Code)
		}
		if !strings.Contains(wSPKRcpt.Body.String(), "SURAT PEMESANAN KENDARAAN (SPK)") {
			t.Errorf("Receipt text missing SPK title: %s", wSPKRcpt.Body.String())
		}
		if !strings.Contains(wSPKRcpt.Body.String(), "NSC FINANCE") {
			t.Errorf("Receipt text missing leasing institution NSC Finance: %s", wSPKRcpt.Body.String())
		}
	})

	// =========================================================================
	// 6. UJI EXECUTIVE DASHBOARD AGGREGATION & METRICS
	// =========================================================================
	t.Run("Executive Dashboard: Metric Summaries & Aggregation Integrity", func(t *testing.T) {
		adminToken := loginHelper(t, "admin", "password123")
		wDash := performRequest(testRouter, "GET", "/api/v1/dashboard/summary", nil, adminToken)
		if wDash.Code != http.StatusOK {
			t.Fatalf("Failed to fetch dashboard summary: code %d", wDash.Code)
		}

		var dashResp struct {
			Success bool `json:"success"`
			Data    struct {
				TodayRevenue       decimal.Decimal `json:"today_revenue"`
				MonthRevenue       decimal.Decimal `json:"month_revenue"`
				MonthBikesServiced int             `json:"month_bikes_serviced"`
				LowStockCount      int             `json:"low_stock_count"`
			} `json:"data"`
		}
		json.Unmarshal(wDash.Body.Bytes(), &dashResp)

		if dashResp.Data.MonthBikesServiced <= 0 {
			t.Errorf("Expected positive bikes serviced count, got %d", dashResp.Data.MonthBikesServiced)
		}
		if dashResp.Data.TodayRevenue.LessThanOrEqual(decimal.Zero) {
			t.Errorf("Expected positive today revenue, got %s", dashResp.Data.TodayRevenue.String())
		}
	})
}
