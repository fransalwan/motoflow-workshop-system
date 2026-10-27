package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
)

type CustomerController struct {
	db *gorm.DB
}

func NewCustomerController(db *gorm.DB) *CustomerController {
	return &CustomerController{db: db}
}

func (ctrl *CustomerController) List(c *gin.Context) {
	search := c.Query("search")
	var customers []models.Customer

	query := ctrl.db.Model(&models.Customer{})
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR plate_number LIKE ? OR motor_model LIKE ?",
			searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if err := query.Order("id desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengambil data pelanggan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Daftar pelanggan & motor", customers))
}

func (ctrl *CustomerController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := ctrl.db.Preload("WorkOrders", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Pelanggan tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Detail pelanggan", customer))
}

func (ctrl *CustomerController) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data pelanggan tidak valid", err.Error()))
		return
	}

	customer := models.Customer{
		Name:         req.Name,
		Phone:        req.Phone,
		Email:        req.Email,
		Address:      req.Address,
		PlateNumber:  req.PlateNumber,
		MotorModel:   req.MotorModel,
		Year:         req.Year,
		FrameNumber:  req.FrameNumber,
		EngineNumber: req.EngineNumber,
		CurrentKM:    req.CurrentKM,
	}

	if err := ctrl.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menambahkan pelanggan", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Pelanggan & unit motor berhasil didaftarkan", customer))
}

func (ctrl *CustomerController) Update(c *gin.Context) {
	id := c.Param("id")
	var customer models.Customer

	if err := ctrl.db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Pelanggan tidak ditemukan", err.Error()))
		return
	}

	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data update tidak valid", err.Error()))
		return
	}

	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Address != "" {
		customer.Address = req.Address
	}
	if req.PlateNumber != "" {
		customer.PlateNumber = req.PlateNumber
	}
	if req.MotorModel != "" {
		customer.MotorModel = req.MotorModel
	}
	if req.Year != 0 {
		customer.Year = req.Year
	}
	if req.FrameNumber != "" {
		customer.FrameNumber = req.FrameNumber
	}
	if req.EngineNumber != "" {
		customer.EngineNumber = req.EngineNumber
	}
	if req.CurrentKM != 0 {
		customer.CurrentKM = req.CurrentKM
	}

	if err := ctrl.db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal memperbarui data pelanggan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Data pelanggan berhasil diperbarui", customer))
}

func (ctrl *CustomerController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.db.Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal menghapus data pelanggan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Pelanggan berhasil dihapus", nil))
}

