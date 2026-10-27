package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
	"motoflow-backend/services"
)

type WorkOrderController struct {
	db               *gorm.DB
	financialService *services.FinancialService
}

func NewWorkOrderController(db *gorm.DB) *WorkOrderController {
	return &WorkOrderController{
		db:               db,
		financialService: services.NewFinancialService(),
	}
}

func (ctrl *WorkOrderController) List(c *gin.Context) {
	status := c.Query("status")
	search := c.Query("search")
	paymentStatus := c.Query("payment_status")

	query := ctrl.db.Model(&models.WorkOrder{}).
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceAdvisor").
		Preload("Items")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("LEFT JOIN customers ON customers.id = work_orders.customer_id").
			Where("work_orders.wo_number LIKE ? OR customers.name LIKE ? OR customers.plate_number LIKE ?",
				searchTerm, searchTerm, searchTerm)
	}

	var workOrders []models.WorkOrder
	if err := query.Order("work_orders.id desc").Find(&workOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data perintah kerja", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Daftar Work Order / Faktur Servis", workOrders))
}

func (ctrl *WorkOrderController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var wo models.WorkOrder

	if err := ctrl.db.
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceAdvisor").
		Preload("Items.Sparepart").
		First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Perintah kerja tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Detail Work Order", wo))
}

func (ctrl *WorkOrderController) Create(c *gin.Context) {
	var req dto.CreateWorkOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Payload Work Order tidak valid", err.Error()))
		return
	}

	var customer models.Customer
	if err := ctrl.db.First(&customer, req.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Pelanggan tidak ditemukan", err.Error()))
		return
	}

	// Update current km if provided
	if req.EntryKM > 0 {
		customer.CurrentKM = req.EntryKM
		ctrl.db.Save(&customer)
	}

	// Retrieve Service Advisor from context
	var saID *uint
	if val, exists := c.Get("userID"); exists {
		id := val.(uint)
		saID = &id
	}

	taxRate := decimal.NewFromInt(11) // Default PPN 11%
	if req.TaxRate != nil {
		taxRate = *req.TaxRate
	}

	orderDiscount := decimal.Zero
	if req.TotalDiscount != nil {
		orderDiscount = *req.TotalDiscount
	}

	// Prepare WorkOrder entity & transaction
	var createdWO models.WorkOrder

	err := ctrl.db.Transaction(func(tx *gorm.DB) error {
		// Collision-free sequential WO number
		woNumber, err := services.GenerateWorkOrderNumber(tx)
		if err != nil {
			return fmt.Errorf("gagal men-generate nomor WO: %w", err)
		}

		// Convert input items to model items
		var orderItems []models.WorkOrderItem
		for _, itm := range req.Items {
			subtotal := ctrl.financialService.CalculateItemSubtotal(itm.Price, itm.Quantity, itm.DiscountPercent)

			orderItem := models.WorkOrderItem{
				SparepartID:        itm.SparepartID,
				ItemType:           itm.ItemType,
				PartNumberSnapshot: itm.PartNumber,
				ItemNameSnapshot:   itm.ItemName,
				Quantity:           itm.Quantity,
				PriceSnapshot:      itm.Price,
				DiscountPercent:    itm.DiscountPercent,
				Subtotal:           subtotal,
				MechanicNotes:      itm.MechanicNotes,
			}

			// If it's a physical part, deduct stock and verify
			if itm.SparepartID != nil && itm.ItemType == models.ItemTypePart {
				var part models.Sparepart
				if err := tx.First(&part, *itm.SparepartID).Error; err != nil {
					return fmt.Errorf("suku cadang ID %d tidak ditemukan", *itm.SparepartID)
				}

				if part.Stock < itm.Quantity {
					return fmt.Errorf("stok suku cadang '%s' tidak mencukupi (sisa: %d, diminta: %d)", part.Name, part.Stock, itm.Quantity)
				}

				part.Stock -= itm.Quantity
				if err := tx.Save(&part).Error; err != nil {
					return fmt.Errorf("gagal memotong stok suku cadang '%s': %w", part.Name, err)
				}
			}

			orderItems = append(orderItems, orderItem)
		}

		// Calculate financial totals
		totals := ctrl.financialService.CalculateWorkOrderTotals(orderItems, orderDiscount, taxRate, decimal.Zero)

		wo := models.WorkOrder{
			WONumber:         woNumber,
			CustomerID:       req.CustomerID,
			MechanicID:       req.MechanicID,
			ServiceAdvisorID: saID,
			EntryKM:          req.EntryKM,
			Complaints:       req.Complaints,
			ServiceDiagnosis: req.ServiceDiagnosis,
			Status:           models.WOStatusDraft,
			PaymentStatus:    models.PaymentStatusUnpaid,
			PaymentMethod:    models.PaymentCash,
			PartsSubtotal:    totals.PartsSubtotal,
			ServicesSubtotal: totals.ServicesSubtotal,
			TotalDiscount:    totals.TotalDiscount,
			TaxRate:          totals.TaxRate,
			TaxAmount:        totals.TaxAmount,
			GrandTotal:       totals.GrandTotal,
			PaidAmount:       decimal.Zero,
			ChangeAmount:     decimal.Zero,
			Notes:            req.Notes,
			Items:            orderItems,
		}

		if err := tx.Create(&wo).Error; err != nil {
			return err
		}

		createdWO = wo
		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Gagal membuat Work Order", err.Error()))
		return
	}

	// Preload full details
	ctrl.db.
		Preload("Customer").
		Preload("Mechanic").
		Preload("ServiceAdvisor").
		Preload("Items.Sparepart").
		First(&createdWO, createdWO.ID)

	c.JSON(http.StatusCreated, dto.SuccessResponse("Work Order berhasil diterbitkan", createdWO))
}

