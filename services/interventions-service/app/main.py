"""
Service de gestion des interventions de maintenance
Technologie: Python FastAPI + MongoDB
Port: 8002
"""

from fastapi import FastAPI, HTTPException, Depends, status
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import uvicorn
from app.database import init_db, close_db
from app.routers import interventions, auth
import os

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Initialisation de la base de données
    await init_db()
    yield
    # Nettoyage lors de l'arrêt
    await close_db()

# Configuration de l'application FastAPI
app = FastAPI(
    title="Service Interventions GMAO",
    description="API de gestion des interventions de maintenance pour ICS",
    version="1.0.0",
    docs_url="/docs",
    redoc_url="/redoc",
    lifespan=lifespan
)

# Configuration CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # En production, spécifier les domaines autorisés
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Inclusion des routers
app.include_router(interventions.router, prefix="/api/interventions", tags=["interventions"])
app.include_router(auth.router, prefix="/api/auth", tags=["authentication"])

@app.get("/")
async def root():
    """Endpoint racine du service"""
    return {
        "service": "Interventions GMAO",
        "version": "1.0.0",
        "status": "operational",
        "technology": "Python FastAPI + MongoDB"
    }

@app.get("/health")
async def health_check():
    """Endpoint de vérification de santé"""
    return {"status": "healthy", "service": "interventions-service"}

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8002,
        reload=True
    )