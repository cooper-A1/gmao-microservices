# Dockerfile pour le service Interventions (Python FastAPI)
FROM python:3.11-slim

# Information sur l'image
LABEL maintainer="ICS GMAO Team"
LABEL service="interventions-service"
LABEL technology="Python FastAPI"

# Variables d'environnement
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV MONGODB_URL=mongodb://mongo-interventions:27017
ENV JWT_SECRET=your-secret-key-here
ENV STOCK_SERVICE_URL=http://stock-service:8004
ENV TECHNICIENS_SERVICE_URL=http://techniciens-service:8003

# Création du répertoire de travail
WORKDIR /app

# Installation des dépendances système
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        curl \
        gcc \
    && rm -rf /var/lib/apt/lists/*

# Copie et installation des dépendances Python
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copie du code source
COPY ./app ./app

# Exposition du port
EXPOSE 8002

# Vérification de santé
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8002/health || exit 1

# Commande de démarrage
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8002", "--reload"]

# .env (fichier d'environnement pour développement)
# MONGODB_URL=mongodb://localhost:27017
# JWT_SECRET=your-super-secret-jwt-key-change-in-production
# STOCK_SERVICE_URL=http://localhost:8004
# TECHNICIENS_SERVICE_URL=http://localhost:8003