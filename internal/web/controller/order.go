package controller

import (
	"doudian/internal/database/model"
	"doudian/internal/service"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrders 获取订单列表（分页+状态筛选+抖店筛选+搜索）
// GET /api/orders
func ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	shopID, _ := strconv.ParseUint(c.Query("shop_id"), 10, 32)
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := service.ListOrders(page, pageSize, status, uint(shopID), keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取订单列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    result,
	})
}

// GetOrder 获取订单详情（含采购单）
// GET /api/orders/:id
func GetOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	order, purchaseOrders, err := service.GetOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data": gin.H{
			"order":           order,
			"purchase_orders": purchaseOrders,
		},
	})
}

// CreateOrder 创建订单
// POST /api/orders
func CreateOrder(c *gin.Context) {
	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    order,
	})
}

// DeleteOrder 删除订单
// DELETE /api/orders/:id
func DeleteOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	if err := service.DeleteOrder(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除订单失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}

// UpdateOrder 更新订单
// PUT /api/orders/:id
func UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	var order model.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	if err := service.UpdateOrder(uint(id), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新订单失败: " + err.Error()})
		return
	}

	updatedOrder, _, _ := service.GetOrder(uint(id))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功", "data": updatedOrder})
}

// ProcessAllOrders 一键匹配所有待处理订单
// POST /api/orders/match
func ProcessAllOrders(c *gin.Context) {
	successCount, failCount, err := service.ProcessAllPendingOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "处理失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "处理完成",
		"data": gin.H{"success_count": successCount, "fail_count": failCount},
	})
}

// ImportOrders CSV导入订单
// POST /api/orders/import
func ImportOrders(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传CSV文件",
		})
		return
	}
	defer file.Close()

	count, err := service.ImportOrders(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "导入失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "导入成功",
		"data": gin.H{
			"count": count,
		},
	})
}

// ExportOrders CSV导出订单
// GET /api/orders/export
func ExportOrders(c *gin.Context) {
	orders, err := service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=orders.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.OrderCSVHeader())

	for i := range orders {
		writer.Write(service.OrderToCSVRow(&orders[i]))
	}
}

// OrderTemplate 下载CSV模板
// GET /api/orders/template
func OrderTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=orders_template.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.OrderCSVHeader())
	writer.Write([]string{"DD20260901001", "SKU-001", "示例商品", "1", "张三", "13800138000", "北京市朝阳区", "1", "pending"})
}
