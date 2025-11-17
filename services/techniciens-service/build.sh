# build.sh - Script de build
#!/bin/bash

echo "🔨 Building Techniciens Service..."

# Installation des dépendances
echo "📦 Installing dependencies..."
npm install

# Tests (si disponibles)
# echo "🧪 Running tests..."
# npm test

# Build de l'image Docker
echo "🐳 Building Docker image..."
docker build -t ics-gmao/techniciens-service:latest .

echo "✅ Build completed successfully!"
echo "🚀 To run: docker run -p 8003:8003 ics-gmao/techniciens-service:latest"