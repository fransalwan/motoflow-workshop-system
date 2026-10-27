package dto

import "motoflow-backend/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type RegisterRequest struct {
	Username string      `json:"username" binding:"required,min=3,max=50"`
	Password string      `json:"password" binding:"required,min=6"`
	FullName string      `json:"full_name" binding:"required"`
	Email    string      `json:"email" binding:"required,email"`
	Role     models.Role `json:"role" binding:"required"`
	Phone    string      `json:"phone"`
}

