"""
Router pour les endpoints des interventions
"""

from fastapi import APIRouter, HTTPException, Depends, status, Query
from typing import List, Optional
from app.models import (
    InterventionCreate, InterventionUpdate, InterventionResponse, 
    InterventionCloture, StatutIntervention
)
from app.services.intervention_service import InterventionService
from app.dependencies import get_current_user

router = APIRouter()

# Injection de dépendance pour le service
def get_intervention_service() -> InterventionService:
    return InterventionService()

@router.post("/", response_model=InterventionResponse, status_code=status.HTTP_201_CREATED)
async def create_intervention(
    intervention: InterventionCreate,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Crée une nouvelle intervention"""
    return await service.create_intervention(intervention)

@router.get("/", response_model=List[InterventionResponse])
async def get_interventions(
    skip: int = Query(0, ge=0),
    limit: int = Query(100, ge=1, le=1000),
    statut: Optional[StatutIntervention] = Query(None),
    machine_id: Optional[int] = Query(None),
    technicien_id: Optional[int] = Query(None),
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Récupère la liste des interventions avec filtres optionnels"""
    return await service.get_all_interventions(skip, limit, statut, machine_id, technicien_id)

@router.get("/{intervention_id}", response_model=InterventionResponse)
async def get_intervention(
    intervention_id: str,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Récupère une intervention par son ID"""
    intervention = await service.get_intervention(intervention_id)
    if not intervention:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )
    return intervention

@router.put("/{intervention_id}", response_model=InterventionResponse)
async def update_intervention(
    intervention_id: str,
    intervention_update: InterventionUpdate,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Met à jour une intervention"""
    intervention = await service.update_intervention(intervention_id, intervention_update)
    if not intervention:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )
    return intervention

@router.delete("/{intervention_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_intervention(
    intervention_id: str,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Supprime une intervention"""
    deleted = await service.delete_intervention(intervention_id)
    if not deleted:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )

@router.post("/{intervention_id}/assign/{technicien_id}", response_model=InterventionResponse)
async def assign_technicien(
    intervention_id: str,
    technicien_id: int,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Assigne un technicien à une intervention"""
    intervention = await service.assign_technicien(intervention_id, technicien_id)
    if not intervention:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )
    return intervention

@router.post("/{intervention_id}/start", response_model=InterventionResponse)
async def start_intervention(
    intervention_id: str,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Démarre une intervention"""
    intervention = await service.start_intervention(intervention_id)
    if not intervention:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )
    return intervention

@router.post("/{intervention_id}/close", response_model=InterventionResponse)
async def close_intervention(
    intervention_id: str,
    cloture_data: InterventionCloture,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Clôture une intervention avec pièces utilisées et coût"""
    intervention = await service.close_intervention(intervention_id, cloture_data)
    if not intervention:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Intervention non trouvée"
        )
    return intervention

@router.get("/machine/{machine_id}", response_model=List[InterventionResponse])
async def get_interventions_by_machine(
    machine_id: int,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Récupère toutes les interventions pour une machine donnée"""
    return await service.get_interventions_by_machine(machine_id)

@router.get("/technicien/{technicien_id}", response_model=List[InterventionResponse])
async def get_interventions_by_technicien(
    technicien_id: int,
    service: InterventionService = Depends(get_intervention_service),
    current_user: dict = Depends(get_current_user)
):
    """Récupère toutes les interventions pour un technicien donné"""
    return await service.get_interventions_by_technicien(technicien_id)