package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
	"motoflow-backend/services"
)

type ShowroomController struct {
	db *gorm.DB
}

func NewShowroomController(db *gorm.DB) *ShowroomController {
	return &ShowroomController{db: db}
}

// GetMotorUnits mengambil katalog unit motor di showroom
func (ctrl *ShowroomController) GetMotorUnits(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
	search := c.Query("search")

	query := ctrl.db.Model(&models.MotorUnit{}).Order("model_name ASC, status ASC")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("model_name LIKE ? OR frame_number LIKE ? OR color LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	var units []models.MotorUnit
	if err := query.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil katalog showroom", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": units})
}

// SimulateCredit menghitung simulasi kredit pembiayaan motor Honda NSC Finance
func (ctrl *ShowroomController) SimulateCredit(c *gin.Context) {
	var req dto.CreditSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload simulasi kredit tidak valid", "error": err.Error()})
		return
	}

	res := services.CalculateCreditSimulation(req.OTRPrice, req.DPPercent, req.TenorMonths)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": res})
}

// CreateSPK menerbitkan Surat Pemesanan Kendaraan baru dan mengalokasikan unit motor
func (ctrl *ShowroomController) CreateSPK(c *gin.Context) {
	var req dto.CreateSPKRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Payload penerbitan SPK tidak valid", "error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		userID, exists = c.Get("user_id")
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Autentikasi sales counter diperlukan"})
		return
	}
	salesID := userID.(uint)

	tx := ctrl.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Cek ketersediaan motor unit fisik
	var unit models.MotorUnit
	if err := tx.First(&unit, req.MotorUnitID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Unit motor tidak ditemukan"})
		return
	}

	if unit.Status != "ready" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("Unit motor %s (%s) tidak siap! Status saat ini: %s", unit.ModelName, unit.FrameNumber, unit.Status)})
		return
	}

	// 2. Generate Nomor SPK Anti-Race Condition
	spkNo, err := services.GenerateSPKNumber(tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghasilkan nomor SPK", "error": err.Error()})
		return
	}

	leasingStatus := "not_applicable"
	if req.PurchaseType == "credit_nsc" {
		leasingStatus = "survey"
	}

	spk := models.SPKUnit{
		SPKNumber:          spkNo,
		CustomerID:         req.CustomerID,
		SalesID:            salesID,
		MotorUnitID:        unit.ID,
		PurchaseType:       req.PurchaseType,
		BookingFee:         req.BookingFee,
		DownPayment:        req.DownPayment,
		TenorMonths:        req.TenorMonths,
		MonthlyInstallment: req.MonthlyInstallment,
		LeasingStatus:      leasingStatus,
		Status:             "draft",
		Notes:              req.Notes,
	}

	if err := tx.Create(&spk).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menerbitkan SPK", "error": err.Error()})
		return
	}

	// Update status unit menjadi booked
	unit.Status = "booked"
	if err := tx.Save(&unit).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengalokasikan unit motor"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal commit database"})
		return
	}

	ctrl.db.Preload("Customer").Preload("Sales").Preload("MotorUnit").First(&spk, spk.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": fmt.Sprintf("Surat Pemesanan Kendaraan %s berhasil diterbitkan", spk.SPKNumber),
		"data":    spk,
	})
}

// GetSPKList mengambil riwayat pemesanan unit kendaraan
func (ctrl *ShowroomController) GetSPKList(c *gin.Context) {
	var list []models.SPKUnit
	query := ctrl.db.Model(&models.SPKUnit{}).
		Preload("Customer").
		Preload("Sales").
		Preload("MotorUnit").
		Order("created_at DESC")

	if search := c.Query("search"); search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("spk_number LIKE ?", searchTerm)
	}

	if err := query.Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil daftar SPK", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// GetSPKReceipt mengembalikan teks konfirmasi pesanan SPK dengan format 48-char (80mm) presisi
func (ctrl *ShowroomController) GetSPKReceipt(c *gin.Context) {
	id := c.Param("id")
	var spk models.SPKUnit

	if err := ctrl.db.Preload("Customer").Preload("Sales").Preload("MotorUnit").First(&spk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "SPK tidak ditemukan"})
		return
	}

	lineLength := 48
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
	sb.WriteString(center("PT NUSANTARA SAKTI GROUP"))
	sb.WriteString(center("DEALER RESMI SEPEDA MOTOR HONDA"))
	sb.WriteString(center("SURAT PEMESANAN KENDARAAN (SPK)"))
	sb.WriteString(doubleSeparator)
	sb.WriteString(row("No. SPK:", spk.SPKNumber))
	sb.WriteString(row("Tanggal:", spk.CreatedAt.Format("02/01/2006 15:04")))
	salesName := "Sales Counter NSS"
	if spk.Sales != nil {
		salesName = spk.Sales.FullName
	}
	sb.WriteString(row("Sales Counter:", salesName))
	sb.WriteString(separator)

	customerName := "-"
	customerPhone := "-"
	if spk.Customer != nil {
		customerName = spk.Customer.Name
		customerPhone = spk.Customer.Phone
	}
	sb.WriteString(row("Nama Pemesan:", customerName))
	sb.WriteString(row("No. HP / WA:", customerPhone))
	sb.WriteString(separator)

	unitName := "-"
	color := "-"
	frameNo := "-"
	engineNo := "-"
	if spk.MotorUnit != nil {
		unitName = spk.MotorUnit.ModelName
		color = spk.MotorUnit.Color
		frameNo = spk.MotorUnit.FrameNumber
		engineNo = spk.MotorUnit.EngineNumber
	}
	sb.WriteString(row("Tipe Motor:", unitName))
	sb.WriteString(row("Warna:", color))
	sb.WriteString(row("No. Rangka:", frameNo))
	sb.WriteString(row("No. Mesin:", engineNo))
	sb.WriteString(separator)

	sb.WriteString(row("Skema Beli:", strings.ToUpper(spk.PurchaseType)))
	if spk.PurchaseType == "credit_nsc" {
		sb.WriteString(row("Lembaga:", "NSC FINANCE"))
		sb.WriteString(row("Uang Muka (DP):", "Rp "+spk.DownPayment.StringFixed(0)))
		sb.WriteString(row("Tenor Cicilan:", fmt.Sprintf("%d Bulan", spk.TenorMonths)))
		sb.WriteString(row("Angsuran/Bln:", "Rp "+spk.MonthlyInstallment.StringFixed(0)))
		sb.WriteString(row("Status Leasing:", strings.ToUpper(spk.LeasingStatus)))
	} else {
		sb.WriteString(row("Tanda Jadi (BF):", "Rp "+spk.BookingFee.StringFixed(0)))
	}
	sb.WriteString(doubleSeparator)
	sb.WriteString(center("* KETENTUAN PENGURUSAN STNK/BPKB AHM *"))
	sb.WriteString(center("1. Estimasi STNK & Plat: 14 hari kerja"))
	sb.WriteString(center("2. Estimasi BPKB: 60-90 hari kerja"))
	sb.WriteString(center("Salam Satu HATI - PT NUSANTARA SAKTI"))
	sb.WriteString(doubleSeparator)
	sb.WriteString("\n")

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"receipt_text": sb.String(),
		"spk":          spk,
	})
}

