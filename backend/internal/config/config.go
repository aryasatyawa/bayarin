package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type JWTConfig struct {
	Secret       string
	ExpireHours  int
	RefreshHours int
}

type AppConfig struct {
	Name              string
	Version           string
	Currency          string
	CurrencyMinorUnit int // 100 untuk IDR (sen)
}

func Load() (*Config, error) {
	// Load .env file (ignore error jika tidak ada, untuk production bisa pakai env vars langsung)
	_ = godotenv.Load()

	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	maxLifetime, _ := strconv.Atoi(getEnv("DB_MAX_LIFETIME", "300"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	redisPoolSize, _ := strconv.Atoi(getEnv("REDIS_POOL_SIZE", "10"))
	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	currencyMinor, _ := strconv.Atoi(getEnv("CURRENCY_MINOR_UNIT", "100"))

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "bayarin_user"),
			Password:     getEnv("DB_PASSWORD", "bayarin_pass_2024"),
			Name:         getEnv("DB_NAME", "bayarin_db"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns: maxOpenConns,
			MaxIdleConns: maxIdleConns,
			MaxLifetime:  time.Duration(maxLifetime) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
			PoolSize: redisPoolSize,
		},
		JWT: JWTConfig{
			Secret:       getEnv("JWT_SECRET", "bayarin-secret-key"),
			ExpireHours:  jwtExpire,
			RefreshHours: jwtExpire * 7, // 7x JWT expire
		},
		App: AppConfig{
			Name:              getEnv("APP_NAME", "Bayarin"),
			Version:           getEnv("APP_VERSION", "1.0.0"),
			Currency:          getEnv("CURRENCY", "IDR"),
			CurrencyMinorUnit: currencyMinor,
		},
	}

	return cfg, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
