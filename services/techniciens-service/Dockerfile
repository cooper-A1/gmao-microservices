
#Dockerfile pour le service Techniciens
FROM node:18-alpine

# Information sur l'image
LABEL maintainer="ICS GMAO Team"
LABEL service="techniciens-service"
LABEL technology="Node.js Express"

# Variables d'environnement
ENV NODE_ENV=production
ENV PORT=8003
ENV DB_HOST=mysql-techniciens
ENV DB_PORT=3306
ENV DB_USER=gmao_user
ENV DB_PASSWORD=gmao_password
ENV DB_NAME=techniciens_db
ENV JWT_SECRET=your-secret-key-here

# Installation des packages système nécessaires
RUN apk add --no-cache curl

# Création du répertoire de travail
WORKDIR /app

# Création du répertoire logs
RUN mkdir -p logs

# Copie des fichiers de configuration
COPY package*.json ./

# Installation des dépendances
RUN npm ci --only=production && npm cache clean --force

# Copie du code source
COPY ./src ./src

# Création d'un utilisateur non-root
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Attribution des permissions
RUN chown -R nodejs:nodejs /app
USER nodejs

# Exposition du port
EXPOSE 8003

# Vérification de santé
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD curl -f http://localhost:8003/health || exit 1

# Commande de démarrage
CMD ["node", "src/server.js"]
