package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
	"motoflow-backend/services"
)

type RetailController struct {
	db *gorm.DB
}

func NewRetailController(db *gorm.DB) *RetailController {
	return &RetailController{db: db}
}

// Checkout menangani transaksi kasir Over-The-Counter (OTC) langsung dengan transaksi ACID & pemotongan stok
func (ctrl *RetailController) Checkout(c *gin.Context) {
	var req dto.RetailCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload checkout retail tidak valid", "error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		userID, exists = c.Get("user_id")
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Autentikasi kasir diperlukan"})
		return
	}
	cashierID := userID.(uint)

	tx := ctrl.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Validasi promo code jika diinput
	promoDiscount := decimal.Zero
	promoCodeClean := strings.TrimSpace(req.PromoCode)
	if promoCodeClean != "" {
		var promo models.Promotion
		if err := tx.Where("promo_code = ? AND is_active = true AND valid_until >= ?", promoCodeClean, time.Now()).First(&promo).Error; err == nil {
			if promo.PromoType == "percentage" {
				// Akan dihitung terhadap subtotal
			} else {
				promoDiscount = promo.DiscountValue
			}
		}
	}

	// 2. Siapkan items dan potong stok suku cadang
	var saleItems []models.RetailSaleItem
	rawSubtotal := decimal.Zero

	for _, itemReq := range req.Items {
		var part models.Sparepart
		if err := tx.First(&part, itemReq.SparepartID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": fmt.Sprintf("Suku cadang ID %d tidak ditemukan", itemReq.SparepartID)})
			return
		}

		if part.Stock < itemReq.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Stok %s (%s) tidak mencukupi! Tersedia: %d, Diminta: %d", part.Name, part.PartNumber, part.Stock, itemReq.Quantity),
			})
			return
		}

		// Potong stok fisik gudang
		part.Stock -= itemReq.Quantity
		if err := tx.Save(&part).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal memperbarui stok suku cadang"})
			return
		}

		itemSubtotal := services.CalculateRetailItemSubtotal(part.SellingPrice, itemReq.Quantity, itemReq.DiscountPercent)
		rawSubtotal = rawSubtotal.Add(itemSubtotal)

		saleItems = append(saleItems, models.RetailSaleItem{
			SparepartID:        part.ID,
			PartNumberSnapshot: part.PartNumber,
			PartNameSnapshot:   part.Name,
			Quantity:           itemReq.Quantity,
			PriceSnapshot:      part.SellingPrice,
			DiscountPercent:    itemReq.DiscountPercent,
			Subtotal:           itemSubtotal,
		})
	}

	// Hitung diskon promo jika tipe persentase
	if promoCodeClean != "" && promoDiscount.IsZero() {
		var promo models.Promotion
		if err := tx.Where("promo_code = ?", promoCodeClean).First(&promo).Error; err == nil && promo.PromoType == "percentage" {
			promoDiscount = rawSubtotal.Mul(promo.DiscountValue.Div(decimal.NewFromInt(100))).Round(2)
		}
	}

	// 3. Hitung kalkulasi desimal total belanja, PPN 11%, dan kembalian
	subtotal, totalDiscount, taxPPN, grandTotal, change := services.CalculateRetailTotals(saleItems, promoDiscount, req.PaidAmount)

	if req.PaidAmount.LessThan(grandTotal) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Jumlah pembayaran kurang! Total tagihan: %s, Dibayarkan: %s", grandTotal.String(), req.PaidAmount.String()),
		})
		return
	}

	// 4. Generate Nomor Faktur Retail Anti-Race Condition
	invoiceNo, err := services.GenerateRetailInvoiceNumber(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat nomor faktur retail", "error": err.Error()})
		return
	}

	customerName := "Konsumen Walk-in"
	if strings.TrimSpace(req.CustomerName) != "" {
		customerName = strings.TrimSpace(req.CustomerName)
	}

	sale := models.RetailSale{
		InvoiceNumber:  invoiceNo,
		CustomerID:     req.CustomerID,
		CustomerName:   customerName,
		CustomerPhone:  req.CustomerPhone,
		CashierID:      cashierID,
		Subtotal:       subtotal,
		DiscountAmount: totalDiscount,
		TaxPPN:         taxPPN,
		GrandTotal:     grandTotal,
		PaidAmount:     req.PaidAmount,
		ChangeAmount:   change,
		PaymentMethod:  models.PaymentMethod(req.PaymentMethod),
		PaymentStatus:  models.PaymentStatusPaid,
		PromoCode:      promoCodeClean,
		Notes:          req.Notes,
		Items:          saleItems,
	}

	if err := tx.Create(&sale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan transaksi retail", "error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal commit transaksi database"})
		return
	}

	// Invalidate dashboard metrics cache
	services.GetCache().InvalidatePrefix("dashboard_")

	// Preload full response
	ctrl.db.Preload("Items.Sparepart").Preload("Cashier").First(&sale, sale.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": fmt.Sprintf("Transaksi Kasir Retail %s berhasil diselesaikan", sale.InvoiceNumber),
		"data":    sale,
	})
}

// List mengambil daftar riwayat transaksi kasir retail
func (ctrl *RetailController) List(c *gin.Context) {
	search := c.Query("search")
	var sales []models.RetailSale

	query := ctrl.db.Model(&models.RetailSale{}).
		Preload("Items.Sparepart").
		Preload("Cashier").
		Preload("Customer").
		Order("created_at DESC")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("invoice_number LIKE ? OR customer_name LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Limit(50).Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil riwayat transaksi retail", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sales})
}

