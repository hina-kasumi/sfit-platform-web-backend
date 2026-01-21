package config

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	App          AppConfig
	Cors         CorsConfig
	Database     DatabaseConfig
	Redis        RedisConfig
	Admin        AdminConfig
	Youtube      YoutubeConfig
	Jwt          JwtConfig
	RefreshToken RefreshTokenConfig
}

type AppConfig struct {
	Port string
}

type CorsConfig struct {
	AllowedOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type AdminConfig struct {
	Username string
	Email    string
	Password string
}

type YoutubeConfig struct {
	APIKey string
	Part   string
	URL    string
}

type JwtConfig struct {
	SecretKey     string
	ExpirationSec int
}

type RefreshTokenConfig struct {
	SecretKey     string
	ExpirationSec int
}

var (
	cfg  *Config
	once sync.Once
)

func Load() (*Config, error) {
	once.Do(func() {
		cfg = &Config{
			App: AppConfig{
				Port: getEnv("APP_PORT", "8080"),
			},
			Cors: CorsConfig{
				AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", ","),
			},
			Database: DatabaseConfig{
				Host:     getEnv("DB_HOST", "localhost"),
				Port:     getEnvAsInt("DB_PORT", 5432),
				User:     getEnv("DB_USER", "user"),
				Password: getEnv("DB_PASSWORD", "password"),
				Name:     getEnv("DB_NAME", "sfitdb"),
			},
			Redis: RedisConfig{
				Address:  getEnv("REDIS_ADDRESS", "localhost:6379"),
				Password: getEnv("REDIS_PASSWORD", ""),
				DB:       getEnvAsInt("REDIS_DB", 0),
			},
			Admin: AdminConfig{
				Username: getEnv("INIT_ADMIN_USERNAME", "admin"),
				Email:    getEnv("INIT_ADMIN_EMAIL", "admin@example.com"),
				Password: getEnv("INIT_ADMIN_PASSWORD", "adminpass"),
			},
			Youtube: YoutubeConfig{
				APIKey: getEnv("YOUTUBE_API_KEY", ""),
				Part:   getEnv("YOUTUBE_PART", "snippet"),
				URL:    getEnv("YOUTUBE_URL", "https://www.googleapis.com/youtube/v3"),
			},
			Jwt: JwtConfig{
				SecretKey:     getEnv("JWT_SECRET_KEY", "secret"),
				ExpirationSec: getEnvAsInt("JWT_EXPIRATION_SEC", 3600),
			},
			RefreshToken: RefreshTokenConfig{
				SecretKey:     getEnv("REFRESH_TOKEN_SECRET_KEY", "refreshsecret"),
				ExpirationSec: getEnvAsInt("REFRESH_TOKEN_EXPIRATION_SEC", 86400),
			},
		}
	})
	return cfg, nil
}

// Helper functions to read environment variables with default values
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsSlice(key, sep string) []string {
	val := getEnv(key, "")
	if val == "" {
		return []string{}
	}

	items := strings.Split(val, sep)

	result := make([]string, 0, len(items))
	for _, item := range items {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
