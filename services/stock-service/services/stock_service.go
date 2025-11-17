package services

import (
    "context"
    "fmt"
    "stock-service/models"
    "strings"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
    "go.uber.org/zap"
)

const (
    PIECE_KEY_PREFIX    = "stock:piece:"
    PIECES_SET_KEY      = "stock:pieces"
    CATEGORY_SET_PREFIX = "stock:category:"
)

type StockService struct {
    redis  *redis.Client
    logger *zap.Logger
}

func NewStockService(redisClient *redis.Client, logger *zap.Logger) *StockService {
    return &StockService{
        redis:  redisClient,
        logger: logger,
    }
}

// CreatePiece crée une nouvelle pièce en stock
func (s *StockService) CreatePiece(piece *models.Piece) error {
    ctx := context.Background()

    // Génération d'un ID si non fourni
    if piece.ID == "" {
        piece.ID = uuid.New().String()
    }

    // Vérification de l'unicité de l'ID
    exists, err := s.redis.Exists(ctx, PIECE_KEY_PREFIX+piece.ID).Result()
    if err != nil {
        return fmt.Errorf("erreur lors de la vérification d'existence: %w", err)
    }
    if exists > 0 {
        return fmt.Errorf("une pièce avec l'ID %s existe déjà", piece.ID)
    }

    // Timestamps
    now := time.Now()
    piece.CreatedAt = now
    piece.UpdatedAt = now

    // Sérialisation
    pieceJSON, err := piece.ToJSON()
    if err != nil {
        return fmt.Errorf("erreur de sérialisation: %w", err)
    }

    // Transaction Redis
    pipe := s.redis.TxPipeline()
    
    // Stockage de la pièce
    pipe.Set(ctx, PIECE_KEY_PREFIX+piece.ID, pieceJSON, 0)
    
    // Ajout à l'ensemble des pièces
    pipe.SAdd(ctx, PIECES_SET_KEY, piece.ID)
    
    // Ajout à l'ensemble de la catégorie
    if piece.Categorie != "" {
        pipe.SAdd(ctx, CATEGORY_SET_PREFIX+strings.ToLower(piece.Categorie), piece.ID)
    }

    // Exécution de la transaction
    _, err = pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("erreur lors de la création: %w", err)
    }

    s.logger.Info("Pièce créée avec succès",
        zap.String("id", piece.ID),
        zap.String("nom", piece.Nom),
        zap.Int("quantite", piece.Quantite))

    return nil
}

// GetPiece récupère une pièce par ID
func (s *StockService) GetPiece(id string) (*models.Piece, error) {
    ctx := context.Background()

    pieceJSON, err := s.redis.Get(ctx, PIECE_KEY_PREFIX+id).Result()
    if err == redis.Nil {
        return nil, fmt.Errorf("pièce non trouvée: %s", id)
    }
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération: %w", err)
    }

    var piece models.Piece
    if err := piece.FromJSON([]byte(pieceJSON)); err != nil {
        return nil, fmt.Errorf("erreur de désérialisation: %w", err)
    }

    return &piece, nil
}

// GetAllPieces récupère toutes les pièces
func (s *StockService) GetAllPieces() ([]models.Piece, error) {
    ctx := context.Background()

    // Récupération de tous les IDs
    pieceIDs, err := s.redis.SMembers(ctx, PIECES_SET_KEY).Result()
    if err != nil {
        return nil, fmt.Errorf("erreur lors de la récupération des IDs: %w", err)
    }

    pieces := make([]models.Piece, 0, len(pieceIDs))

    // Récupération en batch si possible, sinon une par une
    for _, id := range pieceIDs {
        piece, err := s.GetPiece(id)
        if err != nil {
            s.logger.Warn("Impossible de récupérer la pièce", zap.String("id", id), zap.Error(err))
            continue
        }
        pieces = append(pieces, *piece)
    }

    return pieces, nil
}

// UpdatePiece met à jour une pièce existante
func (s *StockService) UpdatePiece(id string, updates *models.UpdatePieceRequest) (*models.Piece, error) {
    // Récupération de la pièce existante
    piece, err := s.GetPiece(id)
    if err != nil {
        return nil, err
    }

    // Application des mises à jour
    if updates.Nom != nil {
        piece.Nom = *updates.Nom
    }
    if updates.Description != nil {
        piece.Description = *updates.Description
    }
    if updates.SeuilMin != nil {
        piece.SeuilMin = *updates.SeuilMin
    }
    if updates.PrixUnitaire != nil {
        piece.PrixUnitaire = *updates.PrixUnitaire
    }
    if updates.Fournisseur != nil {
        piece.Fournisseur = *updates.Fournisseur
    }
    if updates.Emplacement != nil {
        piece.Emplacement = *updates.Emplacement
    }
    if updates.CodeEAN != nil {
        piece.CodeEAN = *updates.CodeEAN
    }
    if updates.Categorie != nil {
        piece.Categorie = *updates.Categorie
    }
    if updates.UniteStock != nil {
        piece.UniteStock = *updates.UniteStock
    }

    // Mise à jour du timestamp
    piece.UpdatedAt = time.Now()

    // Sauvegarde
    ctx := context.Background()
    pieceJSON, err := piece.ToJSON()
    if err != nil {
        return nil, fmt.Errorf("erreur de sérialisation: %w", err)
    }

    if err := s.redis.Set(ctx, PIECE_KEY_PREFIX+id, pieceJSON, 0).Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de la mise à jour: %w", err)
    }

    s.logger.Info("Pièce mise à jour avec succès",
        zap.String("id", piece.ID),
        zap.String("nom", piece.Nom))

    return piece, nil
}

