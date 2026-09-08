package controller

import (
	"doudian/internal/config"
	"doudian/internal/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// ProcessAllOrders 一键处理所有待处理订单
// POST /api/tools/process-all
func ProcessAllOrders(c *gin.Context) {
	successCount, failCount, err := service.ProcessAllPendingOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "处理失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("处理完成，成功 %d 个，失败 %d 个", successCount, failCount),
		"data": gin.H{
			"success_count": successCount,
			"fail_count":    failCount,
		},
	})
}

// BackupDatabase 数据库备份
// POST /api/tools/backup
func BackupDatabase(c *gin.Context) {
	cfg := config.Get()
	dbPath := cfg.DBPath

	// 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "数据库文件不存在",
		})
		return
	}

	// 创建备份目录
	backupDir := filepath.Join(filepath.Dir(dbPath), "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建备份目录失败: " + err.Error(),
		})
		return
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	backupFileName := fmt.Sprintf("doudian_backup_%s.db", timestamp)
	backupPath := filepath.Join(backupDir, backupFileName)

	// 读取原数据库文件
	data, err := os.ReadFile(dbPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取数据库失败: " + err.Error(),
		})
		return
	}

	// 写入备份文件
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "写入备份文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "备份成功",
		"data": gin.H{
			"backup_file": backupFileName,
			"backup_path": backupPath,
		},
	})
}
