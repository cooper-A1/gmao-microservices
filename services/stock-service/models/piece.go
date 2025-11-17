package models

import (
    "encoding/json"
    "time"
)

// Piece représente une pièce détachée en stock
type Piece struct {
    ID           string    `json:"id" redis:"id"`
    Nom          string    `json:"nom" redis:"nom" binding:"required"`
    Description  string    `json:"description" redis:"description"`
    Quantite     int       `json:"quantite" redis:"quantite" binding:"required,min=0"`
    SeuilMin     int       `json:"seuil_min" redis:"seuil_min" binding:"required,min=1"`
    PrixUnitaire float64   `json:"prix_unitaire" redis:"prix_unitaire" binding:"required,gt=0"`
    Fournisseur  string    `json:"fournisseur" redis:"fournisseur"`
    Emplacement  string    `json:"emplacement" redis:"emplacement"`
    CodeEAN      string    `json:"code_ean" redis:"code_ean"`
    Categorie    string    `json:"categorie" redis:"categorie"`
    UniteStock   string    `json:"unite_stock" redis:"unite_stock" binding:"required"`
    CreatedAt    time.Time `json:"created_at" redis:"created_at"`
    UpdatedAt    time.Time `json:"updated_at" redis:"updated_at"`
}

// CreatePieceRequest représente une requête de création de pièce
type CreatePieceRequest struct {
    Nom          string  `json:"nom" binding:"required,min=3,max=200"`
    Description  string  `json:"description" binding:"max=1000"`
    Quantite     int     `json:"quantite" binding:"required,min=0"`
    SeuilMin     int     `json:"seuil_min" binding:"required,min=1"`
    PrixUnitaire float64 `json:"prix_unitaire" binding:"required,gt=0"`
    Fournisseur  string  `json:"fournisseur" binding:"max=200"`
    Emplacement  string  `json:"emplacement" binding:"max=50"`
    CodeEAN      string  `json:"code_ean" binding:"max=50"`
    Categorie    string  `json:"categorie" binding:"required,max=100"`
    UniteStock   string  `json:"unite_stock" binding:"required,max=20"`
}

// UpdatePieceRequest représente une requête de mise à jour de pièce
type UpdatePieceRequest struct {
    Nom          *string  `json:"nom,omitempty" binding:"omitempty,min=3,max=200"`
    Description  *string  `json:"description,omitempty" binding:"omitempty,max=1000"`
    SeuilMin     *int     `json:"seuil_min,omitempty" binding:"omitempty,min=1"`
    PrixUnitaire *float64 `json:"prix_unitaire,omitempty" binding:"omitempty,gt=0"`
    Fournisseur  *string  `json:"fournisseur,omitempty" binding:"omitempty,max=200"`
    Emplacement  *string  `json:"emplacement,omitempty" binding:"omitempty,max=50"`
    CodeEAN      *string  `json:"code_ean,omitempty" binding:"omitempty,max=50"`
    Categorie    *string  `json:"categorie,omitempty" binding:"omitempty,max=100"`
    UniteStock   *string  `json:"unite_stock,omitempty" binding:"omitempty,max=20"`
}

// StockMovementRequest représente une requête de mouvement de stock
type StockMovementRequest struct {
    Quantite int    `json:"quantite" binding:"required,gt=0"`
    Motif    string `json:"motif,omitempty" binding:"max=500"`
}

// AlerteStock représente une alerte de stock faible
type AlerteStock struct {
    PieceID      string  `json:"piece_id"`
    Nom          string  `json:"nom"`
    Quantite     int     `json:"quantite"`
    SeuilMin     int     `json:"seuil_min"`
    Severite     string  `json:"severite"` // "critique", "attention"
    PourcentageStock float64 `json:"pourcentage_stock"`
}

// ToJSON convertit la pièce en JSON
func (p *Piece) ToJSON() ([]byte, error) {
    return json.Marshal(p)
}

// FromJSON crée une pièce depuis du JSON
func (p *Piece) FromJSON(data []byte) error {
    return json.Unmarshal(data, p)
}

// IsLowStock vérifie si la pièce est en stock faible
func (p *Piece) IsLowStock() bool {
    return p.Quantite <= p.SeuilMin
}

// IsCriticalStock vérifie si la pièce est en stock critique
func (p *Piece) IsCriticalStock() bool {
    return p.Quantite <= (p.SeuilMin / 2)
}

// GetStockPercentage calcule le pourcentage de stock par rapport au seuil
func (p *Piece) GetStockPercentage() float64 {
    if p.SeuilMin == 0 {
        return 100.0
    }
    return (float64(p.Quantite) / float64(p.SeuilMin)) * 100.0
}

