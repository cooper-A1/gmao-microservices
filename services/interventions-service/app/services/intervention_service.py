"""
Services métier pour la gestion des interventions
"""

from typing import List, Optional, Dict, Any
from datetime import datetime
from bson import ObjectId
from fastapi import HTTPException, status
from app.database import get_database
from app.models import (
    Intervention, InterventionCreate, InterventionUpdate, 
    InterventionCloture, InterventionResponse, StatutIntervention
)
from app.services.stock_client import StockServiceClient
from app.services.techniciens_client import TechniciensServiceClient

class InterventionService:
    """Service pour gérer les interventions"""
    
    def __init__(self):
        self.stock_client = StockServiceClient()
        self.techniciens_client = TechniciensServiceClient()
    
    async def create_intervention(self, intervention_data: InterventionCreate) -> InterventionResponse:
        """Crée une nouvelle intervention"""
        db = get_database()
        
        # Vérifier la disponibilité du technicien si assigné
        if intervention_data.technicien_id:
            is_available = await self.techniciens_client.check_availability(
                intervention_data.technicien_id, 
                intervention_data.date_planifiee
            )
            if not is_available:
                raise HTTPException(
                    status_code=status.HTTP_400_BAD_REQUEST,
                    detail="Technicien non disponible à cette date"
                )
        
        # Création de l'intervention
        intervention = Intervention(**intervention_data.dict())
        
        # Insertion en base
        result = await db.interventions.insert_one(intervention.dict(by_alias=True, exclude={"id"}))
        
        # Récupération de l'intervention créée
        created_intervention = await db.interventions.find_one({"_id": result.inserted_id})
        
        return self._convert_to_response(created_intervention)
    
    async def get_intervention(self, intervention_id: str) -> Optional[InterventionResponse]:
        """Récupère une intervention par ID"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        
        if not intervention:
            return None
            
        return self._convert_to_response(intervention)
    
    async def get_all_interventions(
        self, 
        skip: int = 0, 
        limit: int = 100,
        statut: Optional[StatutIntervention] = None,
        machine_id: Optional[int] = None,
        technicien_id: Optional[int] = None
    ) -> List[InterventionResponse]:
        """Récupère toutes les interventions avec filtres optionnels"""
        db = get_database()
        
        # Construction de la requête avec filtres
        filter_query = {}
        if statut:
            filter_query["statut"] = statut
        if machine_id:
            filter_query["machine_id"] = machine_id
        if technicien_id:
            filter_query["technicien_id"] = technicien_id
        
        # Exécution de la requête
        cursor = db.interventions.find(filter_query).sort("date_planifiee", -1).skip(skip).limit(limit)
        interventions = await cursor.to_list(length=limit)
        
        return [self._convert_to_response(intervention) for intervention in interventions]
    
    async def update_intervention(self, intervention_id: str, update_data: InterventionUpdate) -> Optional[InterventionResponse]:
        """Met à jour une intervention"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        # Supprime les champs None
        update_dict = {k: v for k, v in update_data.dict().items() if v is not None}
        
        if not update_dict:
            raise HTTPException(status_code=400, detail="Aucune donnée à mettre à jour")
        
        # Vérification de la disponibilité du technicien si changé
        if "technicien_id" in update_dict and "date_planifiee" in update_dict:
            is_available = await self.techniciens_client.check_availability(
                update_dict["technicien_id"], 
                update_dict["date_planifiee"]
            )
            if not is_available:
                raise HTTPException(
                    status_code=status.HTTP_400_BAD_REQUEST,
                    detail="Technicien non disponible à cette date"
                )
        
        # Mise à jour
        result = await db.interventions.update_one(
            {"_id": ObjectId(intervention_id)},
            {"$set": update_dict}
        )
        
        if result.matched_count == 0:
            return None
        
        # Récupération de l'intervention mise à jour
        updated_intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        return self._convert_to_response(updated_intervention)
    
    async def delete_intervention(self, intervention_id: str) -> bool:
        """Supprime une intervention"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        result = await db.interventions.delete_one({"_id": ObjectId(intervention_id)})
        return result.deleted_count > 0
    
    async def assign_technicien(self, intervention_id: str, technicien_id: int) -> Optional[InterventionResponse]:
        """Assigne un technicien à une intervention"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        # Récupération de l'intervention
        intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        if not intervention:
            return None
        
        # Vérification de la disponibilité du technicien
        is_available = await self.techniciens_client.check_availability(
            technicien_id, 
            intervention["date_planifiee"]
        )
        if not is_available:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Technicien non disponible"
            )
        
        # Mise à jour de l'intervention
        await db.interventions.update_one(
            {"_id": ObjectId(intervention_id)},
            {"$set": {"technicien_id": technicien_id}}
        )
        
        # Notification au service techniciens
        await self.techniciens_client.assign_to_intervention(technicien_id, intervention_id)
        
        # Récupération de l'intervention mise à jour
        updated_intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        return self._convert_to_response(updated_intervention)
    
    async def start_intervention(self, intervention_id: str) -> Optional[InterventionResponse]:
        """Démarre une intervention"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        result = await db.interventions.update_one(
            {"_id": ObjectId(intervention_id)},
            {
                "$set": {
                    "statut": StatutIntervention.EN_COURS,
                    "date_debut": datetime.utcnow()
                }
            }
        )
        
        if result.matched_count == 0:
            return None
        
        updated_intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        return self._convert_to_response(updated_intervention)
    
    async def close_intervention(self, intervention_id: str, cloture_data: InterventionCloture) -> Optional[InterventionResponse]:
        """Clôture une intervention avec pièces utilisées et coût"""
        db = get_database()
        
        if not ObjectId.is_valid(intervention_id):
            raise HTTPException(status_code=400, detail="ID intervention invalide")
        
        # Calcul du coût total
        cout_total = sum(piece.quantite * piece.prix_unitaire for piece in cloture_data.pieces_utilisees)
        
        # Mise à jour du stock pour chaque pièce utilisée
        for piece in cloture_data.pieces_utilisees:
            success = await self.stock_client.decrement_stock(piece.piece_id, piece.quantite)
            if not success:
                raise HTTPException(
                    status_code=status.HTTP_400_BAD_REQUEST,
                    detail=f"Impossible de décrémenter le stock pour la pièce {piece.nom}"
                )
        
        # Mise à jour de l'intervention
        update_data = {
            "statut": cloture_data.statut,
            "date_fin": datetime.utcnow(),
            "duree_reelle": cloture_data.duree_reelle,
            "compte_rendu": cloture_data.compte_rendu,
            "pieces_utilisees": [piece.dict() for piece in cloture_data.pieces_utilisees],
            "cout_total": cout_total if cloture_data.cout_total is None else cloture_data.cout_total
        }
        
        result = await db.interventions.update_one(
            {"_id": ObjectId(intervention_id)},
            {"$set": update_data}
        )
        
        if result.matched_count == 0:
            return None
        
        updated_intervention = await db.interventions.find_one({"_id": ObjectId(intervention_id)})
        return self._convert_to_response(updated_intervention)
    
    async def get_interventions_by_machine(self, machine_id: int) -> List[InterventionResponse]:
        """Récupère toutes les interventions pour une machine donnée"""
        db = get_database()
        
        cursor = db.interventions.find({"machine_id": machine_id}).sort("date_planifiee", -1)
        interventions = await cursor.to_list(length=None)
        
        return [self._convert_to_response(intervention) for intervention in interventions]
    
    async def get_interventions_by_technicien(self, technicien_id: int) -> List[InterventionResponse]:
        """Récupère toutes les interventions pour un technicien donné"""
        db = get_database()
        
        cursor = db.interventions.find({"technicien_id": technicien_id}).sort("date_planifiee", -1)
        interventions = await cursor.to_list(length=None)
        
        return [self._convert_to_response(intervention) for intervention in interventions]
    
    def _convert_to_response(self, intervention_doc: dict) -> InterventionResponse:
        """Convertit un document MongoDB en InterventionResponse"""
        intervention_doc["id"] = str(intervention_doc["_id"])
        del intervention_doc["_id"]
        return InterventionResponse(**intervention_doc)
