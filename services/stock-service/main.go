package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "stock-service/config"
    "stock-service/controllers"
    "stock-service/middleware"
    "stock-service/models"
    "stock-service/services"
    "syscall"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "go.uber.org/zap"
)

// @title Service Stock GMAO API
// @version 1.0
// @description API de gestion du stock de pièces détachées pour le système GMAO d'ICS
// @termsOfService http://ics.sn/terms/

// @contact.name ICS GMAO Team
// @contact.email gmao@ics.sn

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8004
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer Token (format: Bearer {token})

func main() {
    // Configuration des logs
    logger := config.InitLogger()
    defer logger.Sync()

    // Chargement de la configuration
    cfg := config.Load()
    
    // Initialisation de Redis
    redisClient := config.InitRedis(cfg)
    defer redisClient.Close()

    // Test de connexion Redis
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := redisClient.Ping(ctx).Err(); err != nil {
        logger.Fatal("Échec de connexion à Redis", zap.Error(err))
    }
    logger.Info("✅ Connexion Redis établie")

    // Initialisation des services
    stockService := services.NewStockService(redisClient, logger)
    
    // Insertion de données de test
    if err := insertTestData(stockService); err != nil {
        logger.Error("Erreur lors de l'insertion des données de test", zap.Error(err))
    }

    // Configuration de Gin
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    router := gin.New()
    
    // Middlewares globaux
    router.Use(gin.Recovery())
    router.Use(middleware.LoggerMiddleware(logger))
    router.Use(middleware.CORSMiddleware())

    // Configuration CORS plus détaillée si nécessaire
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"}, // En production: spécifier les domaines
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Routes publiques
    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "service":     "Service Stock GMAO",
            "version":     "1.0.0",
            "status":      "operational",
            "technology":  "Go Gin + Redis",
            "docs":        "/swagger/index.html",
            "health":      "/health",
        })
    })

    router.GET("/health", func(c *gin.Context) {
        // Test de connexion Redis
        ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        
        redisStatus := "connected"
        if err := redisClient.Ping(ctx).Err(); err != nil {
            redisStatus = "disconnected"
        }
        
        status := http.StatusOK
        if redisStatus == "disconnected" {
            status = http.StatusServiceUnavailable
        }
        
        c.JSON(status, gin.H{
            "status":    map[string]string{"healthy": "healthy", "unhealthy": "unhealthy"}[func() string {
                if status == http.StatusOK { return "healthy" }
                return "unhealthy"
            }()],
            "service":   "stock-service",
            "timestamp": time.Now().Format(time.RFC3339),
            "redis":     redisStatus,
        })
    })

    // Documentation Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Initialisation du contrôleur
    stockController := controllers.NewStockController(stockService, logger)


    // Routes API avec authentification
    apiRoutes := router.Group("/api")
    apiRoutes.Use(middleware.AuthMiddleware(cfg.JWTSecret))
    {
        // Routes pour les pièces détachées
        stock := apiRoutes.Group("/stock")
        {
            stock.GET("", stockController.GetAllPieces)
            stock.POST("", stockController.CreatePiece)
            stock.GET("/:id", stockController.GetPiece)
            stock.PUT("/:id", stockController.UpdatePiece)
            stock.DELETE("/:id", stockController.DeletePiece)
            stock.POST("/:id/increment", stockController.IncrementStock)
            stock.POST("/:id/decrement", stockController.DecrementStock)
            stock.GET("/alerts", stockController.GetLowStockAlerts)
            stock.GET("/search", stockController.SearchPieces)
        }
    }

    // Configuration du serveur
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Démarrage du serveur dans une goroutine
    go func() {
        logger.Info("🚀 Service Stock démarré", 
            zap.String("port", cfg.Port),
            zap.String("env", cfg.Environment))
        logger.Info("📚 Documentation disponible sur http://localhost:" + cfg.Port + "/swagger/index.html")
        
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Échec du démarrage du serveur", zap.Error(err))
        }
    }()

    // Attendre le signal d'arrêt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Arrêt du serveur en cours...")

    // Arrêt gracieux du serveur
    ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Arrêt forcé du serveur", zap.Error(err))
    }
    
    logger.Info("✅ Serveur arrêté proprement")
}

// insertTestData insère des données de test dans Redis
func insertTestData(stockService *services.StockService) error {
    testPieces := []models.Piece{
        {
            ID:              "piece-001",
            Nom:             "Roulement à billes 6205",
            Description:     "Roulement à billes standard pour moteurs électriques",
            Quantite:        25,
            SeuilMin:        5,
            PrixUnitaire:    45.50,
            Fournisseur:     "SKF Sénégal",
            Emplacement:     "A1-B2-C3",
            CodeEAN:         "3276000123456",
            Categorie:       "Roulements",
            UniteStock:      "pièce",
        },
        {
            ID:              "piece-002", 
            Nom:             "Courroie trapézoïdale A50",
            Description:     "Courroie trapézoïdale pour transmission de puissance",
            Quantite:        8,
            SeuilMin:        10,
            PrixUnitaire:    22.75,
            Fournisseur:     "Gates Dakar",
            Emplacement:     "A2-B1-C4",
            CodeEAN:         "3276000234567",
            Categorie:       "Courroies",
            UniteStock:      "pièce",
        },
        {
            ID:              "piece-003",
            Nom:             "Huile hydraulique ISO 68",
            Description:     "Huile hydraulique pour systèmes industriels",
            Quantite:        120,
            SeuilMin:        30,
            PrixUnitaire:    8.90,
            Fournisseur:     "Total Sénégal",
            Emplacement:     "B1-A3-C2",
            CodeEAN:         "3276000345678",
            Categorie:       "Lubrifiants",
            UniteStock:      "litre",
        },
        {
            ID:              "piece-004",
            Nom:             "Contacteur LC1D18",
            Description:     "Contacteur triphasé 18A pour commande moteur",
            Quantite:        3,
            SeuilMin:        8,
            PrixUnitaire:    125.00,
            Fournisseur:     "Schneider Electric",
            Emplacement:     "C1-A2-B1",
            CodeEAN:         "3276000456789",
            Categorie:       "Électrique",
            UniteStock:      "pièce",
        },
        {
            ID:              "piece-005",
            Nom:             "Joint torique NBR 20x3",
            Description:     "Joint d'étanchéité en caoutchouc nitrile",
            Quantite:        150,
            SeuilMin:        25,
            PrixUnitaire:    2.30,
            Fournisseur:     "Parker Hannifin",
            Emplacement:     "A3-B3-C1",
            CodeEAN:         "3276000567890",
            Categorie:       "Joints",
            UniteStock:      "pièce",
        },
    }

    for _, piece := range testPieces {
        if err := stockService.CreatePiece(&piece); err != nil {
            return err
        }
    }

    return nil
}