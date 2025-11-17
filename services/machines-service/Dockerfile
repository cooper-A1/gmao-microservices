# # Dockerfile pour le service Machines (Java Spring Boot)
# FROM openjdk:17-jdk-slim

# # Information sur l'image
# LABEL maintainer="ICS GMAO Team"
# LABEL service="machines-service"
# LABEL technology="Java Spring Boot"

# # Variables d'environnement
# ENV JAVA_OPTS=""

# # Création du répertoire de l'application
# WORKDIR /app

# # Copie du fichier jar (sera créé par Maven)
# COPY target/machines-service-1.0.0.jar app.jar

# # Exposition du port
# EXPOSE 8001

# # Vérification de santé
# HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
#     CMD curl -f http://localhost:8001/api/machines/health || exit 1

# # Commande de démarrage
# ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]

# # Script de build Maven (build.sh)
# # #!/bin/bash
# # echo "Building Machines Service..."
# # mvn clean package -DskipTests
# # echo "Building Docker image..."
# # docker build -t ics-gmao/machines-service:latest .


# Étape 1 : Build
FROM maven:3.9.6-eclipse-temurin-17 AS builder

WORKDIR /app
COPY pom.xml .
COPY src ./src

RUN mvn clean package -DskipTests

# Étape 2 : Run
FROM openjdk:17-jdk-slim

WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar

EXPOSE 8001
ENTRYPOINT ["java", "-jar", "app.jar"]
