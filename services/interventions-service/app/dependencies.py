
"""
Dépendances FastAPI pour l'authentification
"""

from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from typing import Dict
from app.auth import verify_token

security = HTTPBearer()

async def get_current_user(credentials: HTTPAuthorizationCredentials = Depends(security)) -> Dict:
    """Dépendance pour récupérer l'utilisateur actuel depuis le JWT"""
    
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Token invalide",
        headers={"WWW-Authenticate": "Bearer"},
    )
    
    try:
        payload = verify_token(credentials.credentials)
        if payload is None:
            raise credentials_exception
            
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
            
        return {
            "username": username,
            "user_id": payload.get("user_id"),
            "role": payload.get("role", "technicien")
        }
        
    except Exception:
        raise credentials_exception

def require_role(required_role: str):
    """Dépendance pour vérifier le rôle de l'utilisateur"""
    def role_checker(current_user: Dict = Depends(get_current_user)):
        user_role = current_user.get("role", "")
        
        if user_role != required_role and user_role != "admin":
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Permissions insuffisantes"
            )
        return current_user
    return role_checker

# Alias pour les rôles courants
require_admin = require_role("admin")
require_manager = require_role("manager")
