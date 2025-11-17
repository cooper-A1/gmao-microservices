
package config

import (
    "os"
    "github.com/go-redis/redis/v8"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

type Config struct {
    Port        string
    Environment string
    RedisURL    string
    JWTSecret   string
}

func Load() *Config {
    return &Config{
        Port:        getEnv("PORT", "8004"),
        Environment: getEnv("ENVIRONMENT", "development"),
        RedisURL:    getEnv("REDIS_URL", "redis://redis-stock:6379"),
        JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-here"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func InitRedis(cfg *Config) *redis.Client {
    opts, err := redis.ParseURL(cfg.RedisURL)
    if err != nil {
        // Fallback vers configuration manuelle
        opts = &redis.Options{
            Addr:     "redis-stock:6379",
            Password: "", 
            DB:       0,  
        }
    }
    
    return redis.NewClient(opts)
}

func InitLogger() *zap.Logger {
    config := zap.NewProductionConfig()
    
    // Configuration pour le développement
    if os.Getenv("ENVIRONMENT") != "production" {
        config = zap.NewDevelopmentConfig()
        config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
    }
    
    logger, err := config.Build()
    if err != nil {
        panic(err)
    }
    
    return logger
}