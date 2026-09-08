package controller

import (
	"doudian/internal/database/model"
	"doudian/internal/service"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListSuppliers 获取供应商列表（分页+搜索）
// GET /api/suppliers
func ListSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := service.ListSuppliers(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取供应商列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    result,
	})
}

// GetAllSuppliers 获取所有供应商（不分页，用于下拉选择）
// GET /api/suppliers/all
func GetAllSuppliers(c *gin.Context) {
	suppliers, err := service.GetAllSuppliers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取供应商列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    suppliers,
	})
}

// GetSupplier 获取供应商详情
// GET /api/suppliers/:id
func GetSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	supplier, err := service.GetSupplier(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "供应商不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    supplier,
	})
}

// CreateSupplier 创建供应商
// POST /api/suppliers
func CreateSupplier(c *gin.Context) {
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.CreateSupplier(&supplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建供应商失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    supplier,
	})
}

// UpdateSupplier 更新供应商
// PUT /api/suppliers/:id
func UpdateSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.UpdateSupplier(uint(id), &supplier); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新供应商失败: " + err.Error(),
		})
		return
	}

	updatedSupplier, _ := service.GetSupplier(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    updatedSupplier,
	})
}

// DeleteSupplier 删除供应商
// DELETE /api/suppliers/:id
func DeleteSupplier(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	if err := service.DeleteSupplier(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除供应商失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// ImportSuppliers CSV导入供应商
// POST /api/suppliers/import
func ImportSuppliers(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传CSV文件",
		})
		return
	}
	defer file.Close()

	count, err := service.ImportSuppliers(file)
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

// ExportSuppliers CSV导出供应商
// GET /api/suppliers/export
func ExportSuppliers(c *gin.Context) {
	suppliers, err := service.GetAllSuppliers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=suppliers.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入 BOM 防止中文乱码
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 写入表头
	writer.Write(service.SupplierCSVHeader())

	// 写入数据
	for i := range suppliers {
		writer.Write(service.SupplierToCSVRow(&suppliers[i]))
	}
}

// SupplierTemplate 下载CSV模板
// GET /api/suppliers/template
func SupplierTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=suppliers_template.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入 BOM
	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 写入表头
	writer.Write(service.SupplierCSVHeader())

	// 写入示例行
	writer.Write([]string{"示例供应商有限公司", "张经理", "13800138000", "example@test.com", "浙江省义乌市", "https://example.com", "示例备注"})
}
