package server

import (
	"fmt"
	"os"
	"strconv"
)

// Config はサーバーの設定を表します
type Config struct {
	// サーバー設定
	Port            int
	ReadTimeoutSec  int
	WriteTimeoutSec int
	IdleTimeoutSec  int

	// データベース設定
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// JWT設定
	JWTSecret       string
	JWTExpiryHours  int
	JWTRefreshHours int

	// CORS設定
	AllowedOrigins []string
}

// NewConfig は環境変数から設定を読み込みます
func NewConfig() *Config {
	return &Config{
		// サーバー設定
		Port:            getEnvAsInt("SERVER_PORT", 8080),
		ReadTimeoutSec:  getEnvAsInt("SERVER_READ_TIMEOUT", 15),
		WriteTimeoutSec: getEnvAsInt("SERVER_WRITE_TIMEOUT", 15),
		IdleTimeoutSec:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 60),

		// データベース設定
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "task_manager"),

		// JWT設定
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiryHours:  getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		JWTRefreshHours: getEnvAsInt("JWT_REFRESH_HOURS", 168), // 7日

		// CORS設定
		AllowedOrigins: []string{getEnv("CORS_ALLOWED_ORIGINS", "*")},
	}
}

// GetServerAddr はサーバーのアドレスを返します
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf(":%d", c.Port)
}

// getEnvAsInt は環境変数の値を整数として取得します
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getEnv は環境変数の値を取得します
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
