package main

import (
	"doudian/internal/config"
	"doudian/internal/database"
	"doudian/internal/eventbus"
	"doudian/internal/web"
	"log"
)

func main() {
	// 加载配置
	cfg := config.Get()
	log.Printf("Configuration loaded: Host=%s Port=%s DBPath=%s", cfg.Host, cfg.Port, cfg.DBPath)

	// 初始化事件总线
	eventbus.GetBus()
	log.Println("EventBus initialized")

	// 初始化数据库
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 设置 Gin 路由
	r := web.SetupRouter()

	// 构建监听地址
	addr := cfg.Host + ":" + cfg.Port
	log.Printf("Server starting on %s", addr)

	if cfg.SecretPath != "" {
		log.Printf("SecretPath enabled: /%s", cfg.SecretPath)
	}

	// 启动 HTTP 服务器
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