// DeletePiece supprime une pièce
func (s *StockService) DeletePiece(id string) error {
    ctx := context.Background()

    // Récupération de la pièce pour obtenir la catégorie
    piece, err := s.GetPiece(id)
    if err != nil {
        return err
    }

    // Transaction Redis
    pipe := s.redis.TxPipeline()
    
    // Suppression de la pièce
    pipe.Del(ctx, PIECE_KEY_PREFIX+id)
    
    // Suppression de l'ensemble des pièces
    pipe.SRem(ctx, PIECES_SET_KEY, id)
    
    // Suppression de l'ensemble de la catégorie
    if piece.Categorie != "" {
        pipe.SRem(ctx, CATEGORY_SET_PREFIX+strings.ToLower(piece.Categorie), id)
    }

    // Exécution
    _, err = pipe.Exec(ctx)
    if err != nil {
        return fmt.Errorf("erreur lors de la suppression: %w", err)
    }

    s.logger.Info("Pièce supprimée avec succès",
        zap.String("id", id),
        zap.String("nom", piece.Nom))

    return nil
}

// IncrementStock augmente la quantité en stock
func (s *StockService) IncrementStock(id string, quantite int, motif string) (*models.Piece, error) {
    piece, err := s.GetPiece(id)
    if err != nil {
        return nil, err
    }

    oldQuantite := piece.Quantite
    piece.Quantite += quantite
    piece.UpdatedAt = time.Now()

    // Sauvegarde
    ctx := context.Background()
    pieceJSON, err := piece.ToJSON()
    if err != nil {
        return nil, fmt.Errorf("erreur de sérialisation: %w", err)
    }

    if err := s.redis.Set(ctx, PIECE_KEY_PREFIX+id, pieceJSON, 0).Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de la mise à jour: %w", err)
    }

    s.logger.Info("Stock incrémenté",
        zap.String("piece_id", id),
        zap.String("nom", piece.Nom),
        zap.Int("ancien_stock", oldQuantite),
        zap.Int("nouveau_stock", piece.Quantite),
        zap.Int("increment", quantite),
        zap.String("motif", motif))

    return piece, nil
}

// DecrementStock diminue la quantité en stock
func (s *StockService) DecrementStock(id string, quantite int, motif string) (*models.Piece, error) {
    piece, err := s.GetPiece(id)
    if err != nil {
        return nil, err
    }

    if piece.Quantite < quantite {
        return nil, fmt.Errorf("stock insuffisant: disponible=%d, demandé=%d", piece.Quantite, quantite)
    }

    oldQuantite := piece.Quantite
    piece.Quantite -= quantite
    piece.UpdatedAt = time.Now()

    // Sauvegarde
    ctx := context.Background()
    pieceJSON, err := piece.ToJSON()
    if err != nil {
        return nil, fmt.Errorf("erreur de sérialisation: %w", err)
    }

    if err := s.redis.Set(ctx, PIECE_KEY_PREFIX+id, pieceJSON, 0).Err(); err != nil {
        return nil, fmt.Errorf("erreur lors de la mise à jour: %w", err)
    }

    s.logger.Info("Stock décrémenté",
        zap.String("piece_id", id),
        zap.String("nom", piece.Nom),
        zap.Int("ancien_stock", oldQuantite),
        zap.Int("nouveau_stock", piece.Quantite),
        zap.Int("decrement", quantite),
        zap.String("motif", motif))

    return piece, nil
}

// GetLowStockAlerts récupère les alertes de stock faible
func (s *StockService) GetLowStockAlerts() ([]models.AlerteStock, error) {
    pieces, err := s.GetAllPieces()
    if err != nil {
        return nil, err
    }

    alerts := make([]models.AlerteStock, 0)

    for _, piece := range pieces {
        if piece.IsLowStock() {
            severite := "attention"
            if piece.IsCriticalStock() {
                severite = "critique"
            }

            alert := models.AlerteStock{
                PieceID:          piece.ID,
                Nom:              piece.Nom,
                Quantite:         piece.Quantite,
                SeuilMin:         piece.SeuilMin,
                Severite:         severite,
                PourcentageStock: piece.GetStockPercentage(),
            }
            alerts = append(alerts, alert)
        }
    }

    return alerts, nil
}

// SearchPieces recherche des pièces par nom ou description
func (s *StockService) SearchPieces(query string) ([]models.Piece, error) {
    allPieces, err := s.GetAllPieces()
    if err != nil {
        return nil, err
    }

    searchTerm := strings.ToLower(query)
    results := make([]models.Piece, 0)

    for _, piece := range allPieces {
        if strings.Contains(strings.ToLower(piece.Nom), searchTerm) ||
           strings.Contains(strings.ToLower(piece.Description), searchTerm) ||
           strings.Contains(strings.ToLower(piece.Categorie), searchTerm) ||
           strings.Contains(strings.ToLower(piece.CodeEAN), searchTerm) {
            results = append(results, piece)
        }
    }

    return results, nil
}