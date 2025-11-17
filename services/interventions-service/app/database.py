"""
Configuration de la base de données MongoDB
"""

from motor.motor_asyncio import AsyncIOMotorClient
from pymongo import ASCENDING, DESCENDING
import os

# Configuration MongoDB
MONGODB_URL = os.getenv("MONGODB_URL", "mongodb://mongo-interventions:27017")
DATABASE_NAME = "interventions_db"

# Client MongoDB global
client: AsyncIOMotorClient = None
database = None

async def init_db():
    """Initialise la connexion à MongoDB"""
    global client, database
    
    client = AsyncIOMotorClient(MONGODB_URL)
    database = client[DATABASE_NAME]
    
    # Création des index pour optimiser les requêtes
    await create_indexes()
    
    print(f"✅ Connexion MongoDB établie: {MONGODB_URL}")

async def close_db():
    """Ferme la connexion MongoDB"""
    global client
    if client:
        client.close()
        print("✅ Connexion MongoDB fermée")

async def create_indexes():
    """Crée les index pour optimiser les performances"""
    interventions_collection = database.interventions
    
    # Index sur machine_id pour les recherches par machine
    await interventions_collection.create_index([("machine_id", ASCENDING)])
    
    # Index sur technicien_id pour les recherches par technicien
    await interventions_collection.create_index([("technicien_id", ASCENDING)])
    
    # Index sur statut pour les recherches par état
    await interventions_collection.create_index([("statut", ASCENDING)])
    
    # Index sur date_planifiee pour les tris chronologiques
    await interventions_collection.create_index([("date_planifiee", DESCENDING)])
    
    # Index composé pour les recherches complexes
    await interventions_collection.create_index([
        ("machine_id", ASCENDING),
        ("statut", ASCENDING),
        ("date_planifiee", DESCENDING)
    ])

def get_database():
    """Retourne l'instance de la base de données"""
    return database