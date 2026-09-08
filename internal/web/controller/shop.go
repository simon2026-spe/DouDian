package controller

import (
	"doudian/internal/database/model"
	"doudian/internal/service"
	"encoding/csv"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListShops(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := service.ListShops(page, pageSize, keyword, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取抖店列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "data": result})
}

func GetAllShops(c *gin.Context) {
	shops, err := service.GetAllShops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取抖店列表失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "data": shops})
}

func GetShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	shop, err := service.GetShop(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "抖店不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "data": shop})
}

func CreateShop(c *gin.Context) {
	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	if err := service.CreateShop(&shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "创建抖店失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "创建成功", "data": shop})
}

func UpdateShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	var shop model.Shop
	if err := c.ShouldBindJSON(&shop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误: " + err.Error()})
		return
	}

	if err := service.UpdateShop(uint(id), &shop); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "更新抖店失败: " + err.Error()})
		return
	}

	updatedShop, _ := service.GetShop(uint(id))
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "更新成功", "data": updatedShop})
}

func DeleteShop(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "无效的ID"})
		return
	}

	if err := service.DeleteShop(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除抖店失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功"})
}

func ImportShops(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请上传CSV文件"})
		return
	}
	defer file.Close()

	count, err := service.ImportShops(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "导入失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "导入成功", "data": gin.H{"count": count}})
}

func ExportShops(c *gin.Context) {
	shops, err := service.GetAllShops()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "导出失败: " + err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=shops.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.ShopCSVHeader())
	for i := range shops {
		writer.Write(service.ShopToCSVRow(&shops[i]))
	}
}

func ShopTemplate(c *gin.Context) {
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=shops_template.csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	c.Writer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer.Write(service.ShopCSVHeader())
	writer.Write([]string{"抖音旗舰店", "SHOP001", "app_key_example", "app_secret_example", "active", "测试店铺"})
}
