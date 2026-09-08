package controller

import (
	"doudian/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats 获取仪表盘统计数据
// GET /api/dashboard/stats
func GetDashboardStats(c *gin.Context) {
	stats, err := service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取统计数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    stats,
	})
}

// GetRecentOrders 获取最近的待处理订单
// GET /api/dashboard/recent-orders
func GetRecentOrders(c *gin.Context) {
	orders, err := service.GetRecentOrders(5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取最近订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    orders,
	})
}