func (ctrl *WorkOrderController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var wo models.WorkOrder

	if err := ctrl.db.First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Work Order tidak ditemukan", err.Error()))
		return
	}

	var req dto.UpdateWorkOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Status tidak valid", err.Error()))
		return
	}

	wo.Status = req.Status
	if req.Notes != "" {
		wo.Notes = req.Notes
	}

	now := time.Now()
	if req.Status == models.WOStatusCompleted && wo.CompletedAt == nil {
		wo.CompletedAt = &now
	}

	if err := ctrl.db.Save(&wo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengubah status Work Order", err.Error()))
		return
	}

	ctrl.db.Preload("Customer").Preload("Mechanic").Preload("Items").First(&wo, wo.ID)
	c.JSON(http.StatusOK, dto.SuccessResponse("Status Work Order berhasil diperbarui", wo))
}

func (ctrl *WorkOrderController) ProcessPayment(c *gin.Context) {
	id := c.Param("id")
	var wo models.WorkOrder

	if err := ctrl.db.Preload("Items").First(&wo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Work Order tidak ditemukan", err.Error()))
		return
	}

	var req dto.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data pembayaran tidak valid", err.Error()))
		return
	}

	if req.PaidAmount.LessThan(wo.GrandTotal) {
		shortage := wo.GrandTotal.Sub(req.PaidAmount)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse(
			fmt.Sprintf("Jumlah bayar kurang sebesar Rp %s", shortage.StringFixed(2)),
			"INSUFFICIENT_PAYMENT",
		))
		return
	}

	change := req.PaidAmount.Sub(wo.GrandTotal)
	now := time.Now()

	wo.PaymentMethod = req.PaymentMethod
	wo.PaymentStatus = models.PaymentStatusPaid
	wo.Status = models.WOStatusPaid
	wo.PaidAmount = req.PaidAmount
	wo.ChangeAmount = change
	wo.PaidAt = &now

	if req.Notes != "" {
		wo.Notes = req.Notes
	}

	if err := ctrl.db.Save(&wo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memproses pembayaran kasir", err.Error()))
		return
	}

	ctrl.db.Preload("Customer").Preload("Mechanic").Preload("ServiceAdvisor").Preload("Items").First(&wo, wo.ID)
	c.JSON(http.StatusOK, dto.SuccessResponse("Pembayaran kasir berhasil diselesaikan", wo))
}

func (ctrl *WorkOrderController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.db.Delete(&models.WorkOrder{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal membatalkan Work Order", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Work Order berhasil dibatalkan", nil))
}

