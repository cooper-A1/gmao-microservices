"""
Modèles Pydantic pour les interventions
"""

from pydantic import BaseModel, Field, validator
from typing import Optional, List, Dict, Any
from datetime import datetime
from enum import Enum
from bson import ObjectId

class PyObjectId(ObjectId):
    """Classe pour gérer les ObjectId MongoDB avec Pydantic"""
    @classmethod
    def __get_validators__(cls):
        yield cls.validate

    @classmethod
    def validate(cls, v):
        if not ObjectId.is_valid(v):
            raise ValueError("Invalid ObjectId")
        return ObjectId(v)

    @classmethod
    def __modify_schema__(cls, field_schema):
        field_schema.update(type="string")

class TypeIntervention(str, Enum):
    """Types d'interventions possibles"""
    PREVENTIVE = "preventive"
    CORRECTIVE = "corrective"
    PREDICTIVE = "predictive"
    AMELIORATIVE = "ameliorative"

class StatutIntervention(str, Enum):
    """Statuts possibles pour une intervention"""
    PLANIFIEE = "planifiee"
    EN_COURS = "en_cours"
    TERMINEE = "terminee"
    ANNULEE = "annulee"
    REPORTEE = "reportee"

class PieceUtilisee(BaseModel):
    """Modèle pour une pièce utilisée lors d'une intervention"""
    piece_id: str = Field(..., description="ID de la pièce dans le service stock")
    nom: str = Field(..., description="Nom de la pièce")
    quantite: int = Field(..., gt=0, description="Quantité utilisée")
    prix_unitaire: float = Field(..., ge=0, description="Prix unitaire")
    
    @validator('quantite')
    def validate_quantite(cls, v):
        if v <= 0:
            raise ValueError('La quantité doit être positive')
        return v

class InterventionBase(BaseModel):
    """Modèle de base pour une intervention"""
    machine_id: int = Field(..., description="ID de la machine concernée")
    type_intervention: TypeIntervention = Field(..., description="Type d'intervention")
    titre: str = Field(..., min_length=5, max_length=200, description="Titre de l'intervention")
    description: Optional[str] = Field(None, max_length=1000, description="Description détaillée")
    technicien_id: Optional[int] = Field(None, description="ID du technicien assigné")
    date_planifiee: datetime = Field(..., description="Date et heure planifiées")
    duree_estimee: Optional[int] = Field(None, gt=0, description="Durée estimée en minutes")
    priorite: int = Field(1, ge=1, le=5, description="Priorité (1=faible, 5=critique)")

class InterventionCreate(InterventionBase):
    """Modèle pour créer une intervention"""
    pass

class InterventionUpdate(BaseModel):
    """Modèle pour mettre à jour une intervention"""
    type_intervention: Optional[TypeIntervention] = None
    titre: Optional[str] = Field(None, min_length=5, max_length=200)
    description: Optional[str] = Field(None, max_length=1000)
    technicien_id: Optional[int] = None
    date_planifiee: Optional[datetime] = None
    duree_estimee: Optional[int] = Field(None, gt=0)
    priorite: Optional[int] = Field(None, ge=1, le=5)
    statut: Optional[StatutIntervention] = None

class InterventionCloture(BaseModel):
    """Modèle pour clôturer une intervention"""
    statut: StatutIntervention = Field(StatutIntervention.TERMINEE)
    compte_rendu: str = Field(..., min_length=10, description="Compte rendu de l'intervention")
    pieces_utilisees: List[PieceUtilisee] = Field(default=[], description="Liste des pièces utilisées")
    duree_reelle: int = Field(..., gt=0, description="Durée réelle en minutes")
    cout_total: Optional[float] = Field(None, ge=0, description="Coût total (calculé automatiquement)")

class Intervention(InterventionBase):
    """Modèle complet d'une intervention"""
    id: Optional[PyObjectId] = Field(default_factory=PyObjectId, alias="_id")
    statut: StatutIntervention = Field(default=StatutIntervention.PLANIFIEE)
    date_creation: datetime = Field(default_factory=datetime.utcnow)
    date_debut: Optional[datetime] = None
    date_fin: Optional[datetime] = None
    duree_reelle: Optional[int] = None
    compte_rendu: Optional[str] = None
    pieces_utilisees: List[PieceUtilisee] = Field(default=[])
    cout_total: float = Field(default=0.0)
    
    class Config:
        allow_population_by_field_name = True
        arbitrary_types_allowed = True
        json_encoders = {ObjectId: str}

class InterventionResponse(BaseModel):
    """Modèle de réponse pour une intervention"""
    id: str
    machine_id: int
    type_intervention: TypeIntervention
    titre: str
    description: Optional[str]
    technicien_id: Optional[int]
    statut: StatutIntervention
    date_planifiee: datetime
    date_creation: datetime
    date_debut: Optional[datetime]
    date_fin: Optional[datetime]
    duree_estimee: Optional[int]
    duree_reelle: Optional[int]
    priorite: int
    compte_rendu: Optional[str]
    pieces_utilisees: List[PieceUtilisee]
    cout_total: float