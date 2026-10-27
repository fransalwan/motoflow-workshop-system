package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
)

type SparepartController struct {
	db *gorm.DB
}

func NewSparepartController(db *gorm.DB) *SparepartController {
	return &SparepartController{db: db}
}

func (ctrl *SparepartController) List(c *gin.Context) {
	search := c.Query("search")
	categoryID := c.Query("category_id")
	itemType := c.Query("type")
	lowStockOnly := c.Query("low_stock") == "true"

	query := ctrl.db.Model(&models.Sparepart{}).Preload("Category")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("part_number LIKE ? OR name LIKE ?", searchTerm, searchTerm)
	}

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if itemType != "" {
		query = query.Where("type = ?", itemType)
	}

	if lowStockOnly {
		// Only parts can have physical low stock alerts (not services)
		query = query.Where("type = ? AND stock <= min_stock", models.ItemTypePart)
	}

	var items []models.Sparepart
	if err := query.Order("id desc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data suku cadang", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Daftar suku cadang HGP & Jasa Servis", items))
}

func (ctrl *SparepartController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var item models.Sparepart

	if err := ctrl.db.Preload("Category").First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Suku cadang tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Detail suku cadang", item))
}

func (ctrl *SparepartController) Create(c *gin.Context) {
	var req dto.CreateSparepartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data suku cadang tidak valid", err.Error()))
		return
	}

	item := models.Sparepart{
		PartNumber:   req.PartNumber,
		Name:         req.Name,
		CategoryID:   req.CategoryID,
		Type:         req.Type,
		CostPrice:    req.CostPrice,
		SellingPrice: req.SellingPrice,
		Stock:        req.Stock,
		MinStock:     req.MinStock,
		Unit:         req.Unit,
		Location:     req.Location,
		IsActive:     true,
	}

	if err := ctrl.db.Create(&item).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Nomor part sudah terdaftar atau kategori tidak valid", err.Error()))
		return
	}

	ctrl.db.Preload("Category").First(&item, item.ID)
	c.JSON(http.StatusCreated, dto.SuccessResponse("Suku cadang berhasil ditambahkan", item))
}

func (ctrl *SparepartController) Update(c *gin.Context) {
	id := c.Param("id")
	var item models.Sparepart

	if err := ctrl.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Suku cadang tidak ditemukan", err.Error()))
		return
	}

	var req dto.UpdateSparepartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data update tidak valid", err.Error()))
		return
	}

	if req.PartNumber != "" {
		item.PartNumber = req.PartNumber
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.CategoryID != 0 {
		item.CategoryID = req.CategoryID
	}
	if req.Type != "" {
		item.Type = req.Type
	}
	if !req.CostPrice.IsZero() {
		item.CostPrice = req.CostPrice
	}
	if !req.SellingPrice.IsZero() {
		item.SellingPrice = req.SellingPrice
	}
	item.Stock = req.Stock
	if req.MinStock != 0 {
		item.MinStock = req.MinStock
	}
	if req.Unit != "" {
		item.Unit = req.Unit
	}
	if req.Location != "" {
		item.Location = req.Location
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := ctrl.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memperbarui suku cadang", err.Error()))
		return
	}

	ctrl.db.Preload("Category").First(&item, item.ID)
	c.JSON(http.StatusOK, dto.SuccessResponse("Suku cadang berhasil diperbarui", item))
}

func (ctrl *SparepartController) UpdateStock(c *gin.Context) {
	id := c.Param("id")
	var item models.Sparepart

	if err := ctrl.db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Suku cadang tidak ditemukan", err.Error()))
		return
	}

	var req dto.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data kuantitas tidak valid", err.Error()))
		return
	}

	newStock := item.Stock + req.Quantity
	if newStock < 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Stok tidak boleh kurang dari nol", "INSUFFICIENT_STOCK"))
		return
	}

	item.Stock = newStock
	if err := ctrl.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menyesuaikan stok", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Penyesuaian stok berhasil", item))
}

func (ctrl *SparepartController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.db.Delete(&models.Sparepart{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghapus suku cadang", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Suku cadang berhasil dihapus", nil))
}

func (ctrl *SparepartController) ListCategories(c *gin.Context) {
	var categories []models.Category
	if err := ctrl.db.Order("name asc").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data kategori", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Daftar kategori suku cadang & servis", categories))
}

