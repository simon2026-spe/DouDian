package controller

import (
	"doudian/internal/database/model"
	"doudian/internal/service"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListPurchaseOrders 获取采购单列表（分页）
// GET /api/purchase-orders
func ListPurchaseOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := service.ListPurchaseOrders(page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取采购单列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    result,
	})
}

// GetPurchaseOrder 获取采购单详情
// GET /api/purchase-orders/:id
func GetPurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	purchaseOrder, err := service.GetPurchaseOrder(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "采购单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    purchaseOrder,
	})
}

// GeneratePurchaseOrder 为订单生成采购单
// POST /api/purchase-orders/generate/:order_id
func GeneratePurchaseOrder(c *gin.Context) {
	orderIDStr := c.Param("order_id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的订单ID",
		})
		return
	}

	purchaseOrder, err := service.GeneratePurchaseOrder(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成采购单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "生成成功",
		"data":    purchaseOrder,
	})
}

// CreatePurchaseOrder 手动新建采购单
// POST /api/purchase-orders
func CreatePurchaseOrder(c *gin.Context) {
	var purchaseOrder model.PurchaseOrder
	if err := c.ShouldBindJSON(&purchaseOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.CreatePurchaseOrder(&purchaseOrder); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建采购单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    purchaseOrder,
	})
}

// UpdatePurchaseOrder 更新采购单状态+物流单号
// PUT /api/purchase-orders/:id
func UpdatePurchaseOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	var req struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.UpdatePurchaseOrder(uint(id), req.Status, req.TrackingNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新采购单失败: " + err.Error(),
		})
		return
	}

	updatedPO, _ := service.GetPurchaseOrder(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    updatedPO,
	})
}

// UpdatePurchaseOrderStatus 更新采购单状态
// PUT /api/purchase-orders/:id/status
func UpdatePurchaseOrderStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	var req struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	if err := service.UpdatePurchaseOrder(uint(id), req.Status, req.TrackingNumber); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新采购单失败: " + err.Error()})
		return
	}

	updatedPO, _ := service.GetPurchaseOrder(uint(id))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功", "data": updatedPO})
}

// ExportPurchaseOrders CSV导出采购单
// GET /api/purchase-orders/export
func ExportPurchaseOrders(c *gin.Context) {
	purchaseOrders, err := service.GetAllPurchaseOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=purchase_orders.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.PurchaseOrderCSVHeader())

	for i := range purchaseOrders {
		writer.Write(service.PurchaseOrderToCSVRow(&purchaseOrders[i]))
	}
}
