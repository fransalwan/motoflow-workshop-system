package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"motoflow-backend/controllers"
	"motoflow-backend/dto"
	"motoflow-backend/middleware"
	"motoflow-backend/models"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Global CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"service":   "motoflow-backend",
			"ahass":     "PT Nusantara Sakti Group",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Controllers
	authCtrl := controllers.NewAuthController(db)
	customerCtrl := controllers.NewCustomerController(db)
	sparepartCtrl := controllers.NewSparepartController(db)
	mechanicCtrl := controllers.NewMechanicController(db)
	workOrderCtrl := controllers.NewWorkOrderController(db)
	dashboardCtrl := controllers.NewDashboardController(db)
	printCtrl := controllers.NewPrintController(db)

	api := r.Group("/api/v1")
	{
		// Public Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", authCtrl.Login)
		}

		// Authenticated Routes
		authRequired := api.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			// Auth Profile & Registration
			authRequired.GET("/auth/profile", authCtrl.Profile)
			authRequired.POST("/auth/register", middleware.RequireRole(models.RoleAdmin), authCtrl.Register)

			// Executive Dashboard
			authRequired.GET("/dashboard/summary", dashboardCtrl.GetSummary)

			// Customers & Honda Motorcycles
			customers := authRequired.Group("/customers")
			{
				customers.GET("", customerCtrl.List)
				customers.GET("/:id", customerCtrl.GetByID)
				customers.POST("", customerCtrl.Create)
				customers.PUT("/:id", customerCtrl.Update)
				customers.DELETE("/:id", customerCtrl.Delete)
			}

			// Spareparts & Services (Honda Genuine Parts)
			spareparts := authRequired.Group("/spareparts")
			{
				spareparts.GET("", sparepartCtrl.List)
				spareparts.GET("/categories", sparepartCtrl.ListCategories)
				spareparts.GET("/:id", sparepartCtrl.GetByID)
				spareparts.POST("", sparepartCtrl.Create)
				spareparts.PUT("/:id", sparepartCtrl.Update)
				spareparts.PATCH("/:id/stock", sparepartCtrl.UpdateStock)
				spareparts.DELETE("/:id", sparepartCtrl.Delete)
			}

			// Mechanics
			mechanics := authRequired.Group("/mechanics")
			{
				mechanics.GET("", mechanicCtrl.List)
				mechanics.GET("/:id", mechanicCtrl.GetByID)
				mechanics.POST("", mechanicCtrl.Create)
				mechanics.PUT("/:id", mechanicCtrl.Update)
				mechanics.DELETE("/:id", mechanicCtrl.Delete)
			}

			// Work Orders (POS Billing & Service Orders)
			workOrders := authRequired.Group("/work-orders")
			{
				workOrders.GET("", workOrderCtrl.List)
				workOrders.GET("/:id", workOrderCtrl.GetByID)
				workOrders.POST("", workOrderCtrl.Create)
				workOrders.PATCH("/:id/status", workOrderCtrl.UpdateStatus)
				workOrders.POST("/:id/pay", workOrderCtrl.ProcessPayment)
				workOrders.DELETE("/:id", workOrderCtrl.Delete)

				// Print Engine (Thermal Receipt & A4 PDF)
				workOrders.GET("/:id/receipt", printCtrl.ThermalReceipt)
				workOrders.GET("/:id/invoice-pdf", printCtrl.ServiceInvoicePDF)
			}
		}
	}

	// 404 fallback
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, dto.ErrorResponse("Endpoint tidak ditemukan", "NOT_FOUND"))
	})

	return r
}

