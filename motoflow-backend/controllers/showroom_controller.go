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

	userID, exists := c.Get("user_id")
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

// GetSPKReceipt mengembalikan teks konfirmasi pesanan SPK
func (ctrl *ShowroomController) GetSPKReceipt(c *gin.Context) {
	id := c.Param("id")
	var spk models.SPKUnit

	if err := ctrl.db.Preload("Customer").Preload("Sales").Preload("MotorUnit").First(&spk, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "SPK tidak ditemukan"})
		return
	}

	var sb strings.Builder
	sb.WriteString("====================================================\n")
	sb.WriteString("             PT NUSANTARA SAKTI GROUP\n")
	sb.WriteString("        DEALER RESMI SEPEDA MOTOR HONDA\n")
	sb.WriteString("        SURAT PEMESANAN KENDARAAN (SPK)\n")
	sb.WriteString("====================================================\n")
	sb.WriteString(fmt.Sprintf("No. SPK        : %s\n", spk.SPKNumber))
	sb.WriteString(fmt.Sprintf("Tanggal        : %s\n", spk.CreatedAt.Format("02/01/2006 15:04")))
	salesName := "Sales Counter NSS"
	if spk.Sales != nil {
		salesName = spk.Sales.FullName
	}
	sb.WriteString(fmt.Sprintf("Sales Counter  : %s\n", salesName))
	sb.WriteString("----------------------------------------------------\n")
	customerName := "-"
	customerPhone := "-"
	if spk.Customer != nil {
		customerName = spk.Customer.Name
		customerPhone = spk.Customer.Phone
	}
	sb.WriteString(fmt.Sprintf("Nama Pemesan   : %s\n", customerName))
	sb.WriteString(fmt.Sprintf("No. HP / WA    : %s\n", customerPhone))
	sb.WriteString("----------------------------------------------------\n")
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
	sb.WriteString(fmt.Sprintf("Tipe Motor     : %s\n", unitName))
	sb.WriteString(fmt.Sprintf("Warna          : %s\n", color))
	sb.WriteString(fmt.Sprintf("Nomor Rangka   : %s\n", frameNo))
	sb.WriteString(fmt.Sprintf("Nomor Mesin    : %s\n", engineNo))
	sb.WriteString("----------------------------------------------------\n")
	sb.WriteString(fmt.Sprintf("Skema Pembelian: %s\n", strings.ToUpper(spk.PurchaseType)))
	if spk.PurchaseType == "credit_nsc" {
		sb.WriteString(fmt.Sprintf("Lembaga Leasing: PT NUSA SURYA CIPTADANA (NSC FINANCE)\n"))
		sb.WriteString(fmt.Sprintf("Uang Muka (DP) : %s\n", formatRupiah(spk.DownPayment)))
		sb.WriteString(fmt.Sprintf("Tenor Cicilan  : %d Bulan\n", spk.TenorMonths))
		sb.WriteString(fmt.Sprintf("Angsuran/Bulan : %s / bulan\n", formatRupiah(spk.MonthlyInstallment)))
		sb.WriteString(fmt.Sprintf("Status Leasing : %s\n", strings.ToUpper(spk.LeasingStatus)))
	} else {
		sb.WriteString(fmt.Sprintf("Tanda Jadi (BF): %s\n", formatRupiah(spk.BookingFee)))
	}
	sb.WriteString("====================================================\n")
	sb.WriteString("  * SYARAT & KETENTUAN PENGURUSAN STNK/BPKB AHM *\n")
	sb.WriteString(" 1. Estimasi STNK & Plat Nomor: 14 hari kerja\n")
	sb.WriteString(" 2. Estimasi BPKB: 60-90 hari kerja via Samsat\n")
	sb.WriteString("        Salam Satu HATI - Honda Juara!\n")
	sb.WriteString("====================================================\n")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"receipt_text": sb.String(),
		"spk": spk,
	})
}
