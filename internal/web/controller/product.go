package controller

import (
	"doudian/internal/database/model"
	"doudian/internal/service"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListProducts 获取商品列表（分页+供应商筛选+搜索）
// GET /api/products
func ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	supplierID, _ := strconv.ParseUint(c.Query("supplier_id"), 10, 32)
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := service.ListProducts(page, pageSize, uint(supplierID), keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取商品列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    result,
	})
}

// GetProduct 获取商品详情
// GET /api/products/:id
func GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	product, err := service.GetProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    product,
	})
}

// GetProductBySKU 根据SKU查询商品
// GET /api/products/lookup/:sku
func GetProductBySKU(c *gin.Context) {
	sku := c.Param("sku")

	product, err := service.GetProductBySKU(sku)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    product,
	})
}

// GetProductsBySupplier 获取供应商的商品列表
// GET /api/products/supplier/:supplier_id
func GetProductsBySupplier(c *gin.Context) {
	supplierIDStr := c.Param("supplier_id")
	supplierID, err := strconv.ParseUint(supplierIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的供应商ID",
		})
		return
	}

	products, err := service.GetProductsBySupplier(uint(supplierID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取商品列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    products,
	})
}

// CreateProduct 创建商品
// POST /api/products
func CreateProduct(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    product,
	})
}

// UpdateProduct 更新商品
// PUT /api/products/:id
func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := service.UpdateProduct(uint(id), &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新商品失败: " + err.Error(),
		})
		return
	}

	updatedProduct, _ := service.GetProduct(uint(id))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    updatedProduct,
	})
}

// DeleteProduct 删除商品
// DELETE /api/products/:id
func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的ID",
		})
		return
	}

	if err := service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除商品失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// ImportProducts CSV导入商品
// POST /api/products/import
func ImportProducts(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传CSV文件",
		})
		return
	}
	defer file.Close()

	count, err := service.ImportProducts(file)
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

// ExportProducts CSV导出商品
// GET /api/products/export
func ExportProducts(c *gin.Context) {
	products, err := service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "导出失败: " + err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=products.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.ProductCSVHeader())

	for i := range products {
		writer.Write(service.ProductToCSVRow(&products[i]))
	}
}

// ProductTemplate 下载CSV模板
// GET /api/products/template
func ProductTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=products_template.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.ProductCSVHeader())
	writer.Write([]string{"1", "SKU-001", "示例商品", "默认规格", "9.90", "19.90", "100", "https://example.com/image.jpg", "示例备注"})
}
