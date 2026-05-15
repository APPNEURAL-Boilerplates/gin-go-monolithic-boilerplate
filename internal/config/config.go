package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName         string
	AppEnv          string
	Port            string
	GinMode         string
	LogLevel        slog.Level
	AllowedOrigins  []string
	ShutdownTimeout time.Duration
}

func Default() Config {
	return Config{
		AppName:         "gin-monolithic-boilerplate",
		AppEnv:          "local",
		Port:            "8080",
		GinMode:         "debug",
		LogLevel:        slog.LevelInfo,
		AllowedOrigins:  []string{"*"},
		ShutdownTimeout: 10 * time.Second,
	}
}

func Load() (Config, error) {
	cfg := Default()

	cfg.AppName = getEnv("APP_NAME", cfg.AppName)
	cfg.AppEnv = getEnv("APP_ENV", cfg.AppEnv)
	cfg.Port = getEnv("PORT", cfg.Port)
	cfg.GinMode = getEnv("GIN_MODE", cfg.GinMode)

	logLevel, err := parseLogLevel(getEnv("LOG_LEVEL", "info"))
	if err != nil {
		return cfg, err
	}
	cfg.LogLevel = logLevel

	cfg.AllowedOrigins = splitAndTrim(getEnv("CORS_ALLOWED_ORIGINS", "*"))

	shutdownSeconds, err := strconv.Atoi(getEnv("SHUTDOWN_TIMEOUT_SECONDS", "10"))
	if err != nil || shutdownSeconds <= 0 {
		return cfg, fmt.Errorf("SHUTDOWN_TIMEOUT_SECONDS must be a positive integer")
	}
	cfg.ShutdownTimeout = time.Duration(shutdownSeconds) * time.Second

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitAndTrim(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}

func parseLogLevel(value string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info", "":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unsupported LOG_LEVEL %q", value)
	}
}
