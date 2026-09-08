package service

import (
	"doudian/internal/database"
	"doudian/internal/database/model"
	"time"
)

// DashboardStats 仪表盘统计数据
type DashboardStats struct {
	TotalOrders       int64   `json:"total_orders"`
	PendingOrders     int64   `json:"pending_orders"`
	TotalSuppliers    int64   `json:"total_suppliers"`
	TotalProducts     int64   `json:"total_products"`
	TotalShops        int64   `json:"total_shops"`
	TotalPurchaseOrders int64 `json:"total_purchase_orders"`
	TodayOrders       int64   `json:"today_orders"`
	TodayRevenue      float64 `json:"today_revenue"`
}

// GetDashboardStats 获取仪表盘统计数据
func GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats

	// 供应商数
	database.DB.Model(&model.Supplier{}).Count(&stats.TotalSuppliers)

	// 商品数
	database.DB.Model(&model.Product{}).Count(&stats.TotalProducts)

	// 抖店数
	database.DB.Model(&model.Shop{}).Count(&stats.TotalShops)

	// 总订单数
	database.DB.Model(&model.Order{}).Count(&stats.TotalOrders)

	// 待处理订单数
	database.DB.Model(&model.Order{}).Where("status = ?", model.OrderStatusPending).Count(&stats.PendingOrders)

	// 采购单数
	database.DB.Model(&model.PurchaseOrder{}).Count(&stats.TotalPurchaseOrders)

	// 今日订单数
	today := time.Now().Format("2006-01-02")
	database.DB.Model(&model.Order{}).Where("DATE(created_at) = ?", today).Count(&stats.TodayOrders)

	// 今日营收（采购单价 * 数量 的总和）
	type RevenueResult struct {
		Total float64
	}
	var rev RevenueResult
	database.DB.Model(&model.PurchaseOrder{}).
		Select("COALESCE(SUM(purchase_price * quantity), 0) as total").
		Where("DATE(created_at) = ?", today).
		Scan(&rev)
	stats.TodayRevenue = rev.Total

	return &stats, nil
}

// GetRecentOrders 获取最近的订单
func GetRecentOrders(limit int) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}
