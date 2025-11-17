"""
Client HTTP pour communiquer avec le service Techniciens
"""

import httpx
import os
from typing import Optional
from datetime import datetime

class TechniciensServiceClient:
    """Client pour le service Techniciens (Node.js + MySQL)"""
    
    def __init__(self):
        self.base_url = os.getenv("TECHNICIENS_SERVICE_URL", "http://techniciens-service:8003")
        self.timeout = 30.0
    
    async def check_availability(self, technicien_id: int, date_intervention: datetime) -> bool:
        """Vérifie la disponibilité d'un technicien"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.get(
                    f"{self.base_url}/api/techniciens/{technicien_id}/availability",
                    params={"date": date_intervention.isoformat()}
                )
                if response.status_code == 200:
                    data = response.json()
                    return data.get("available", False)
                return False
        except Exception as e:
            print(f"Erreur lors de la vérification de disponibilité: {e}")
            return True  # Par défaut, considérer comme disponible
    
    async def assign_to_intervention(self, technicien_id: int, intervention_id: str) -> bool:
        """Notifie le service techniciens de l'assignation"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.post(
                    f"{self.base_url}/api/techniciens/{technicien_id}/assign",
                    json={"intervention_id": intervention_id}
                )
                return response.status_code == 200
        except Exception as e:
            print(f"Erreur lors de l'assignation: {e}")
            return False
    
    async def get_technicien_info(self, technicien_id: int) -> Optional[dict]:
        """Récupère les informations d'un technicien"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.get(f"{self.base_url}/api/techniciens/{technicien_id}")
                if response.status_code == 200:
                    return response.json()
                return None
        except Exception as e:
            print(f"Erreur lors de la récupération des infos technicien: {e}")
            return None
