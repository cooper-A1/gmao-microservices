
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// LoggerMiddleware crée un middleware de logging avec zap
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return gin.LoggerWithWriter(gin.DefaultWriter, "/health") // Exclure /health des logs
}

// ZapLoggerMiddleware est un middleware de logging plus détaillé
func ZapLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Traitement de la requête
        c.Next()

        // Calcul du temps de traitement
        latency := time.Since(start)

        // Construction du chemin complet
        if raw != "" {
            path = path + "?" + raw
        }

        // Log de la requête
        logger.Info("Request processed",
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("ip", c.ClientIP()),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.Int("body_size", c.Writer.Size()),
        )
    }
}