"""
Client HTTP pour communiquer avec le service Stock
"""

import httpx
import os
from typing import Optional

class StockServiceClient:
    """Client pour le service Stock (Go + Redis)"""
    
    def __init__(self):
        self.base_url = os.getenv("STOCK_SERVICE_URL", "http://stock-service:8004")
        self.timeout = 30.0
    
    async def decrement_stock(self, piece_id: str, quantite: int) -> bool:
        """Décrémente le stock d'une pièce"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.post(
                    f"{self.base_url}/api/stock/{piece_id}/decrement",
                    json={"quantite": quantite}
                )
                return response.status_code == 200
        except Exception as e:
            print(f"Erreur lors de la communication avec le service Stock: {e}")
            return False
    
    async def get_piece_info(self, piece_id: str) -> Optional[dict]:
        """Récupère les informations d'une pièce"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.get(f"{self.base_url}/api/stock/{piece_id}")
                if response.status_code == 200:
                    return response.json()
                return None
        except Exception as e:
            print(f"Erreur lors de la récupération des infos pièce: {e}")
            return None
    
    async def check_stock_availability(self, piece_id: str, quantite_requise: int) -> bool:
        """Vérifie la disponibilité du stock"""
        try:
            piece_info = await self.get_piece_info(piece_id)
            if piece_info:
                return piece_info.get("quantite", 0) >= quantite_requise
            return False
        except Exception:
            return False