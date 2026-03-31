package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	JWTSecret                string
	ServerPort               string
	MinioEndpoint            string
	MinioAccessKey           string
	MinioSecretKey           string
	MinioBucket              string
	MinioUseSSL              bool
	CleanupSoftDeleteDays    int
	CleanupCartStaleDays     int
	CleanupOrderArchiveDays  int
	CleanupChatRetentionDays int
	CleanupNotifReadDays     int
	CleanupNotifUnreadDays   int
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = Config{
		DBHost:                   getEnv("DB_HOST", "localhost"),
		DBPort:                   getEnv("DB_PORT", "3306"),
		DBUser:                   getEnv("DB_USER", "root"),
		DBPassword:               getEnv("DB_PASSWORD", ""),
		DBName:                   getEnv("DB_NAME", "bycigar"),
		JWTSecret:                getEnv("JWT_SECRET", "your-super-secret-jwt-key"),
		ServerPort:               getEnv("SERVER_PORT", "3000"),
		MinioEndpoint:            getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey:           getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:           getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		MinioBucket:              getEnv("MINIO_BUCKET", "bycigar"),
		MinioUseSSL:              getEnvBool("MINIO_USE_SSL", false),
		CleanupSoftDeleteDays:    getEnvInt("CLEANUP_SOFT_DELETE_DAYS", 30),
		CleanupCartStaleDays:     getEnvInt("CLEANUP_CART_STALE_DAYS", 90),
		CleanupOrderArchiveDays:  getEnvInt("CLEANUP_ORDER_ARCHIVE_DAYS", 365),
		CleanupChatRetentionDays: getEnvInt("CLEANUP_CHAT_RETENTION_DAYS", 30),
		CleanupNotifReadDays:     getEnvInt("CLEANUP_NOTIF_READ_DAYS", 60),
		CleanupNotifUnreadDays:   getEnvInt("CLEANUP_NOTIF_UNREAD_DAYS", 120),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err == nil {
			return b
		}
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err == nil && n > 0 {
			return n
		}
	}
	return defaultValue
}
