package controller

import (
	"doudian/internal/config"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ProcessAllOrdersTool 工具页一键处理
func ProcessAllOrdersTool(c *gin.Context) {
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

// ListBackups 获取备份列表
// GET /api/tools/backups
func ListBackups(c *gin.Context) {
	cfg := config.Get()
	backupDir := filepath.Join(filepath.Dir(cfg.DBPath), "backups")

	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "暂无备份",
			"data":    []interface{}{},
		})
		return
	}

	type BackupInfo struct {
		ID        string `json:"id"`
		Filename  string `json:"filename"`
		Size      int64  `json:"size"`
		CreatedAt string `json:"created_at"`
	}

	var backups []BackupInfo
	for i, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		backups = append(backups, BackupInfo{
			ID:        strconv.Itoa(i + 1),
			Filename:  file.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    backups,
	})
}

// DownloadBackup 下载备份文件
// GET /api/tools/backup/:filename/download
func DownloadBackup(c *gin.Context) {
	filename := c.Param("filename")
	cfg := config.Get()
	backupPath := filepath.Join(filepath.Dir(cfg.DBPath), "backups", filename)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "备份文件不存在",
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.File(backupPath)
}

// DeleteBackup 删除备份文件
// DELETE /api/tools/backup/:id
func DeleteBackup(c *gin.Context) {
	idStr := c.Param("id")
	cfg := config.Get()
	backupDir := filepath.Join(filepath.Dir(cfg.DBPath), "backups")

	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取备份目录失败",
		})
		return
	}

	idx, err := strconv.Atoi(idStr)
	if err != nil || idx < 1 || idx > len(files) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的备份ID",
		})
		return
	}

	filename := files[idx-1].Name()
	backupPath := filepath.Join(backupDir, filename)

	if err := os.Remove(backupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除备份失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}
