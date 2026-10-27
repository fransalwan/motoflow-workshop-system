package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/middleware"
	"motoflow-backend/models"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db: db}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request body", err.Error()))
		return
	}

	var user models.User
	if err := ctrl.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Username atau kata sandi tidak valid", "INVALID_CREDENTIALS"))
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, dto.ErrorResponse("Akun pengguna nonaktif. Hubungi Admin.", "ACCOUNT_INACTIVE"))
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Username atau kata sandi tidak valid", "INVALID_CREDENTIALS"))
		return
	}

	token, err := middleware.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal membuat token autentikasi", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Login berhasil", dto.LoginResponse{
		Token: token,
		User:  user,
	}))
}

func (ctrl *AuthController) Profile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var user models.User
	if err := ctrl.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Pengguna tidak ditemukan", err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("Data profil pengguna", user))
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Data registrasi tidak lengkap", err.Error()))
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
		Phone:    req.Phone,
		IsActive: true,
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Gagal mengenkripsi kata sandi", err.Error()))
		return
	}

	if err := ctrl.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Username atau Email sudah terdaftar", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse("Pengguna berhasil didaftarkan", user))
}