// GetByID mengambil detail transaksi kasir retail
func (ctrl *RetailController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var sale models.RetailSale

	if err := ctrl.db.Preload("Items.Sparepart").Preload("Cashier").Preload("Customer").First(&sale, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Transaksi retail tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sale})
}

// GetThermalReceipt menghasilkan struk kasir thermal 58mm/80mm ESC/POS text untuk retail OTC
func (ctrl *RetailController) GetThermalReceipt(c *gin.Context) {
	id := c.Param("id")
	widthStr := c.DefaultQuery("width", "58")
	width := 58
	if widthStr == "80" {
		width = 80
	}

	var sale models.RetailSale
	if err := ctrl.db.Preload("Items").Preload("Cashier").First(&sale, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Transaksi tidak ditemukan"})
		return
	}

	lineLength := 32 // 58mm standard width
	if width == 80 {
		lineLength = 48
	}

	center := func(s string) string {
		if len(s) >= lineLength {
			return s[:lineLength] + "\n"
		}
		spaces := (lineLength - len(s)) / 2
		return strings.Repeat(" ", spaces) + s + "\n"
	}

	row := func(left, right string) string {
		maxLeft := lineLength - len(right) - 1
		if maxLeft < 0 {
			maxLeft = 0
		}
		if len(left) > maxLeft {
			left = left[:maxLeft]
		}
		spacing := lineLength - len(left) - len(right)
		if spacing < 1 {
			spacing = 1
		}
		return left + strings.Repeat(" ", spacing) + right + "\n"
	}

	separator := strings.Repeat("-", lineLength) + "\n"
	doubleSeparator := strings.Repeat("=", lineLength) + "\n"

	var sb strings.Builder
	sb.WriteString(doubleSeparator)
	sb.WriteString(center("AHASS NUSANTARA SAKTI"))
	sb.WriteString(center("COUNTER RETAIL HGP RESMI"))
	sb.WriteString(center("Jl. Cempaka Putih Raya No. 128"))
	sb.WriteString(center("Telp: (021) 420-5678"))
	sb.WriteString(doubleSeparator)

	sb.WriteString(row("No. Faktur:", sale.InvoiceNumber))
	sb.WriteString(row("Tanggal:", sale.CreatedAt.Format("02/01/2006 15:04")))
	cashierName := "Kasir AHASS"
	if sale.Cashier != nil {
		cashierName = sale.Cashier.FullName
	}
	sb.WriteString(row("Kasir:", cashierName))
	sb.WriteString(row("Pelanggan:", sale.CustomerName))
	sb.WriteString(separator)

	sb.WriteString(row("ITEM SUKU CADANG", "SUBTOTAL"))
	sb.WriteString(separator)

	for _, it := range sale.Items {
		name := it.PartNameSnapshot
		sb.WriteString(name + "\n")
		qtyPrice := fmt.Sprintf(" %dx Rp %s", it.Quantity, it.PriceSnapshot.StringFixed(0))
		sub := "Rp " + it.Subtotal.StringFixed(0)
		sb.WriteString(row(qtyPrice, sub))
	}

	sb.WriteString(separator)
	sb.WriteString(row("Subtotal Part:", "Rp "+sale.Subtotal.StringFixed(0)))
	if sale.DiscountAmount.GreaterThan(decimal.Zero) {
		sb.WriteString(row("Diskon Promo:", "-Rp "+sale.DiscountAmount.StringFixed(0)))
	}
	sb.WriteString(row("PPN (11%):", "Rp "+sale.TaxPPN.StringFixed(0)))
	sb.WriteString(doubleSeparator)
	sb.WriteString(row("TOTAL BAYAR:", "Rp "+sale.GrandTotal.StringFixed(0)))
	sb.WriteString(row("Metode Bayar:", strings.ToUpper(string(sale.PaymentMethod))))
	sb.WriteString(row("Diterima:", "Rp "+sale.PaidAmount.StringFixed(0)))
	sb.WriteString(row("Kembalian:", "Rp "+sale.ChangeAmount.StringFixed(0)))
	sb.WriteString(doubleSeparator)
	sb.WriteString(center("100% HONDA GENUINE PARTS"))
	sb.WriteString(center("Jaminan Kualitas Astra Honda"))
	sb.WriteString(center("Terima Kasih Atas Kunjungan Anda"))
	sb.WriteString(center("PT NUSANTARA SAKTI GROUP"))
	sb.WriteString(doubleSeparator)
	sb.WriteString("\n")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"receipt_text": sb.String(),
		"sale": sale,
	})
}

// GetPromotions mengambil promo & bundling aktif
func (ctrl *RetailController) GetPromotions(c *gin.Context) {
	var promos []models.Promotion
	if err := ctrl.db.Where("is_active = true AND valid_until >= ?", time.Now()).Find(&promos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil daftar promo", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": promos})
}

func formatRupiah(d decimal.Decimal) string {
	val := d.IntPart()
	str := fmt.Sprintf("%d", val)
	n := len(str)
	if n <= 3 {
		return "Rp " + str
	}
	var res []byte
	rem := n % 3
	if rem > 0 {
		res = append(res, str[:rem]...)
		if n > rem {
			res = append(res, '.')
		}
	}
	for i := rem; i < n; i += 3 {
		res = append(res, str[i:i+3]...)
		if i+3 < n {
			res = append(res, '.')
		}
	}
	return "Rp " + string(res)
}
