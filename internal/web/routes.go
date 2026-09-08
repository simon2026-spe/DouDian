package web

import (
	"doudian/internal/config"
	"doudian/internal/web/controller"
	"doudian/internal/web/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置Gin路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 应用 CORS 中间件
	r.Use(middleware.CORS())

	// 获取配置
	cfg := config.Get()
	secretPath := strings.Trim(cfg.SecretPath, "/")

	// 基础路径组（考虑 SecretPath）
	var baseGroup *gin.RouterGroup
	if secretPath != "" {
		baseGroup = r.Group("/" + secretPath)
	} else {
		baseGroup = r.Group("/")
	}

	// API 路由组
	api := baseGroup.Group("/api")
	{
		// 认证相关（不需要登录）
		api.POST("/login", controller.Login)
		api.POST("/logout", controller.Logout)

		// 需要认证的 API
		authAPI := api.Group("")
		authAPI.Use(middleware.AuthRequired())
		{
			// 仪表盘
			dashboard := authAPI.Group("/dashboard")
			{
				dashboard.GET("/stats", controller.GetDashboardStats)
				dashboard.GET("/recent-orders", controller.GetRecentOrders)
			}

			// 抖店
			shops := authAPI.Group("/shops")
			{
				shops.GET("", controller.ListShops)
				shops.GET("/all", controller.GetAllShops)
				shops.GET("/:id", controller.GetShop)
				shops.POST("", controller.CreateShop)
				shops.PUT("/:id", controller.UpdateShop)
				shops.DELETE("/:id", controller.DeleteShop)
			}

			// 供应商
			suppliers := authAPI.Group("/suppliers")
			{
				suppliers.GET("", controller.ListSuppliers)
				suppliers.GET("/all", controller.GetAllSuppliers)
				suppliers.GET("/template", controller.SupplierTemplate)
				suppliers.GET("/export", controller.ExportSuppliers)
				suppliers.POST("/import", controller.ImportSuppliers)
				suppliers.GET("/:id", controller.GetSupplier)
				suppliers.POST("", controller.CreateSupplier)
				suppliers.PUT("/:id", controller.UpdateSupplier)
				suppliers.DELETE("/:id", controller.DeleteSupplier)
			}

			// 商品
			products := authAPI.Group("/products")
			{
				products.GET("", controller.ListProducts)
				products.GET("/template", controller.ProductTemplate)
				products.GET("/export", controller.ExportProducts)
				products.POST("/import", controller.ImportProducts)
				products.GET("/lookup/:sku", controller.GetProductBySKU)
				products.GET("/supplier/:supplier_id", controller.GetProductsBySupplier)
				products.GET("/:id", controller.GetProduct)
				products.POST("", controller.CreateProduct)
				products.PUT("/:id", controller.UpdateProduct)
				products.DELETE("/:id", controller.DeleteProduct)
			}

			// 订单
			orders := authAPI.Group("/orders")
			{
				orders.GET("", controller.ListOrders)
				orders.GET("/template", controller.OrderTemplate)
				orders.GET("/export", controller.ExportOrders)
				orders.POST("/import", controller.ImportOrders)
				orders.GET("/:id", controller.GetOrder)
				orders.POST("", controller.CreateOrder)
				orders.DELETE("/:id", controller.DeleteOrder)
			}

			// 采购单
			purchaseOrders := authAPI.Group("/purchase-orders")
			{
				purchaseOrders.GET("", controller.ListPurchaseOrders)
				purchaseOrders.GET("/export", controller.ExportPurchaseOrders)
				purchaseOrders.POST("/generate/:order_id", controller.GeneratePurchaseOrder)
				purchaseOrders.GET("/:id", controller.GetPurchaseOrder)
				purchaseOrders.POST("", controller.CreatePurchaseOrder)
				purchaseOrders.PUT("/:id", controller.UpdatePurchaseOrder)
			}

			// 工具
			tools := authAPI.Group("/tools")
			{
				tools.POST("/process-all", controller.ProcessAllOrders)
				tools.POST("/backup", controller.BackupDatabase)
				tools.POST("/backup/create", controller.BackupDatabase)
				tools.GET("/backups", controller.ListBackups)
				tools.GET("/backup/:filename/download", controller.DownloadBackup)
				tools.DELETE("/backup/:id", controller.DeleteBackup)
			}
		}
	}

	// 静态资源（前端构建产物的 assets 目录）
	baseGroup.StaticFS("/assets", http.Dir("./static/assets"))

	// SPA 前端路由回退：所有非 API 的 GET 请求都返回 index.html
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// API 请求返回 404 JSON
		if strings.Contains(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "API endpoint not found",
			})
			return
		}

		// 静态资源请求返回 404
		if strings.Contains(path, "/assets/") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Asset not found",
			})
			return
		}

		// 其他所有路径返回 index.html（SPA 前端路由）
		c.File("./static/index.html")
	})

	return r
}
