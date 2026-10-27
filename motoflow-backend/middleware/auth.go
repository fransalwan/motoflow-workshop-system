package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"motoflow-backend/config"
	"motoflow-backend/dto"
	"motoflow-backend/models"
)

type JWTClaims struct {
	UserID   uint        `json:"user_id"`
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "motoflow-ahass",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(config.AppConfig.JWTSecret)
	return token.SignedString(secret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Authorization header is required", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid authorization header format", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		secret := []byte(config.AppConfig.JWTSecret)

		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Invalid or expired token", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse("Failed to parse token claims", "UNAUTHORIZED"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireRole(allowedRoles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, dto.ErrorResponse("Access denied: role not found", "FORBIDDEN"))
			c.Abort()
			return
		}

		userRole := roleVal.(models.Role)
		// Admin always has full access
		if userRole == models.RoleAdmin {
			c.Next()
			return
		}

		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, dto.ErrorResponse("Insufficient privileges for this action", "FORBIDDEN"))
		c.Abort()
	}
}

