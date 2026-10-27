package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
)

type MechanicController struct {
	db *gorm.DB
}

func NewMechanicController(db *gorm.DB) *MechanicController {
	return &MechanicController{db: db}
}

func (ctrl *MechanicController) List(c *gin.Context) {
	availableOnly := c.Query("available") == "true"
	var mechanics []models.Mechanic

	query := ctrl.db.Model(&models.Mechanic{})
	if availableOnly {
		query = query.Where("is_available = ?", true)
	}

	if err := query.Order("pit_number asc").Find(&mechanics).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data mekanik", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Daftar mekanik AHASS", mechanics))
}

func (ctrl *MechanicController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var mechanic models.Mechanic

	if err := ctrl.db.First(&mechanic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Mekanik tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Detail mekanik", mechanic))
}

func (ctrl *MechanicController) Create(c *gin.Context) {
	var req dto.CreateMechanicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data mekanik tidak valid", err.Error()))
		return
	}

	mechanic := models.Mechanic{
		Code:           req.Code,
		Name:           req.Name,
		Phone:          req.Phone,
		Specialization: req.Specialization,
		PitNumber:      req.PitNumber,
		IsAvailable:    req.IsAvailable,
	}

	if err := ctrl.db.Create(&mechanic).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Kode mekanik sudah terdaftar", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Mekanik berhasil didaftarkan", mechanic))
}

func (ctrl *MechanicController) Update(c *gin.Context) {
	id := c.Param("id")
	var mechanic models.Mechanic

	if err := ctrl.db.First(&mechanic, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Mekanik tidak ditemukan", err.Error()))
		return
	}

	var req dto.UpdateMechanicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data update tidak valid", err.Error()))
		return
	}

	if req.Code != "" {
		mechanic.Code = req.Code
	}
	if req.Name != "" {
		mechanic.Name = req.Name
	}
	if req.Phone != "" {
		mechanic.Phone = req.Phone
	}
	if req.Specialization != "" {
		mechanic.Specialization = req.Specialization
	}
	if req.PitNumber != 0 {
		mechanic.PitNumber = req.PitNumber
	}
	if req.IsAvailable != nil {
		mechanic.IsAvailable = *req.IsAvailable
	}

	if err := ctrl.db.Save(&mechanic).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memperbarui data mekanik", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Data mekanik berhasil diperbarui", mechanic))
}

func (ctrl *MechanicController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.db.Delete(&models.Mechanic{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghapus data mekanik", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Mekanik berhasil dihapus", nil))
}

