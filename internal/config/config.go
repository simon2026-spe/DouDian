package config

import (
	"os"
	"sync"
)

// Config 应用配置
type Config struct {
	Host          string
	Port          string
	DBPath        string
	SecretKey     string
	SecretPath    string
	JWTSecret     string
	LoginUsername string
	LoginPassword string
}

var (
	instance *Config
	once     sync.Once
)

// Load 加载配置（单例模式）
func Load() *Config {
	once.Do(func() {
		instance = &Config{
			Host:          getEnv("HOST", "0.0.0.0"),
			Port:          getEnv("PORT", "2095"),
			DBPath:        getEnv("DB_PATH", "./instance/doudian.db"),
			SecretKey:     getEnv("SECRET_KEY", "doudian-secret-key-change-in-production"),
			SecretPath:    getEnv("SECRET_PATH", ""),
			JWTSecret:     getEnv("JWT_SECRET", "doudian-jwt-secret-change-in-production"),
			LoginUsername: getEnv("LOGIN_USERNAME", "admin"),
			LoginPassword: getEnv("LOGIN_PASSWORD", "admin123"),
		}
	})
	return instance
}

// Get 获取配置单例
func Get() *Config {
	return Load()
}

// getEnv 读取环境变量，带默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
