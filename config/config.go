package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

type TokenConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ExpirationTime time.Duration
}

type JWTConfig struct {
	AccessToken  TokenConfig
	RefreshToken TokenConfig
	Algorithm    string
	Issuer       string
}

var (
	AppConfig Config
	DB        *gorm.DB
)

func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	AppConfig = Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "password"),
			DBName:       getEnv("DB_NAME", "golang"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxLifetime:  time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 60)) * time.Minute,
		},
		JWT: JWTConfig{
			Algorithm: getEnv("JWT_ALGORITHM", "RS256"),
			Issuer:    getEnv("JWT_ISSUER", "golang"),
			AccessToken: TokenConfig{
				PrivateKeyPath: getEnv("JWT_ACCESS_PRIVATE_KEY_PATH", "keys/access_private.pem"),
				PublicKeyPath:  getEnv("JWT_ACCESS_PUBLIC_KEY_PATH", "keys/access_public.pem"),
				ExpirationTime: time.Duration(getEnvAsInt("JWT_ACCESS_EXPIRATION_TIME", 15)) * time.Minute,
			},
			RefreshToken: TokenConfig{
				PrivateKeyPath: getEnv("JWT_REFRESH_PRIVATE_KEY_PATH", "keys/refresh_private.pem"),
				PublicKeyPath:  getEnv("JWT_REFRESH_PUBLIC_KEY_PATH", "keys/refresh_public.pem"),
				ExpirationTime: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRATION_TIME", 30)) * 24 * time.Hour, // 30 days
			},
		},
	}

	initDB()
}

func initDB() {
	var err error
	dsn := GetDSN(&AppConfig.Database)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to configure connection pool")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// CloseDB closes the database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("error getting database instance: %v", err)
	}
	return sqlDB.Close()
}

func GetDSN(c *DatabaseConfig) string {
	value := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
	return value
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
