
"""
Router pour l'authentification
"""

from fastapi import APIRouter, HTTPException, status, Depends
from fastapi.security import OAuth2PasswordRequestForm
from pydantic import BaseModel
from typing import Dict
from app.dependencies import get_current_user
from app.auth import verify_password, create_access_token, get_password_hash
from app.database import get_database

router = APIRouter()

class TokenResponse(BaseModel):
    """Modèle de réponse pour un token"""
    access_token: str
    token_type: str
    user_info: Dict

class UserCreate(BaseModel):
    """Modèle pour créer un utilisateur"""
    username: str
    email: str
    password: str
    role: str = "technicien"

# Utilisateurs par défaut pour la démo (en production, utiliser une vraie base)
DEMO_USERS = {
    "admin": {
        "username": "admin",
        "email": "admin@ics.sn",
        "hashed_password": get_password_hash("admin123"),
        "role": "admin",
        "user_id": 1
    },
    "manager": {
        "username": "manager",
        "email": "manager@ics.sn", 
        "hashed_password": get_password_hash("manager123"),
        "role": "manager",
        "user_id": 2
    },
    "tech1": {
        "username": "tech1",
        "email": "tech1@ics.sn",
        "hashed_password": get_password_hash("tech123"),
        "role": "technicien",
        "user_id": 3
    }
}

@router.post("/login", response_model=TokenResponse)
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    """Endpoint de connexion"""
    
    # Recherche de l'utilisateur
    user = DEMO_USERS.get(form_data.username)
    
    if not user or not verify_password(form_data.password, user["hashed_password"]):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Nom d'utilisateur ou mot de passe incorrect",
            headers={"WWW-Authenticate": "Bearer"},
        )
    
    # Création du token
    access_token = create_access_token(
        data={
            "sub": user["username"],
            "user_id": user["user_id"],
            "role": user["role"]
        }
    )
    
    return TokenResponse(
        access_token=access_token,
        token_type="bearer",
        user_info={
            "username": user["username"],
            "email": user["email"],
            "role": user["role"],
            "user_id": user["user_id"]
        }
    )

@router.get("/me")
async def get_current_user_info(current_user: Dict = Depends(get_current_user)):
    """Récupère les informations de l'utilisateur connecté"""
    return {
        "username": current_user["username"],
        "role": current_user["role"],
        "user_id": current_user["user_id"]
    }