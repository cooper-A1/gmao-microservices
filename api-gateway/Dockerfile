FROM nginx:alpine

# Information sur l'image
LABEL maintainer="ICS GMAO Team"
LABEL service="api-gateway"
LABEL technology="Nginx"

# Copie de la configuration Nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf


# Création des répertoires pour les logs si nécessaire
RUN mkdir -p /var/log/nginx

# Exposition du port
EXPOSE 80

# Vérification de santé
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

# Commande par défaut
CMD ["nginx", "-g", "daemon off;"]