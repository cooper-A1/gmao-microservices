package controllers

import (
    "net/http"
    "stock-service/models"
    "stock-service/services"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type StockController struct {
    stockService *services.StockService
    logger       *zap.Logger
}

func NewStockController(stockService *services.StockService, logger *zap.Logger) *StockController {
    return &StockController{
        stockService: stockService,
        logger:       logger,
    }
}

// GetAllPieces récupère toutes les pièces en stock
// @Summary Récupérer toutes les pièces
// @Description Retourne la liste complète des pièces détachées en stock
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Liste des pièces"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock [get]
func (sc *StockController) GetAllPieces(c *gin.Context) {
    pieces, err := sc.stockService.GetAllPieces()
    if err != nil {
        sc.logger.Error("Erreur lors de la récupération des pièces", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la récupération des pièces",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Pièces récupérées avec succès",
        "data": pieces,
        "count": len(pieces),
    })
}

// CreatePiece crée une nouvelle pièce détachée
// @Summary Créer une nouvelle pièce
// @Description Ajoute une nouvelle pièce détachée au stock
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param piece body models.CreatePieceRequest true "Données de la pièce"
// @Success 201 {object} map[string]interface{} "Pièce créée"
// @Failure 400 {object} map[string]interface{} "Données invalides"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock [post]
func (sc *StockController) CreatePiece(c *gin.Context) {
    var req models.CreatePieceRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        sc.logger.Warn("Données invalides pour création de pièce", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Données invalides",
            "details": err.Error(),
        })
        return
    }

    // Conversion vers le modèle Piece
    piece := &models.Piece{
        Nom:          req.Nom,
        Description:  req.Description,
        Quantite:     req.Quantite,
        SeuilMin:     req.SeuilMin,
        PrixUnitaire: req.PrixUnitaire,
        Fournisseur:  req.Fournisseur,
        Emplacement:  req.Emplacement,
        CodeEAN:      req.CodeEAN,
        Categorie:    req.Categorie,
        UniteStock:   req.UniteStock,
    }

    if err := sc.stockService.CreatePiece(piece); err != nil {
        sc.logger.Error("Erreur lors de la création de la pièce", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la création de la pièce",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Pièce créée avec succès",
        "data": piece,
    })
}

// GetPiece récupère une pièce par ID
// @Summary Récupérer une pièce par ID
// @Description Retourne les détails d'une pièce spécifique
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la pièce"
// @Success 200 {object} map[string]interface{} "Détails de la pièce"
// @Failure 404 {object} map[string]interface{} "Pièce non trouvée"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/{id} [get]
func (sc *StockController) GetPiece(c *gin.Context) {
    id := c.Param("id")

    piece, err := sc.stockService.GetPiece(id)
    if err != nil {
        if err.Error() == "pièce non trouvée: "+id {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Pièce non trouvée",
                "piece_id": id,
            })
            return
        }

        sc.logger.Error("Erreur lors de la récupération de la pièce", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la récupération de la pièce",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Pièce trouvée",
        "data": piece,
    })
}

// UpdatePiece met à jour une pièce existante
// @Summary Mettre à jour une pièce
// @Description Met à jour les informations d'une pièce existante
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la pièce"
// @Param piece body models.UpdatePieceRequest true "Données à mettre à jour"
// @Success 200 {object} map[string]interface{} "Pièce mise à jour"
// @Failure 400 {object} map[string]interface{} "Données invalides"
// @Failure 404 {object} map[string]interface{} "Pièce non trouvée"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/{id} [put]
func (sc *StockController) UpdatePiece(c *gin.Context) {
    id := c.Param("id")
    var req models.UpdatePieceRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        sc.logger.Warn("Données invalides pour mise à jour de pièce", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Données invalides",
            "details": err.Error(),
        })
        return
    }

    piece, err := sc.stockService.UpdatePiece(id, &req)
    if err != nil {
        if err.Error() == "pièce non trouvée: "+id {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Pièce non trouvée",
                "piece_id": id,
            })
            return
        }

        sc.logger.Error("Erreur lors de la mise à jour de la pièce", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la mise à jour de la pièce",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Pièce mise à jour avec succès",
        "data": piece,
    })
}

// DeletePiece supprime une pièce
// @Summary Supprimer une pièce
// @Description Supprime une pièce du stock
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la pièce"
// @Success 204 "Pièce supprimée"
// @Failure 404 {object} map[string]interface{} "Pièce non trouvée"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/{id} [delete]
func (sc *StockController) DeletePiece(c *gin.Context) {
    id := c.Param("id")

    if err := sc.stockService.DeletePiece(id); err != nil {
        if err.Error() == "pièce non trouvée: "+id {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Pièce non trouvée",
                "piece_id": id,
            })
            return
        }

        sc.logger.Error("Erreur lors de la suppression de la pièce", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la suppression de la pièce",
            "details": err.Error(),
        })
        return
    }

    c.Status(http.StatusNoContent)
}

