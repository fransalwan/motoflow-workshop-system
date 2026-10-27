package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"motoflow-backend/dto"
	"motoflow-backend/models"
	"motoflow-backend/services"
)

type DashboardController struct {
	db *gorm.DB
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{db: db}
}

type TopPartStat struct {
	PartNumber   string          `json:"part_number"`
	PartName     string          `json:"part_name"`
	TotalSoldQty int             `json:"total_sold_qty"`
	TotalRevenue decimal.Decimal `json:"total_revenue"`
}

type MechanicStat struct {
	MechanicID    uint            `json:"mechanic_id"`
	Code          string          `json:"code"`
	Name          string          `json:"name"`
	PitNumber     int             `json:"pit_number"`
	CompletedJobs int64           `json:"completed_jobs"`
	LaborRevenue  decimal.Decimal `json:"labor_revenue"`
}

type DailyRevenueStat struct {
	Date         string          `json:"date"`
	TotalRevenue decimal.Decimal `json:"total_revenue"`
	BikeCount    int64           `json:"bike_count"`
}

type DashboardSummary struct {
	TodayRevenue      decimal.Decimal    `json:"today_revenue"`
	TodayBikesServiced int64              `json:"today_bikes_serviced"`
	MonthRevenue      decimal.Decimal    `json:"month_revenue"`
	MonthBikesServiced int64              `json:"month_bikes_serviced"`
	ActiveJobsCount   int64              `json:"active_jobs_count"`
	LowStockCount     int64              `json:"low_stock_count"`
	StatusCounts      map[string]int64   `json:"status_counts"`
	TopSellingParts   []TopPartStat      `json:"top_selling_parts"`
	MechanicRankings  []MechanicStat     `json:"mechanic_rankings"`
	DailyTrend        []DailyRevenueStat `json:"daily_trend"`
}

func (ctrl *DashboardController) GetSummary(c *gin.Context) {
	cacheKey := "dashboard:summary"
	if cached, found := services.GetCache().Get(cacheKey); found {
		if summary, ok := cached.(DashboardSummary); ok {
			c.JSON(http.StatusOK, dto.SuccessResponse("Ringkasan dashboard eksekutif (cached)", summary))
			return
		}
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 1. Today's Revenue and Bike Count
	var todayRevenue decimal.Decimal
	var todayBikes int64

	ctrl.db.Model(&models.WorkOrder{}).
		Where("payment_status = ? AND paid_at >= ?", models.PaymentStatusPaid, startOfDay).
		Select("COALESCE(SUM(grand_total), 0)").
		Scan(&todayRevenue)

	ctrl.db.Model(&models.WorkOrder{}).
		Where("created_at >= ?", startOfDay).
		Count(&todayBikes)

	// 2. Month-to-Date Revenue and Bike Count
	var monthRevenue decimal.Decimal
	var monthBikes int64

	ctrl.db.Model(&models.WorkOrder{}).
		Where("payment_status = ? AND paid_at >= ?", models.PaymentStatusPaid, startOfMonth).
		Select("COALESCE(SUM(grand_total), 0)").
		Scan(&monthRevenue)

	ctrl.db.Model(&models.WorkOrder{}).
		Where("created_at >= ?", startOfMonth).
		Count(&monthBikes)

	// 3. Active Jobs in AHASS Pits
	var activeJobs int64
	ctrl.db.Model(&models.WorkOrder{}).
		Where("status IN ?", []string{string(models.WOStatusInProgress), string(models.WOStatusDraft)}).
		Count(&activeJobs)

	// 4. Low Stock Alerts (Honda Genuine Parts)
	var lowStockCount int64
	ctrl.db.Model(&models.Sparepart{}).
		Where("type = ? AND stock <= min_stock", models.ItemTypePart).
		Count(&lowStockCount)

	// 5. Status Breakdown
	type StatusCountRow struct {
		Status string
		Count  int64
	}
	var statusRows []StatusCountRow
	ctrl.db.Model(&models.WorkOrder{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusRows)

	statusMap := map[string]int64{
		"draft":       0,
		"in_progress": 0,
		"completed":   0,
		"paid":        0,
		"cancelled":   0,
	}
	for _, row := range statusRows {
		statusMap[row.Status] = row.Count
	}

	// 6. Top 5 Honda Genuine Parts Terlaris
	var topParts []TopPartStat
	ctrl.db.Model(&models.WorkOrderItem{}).
		Select("part_number_snapshot as part_number, item_name_snapshot as part_name, SUM(quantity) as total_sold_qty, SUM(subtotal) as total_revenue").
		Where("item_type = ?", models.ItemTypePart).
		Group("part_number_snapshot, item_name_snapshot").
		Order("total_revenue desc").
		Limit(5).
		Scan(&topParts)

	// 7. Mechanic Productivity Rankings
	var mechanics []models.Mechanic
	ctrl.db.Find(&mechanics)

	var mechanicRankings []MechanicStat
	for _, m := range mechanics {
		var jobCount int64
		var laborRev decimal.Decimal

		ctrl.db.Model(&models.WorkOrder{}).
			Where("mechanic_id = ? AND payment_status = ?", m.ID, models.PaymentStatusPaid).
			Count(&jobCount)

		ctrl.db.Model(&models.WorkOrder{}).
			Where("mechanic_id = ? AND payment_status = ?", m.ID, models.PaymentStatusPaid).
			Select("COALESCE(SUM(services_subtotal), 0)").
			Scan(&laborRev)

		mechanicRankings = append(mechanicRankings, MechanicStat{
			MechanicID:    m.ID,
			Code:          m.Code,
			Name:          m.Name,
			PitNumber:     m.PitNumber,
			CompletedJobs: jobCount,
			LaborRevenue:  laborRev,
		})
	}

	// 8. 7-Day Trend
	var dailyTrend []DailyRevenueStat
	for i := 6; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		dayEnd := dayStart.Add(24 * time.Hour)

		var dayRev decimal.Decimal
		var dayCount int64

		ctrl.db.Model(&models.WorkOrder{}).
			Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", models.PaymentStatusPaid, dayStart, dayEnd).
			Select("COALESCE(SUM(grand_total), 0)").
			Scan(&dayRev)

		ctrl.db.Model(&models.WorkOrder{}).
			Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
			Count(&dayCount)

		dailyTrend = append(dailyTrend, DailyRevenueStat{
			Date:         day.Format("02 Jan"),
			TotalRevenue: dayRev,
			BikeCount:    dayCount,
		})
	}

	summary := DashboardSummary{
		TodayRevenue:       todayRevenue,
		TodayBikesServiced: todayBikes,
		MonthRevenue:       monthRevenue,
		MonthBikesServiced: monthBikes,
		ActiveJobsCount:    activeJobs,
		LowStockCount:      lowStockCount,
		StatusCounts:       statusMap,
		TopSellingParts:    topParts,
		MechanicRankings:   mechanicRankings,
		DailyTrend:         dailyTrend,
	}

	services.GetCache().Set(cacheKey, summary, 1*time.Minute)

	c.JSON(http.StatusOK, dto.SuccessResponse("Ringkasan dashboard eksekutif", summary))
}