// IncrementStock augmente la quantité d'une pièce
// @Summary Incrémenter le stock
// @Description Augmente la quantité en stock d'une pièce
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la pièce"
// @Param movement body models.StockMovementRequest true "Données du mouvement"
// @Success 200 {object} map[string]interface{} "Stock incrémenté"
// @Failure 400 {object} map[string]interface{} "Données invalides"
// @Failure 404 {object} map[string]interface{} "Pièce non trouvée"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/{id}/increment [post]
func (sc *StockController) IncrementStock(c *gin.Context) {
    id := c.Param("id")
    var req models.StockMovementRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        sc.logger.Warn("Données invalides pour incrémentation de stock", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Données invalides",
            "details": err.Error(),
        })
        return
    }

    piece, err := sc.stockService.IncrementStock(id, req.Quantite, req.Motif)
    if err != nil {
        if err.Error() == "pièce non trouvée: "+id {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Pièce non trouvée",
                "piece_id": id,
            })
            return
        }

        sc.logger.Error("Erreur lors de l'incrémentation du stock", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de l'incrémentation du stock",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Stock incrémenté avec succès",
        "data": piece,
        "mouvement": gin.H{
            "type": "increment",
            "quantite": req.Quantite,
            "motif": req.Motif,
        },
    })
}

// DecrementStock diminue la quantité d'une pièce
// @Summary Décrémenter le stock
// @Description Diminue la quantité en stock d'une pièce
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID de la pièce"
// @Param movement body models.StockMovementRequest true "Données du mouvement"
// @Success 200 {object} map[string]interface{} "Stock décrémenté"
// @Failure 400 {object} map[string]interface{} "Données invalides ou stock insuffisant"
// @Failure 404 {object} map[string]interface{} "Pièce non trouvée"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/{id}/decrement [post]
func (sc *StockController) DecrementStock(c *gin.Context) {
    id := c.Param("id")
    var req models.StockMovementRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        sc.logger.Warn("Données invalides pour décrémentation de stock", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Données invalides",
            "details": err.Error(),
        })
        return
    }

    piece, err := sc.stockService.DecrementStock(id, req.Quantite, req.Motif)
    if err != nil {
        if err.Error() == "pièce non trouvée: "+id {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "Pièce non trouvée",
                "piece_id": id,
            })
            return
        }

        // Vérification si c'est un problème de stock insuffisant
        if err.Error()[:17] == "stock insuffisant" {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Stock insuffisant",
                "details": err.Error(),
            })
            return
        }

        sc.logger.Error("Erreur lors de la décrémentation du stock", zap.String("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la décrémentation du stock",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Stock décrémenté avec succès",
        "data": piece,
        "mouvement": gin.H{
            "type": "decrement",
            "quantite": req.Quantite,
            "motif": req.Motif,
        },
    })
}

// GetLowStockAlerts récupère les alertes de stock faible
// @Summary Récupérer les alertes de stock
// @Description Retourne les pièces en stock faible ou critique
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Alertes de stock"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/alerts [get]
func (sc *StockController) GetLowStockAlerts(c *gin.Context) {
    alerts, err := sc.stockService.GetLowStockAlerts()
    if err != nil {
        sc.logger.Error("Erreur lors de la récupération des alertes", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la récupération des alertes",
            "details": err.Error(),
        })
        return
    }

    // Comptage par sévérité
    critiques := 0
    attentions := 0
    for _, alert := range alerts {
        if alert.Severite == "critique" {
            critiques++
        } else {
            attentions++
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Alertes de stock récupérées",
        "data": alerts,
        "summary": gin.H{
            "total": len(alerts),
            "critiques": critiques,
            "attentions": attentions,
        },
    })
}

// SearchPieces recherche des pièces
// @Summary Rechercher des pièces
// @Description Recherche des pièces par nom, description ou catégorie
// @Tags Stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Terme de recherche"
// @Success 200 {object} map[string]interface{} "Résultats de recherche"
// @Failure 400 {object} map[string]interface{} "Paramètre de recherche manquant"
// @Failure 500 {object} map[string]interface{} "Erreur interne"
// @Router /stock/search [get]
func (sc *StockController) SearchPieces(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Paramètre de recherche 'q' requis",
        })
        return
    }

    pieces, err := sc.stockService.SearchPieces(query)
    if err != nil {
        sc.logger.Error("Erreur lors de la recherche", zap.String("query", query), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Erreur lors de la recherche",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Recherche effectuée avec succès",
        "query": query,
        "data": pieces,
        "count": len(pieces),
    })
}