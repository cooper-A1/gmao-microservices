#!/bin/bash
echo "🚀 Démarrage du système GMAO ICS..."

# Vérification de Docker et Docker Compose
if ! command -v docker &> /dev/null; then
    echo "❌ Docker n'est pas installé"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose n'est pas installé"
    exit 1
fi

# Création du fichier .env s'il n'existe pas
if [ ! -f .env ]; then
    echo "📝 Création du fichier .env..."
    cp .env.example .env
    echo "⚠️  IMPORTANT: Modifier le fichier .env avec vos paramètres"
fi

# Build et démarrage des services
echo "🔨 Build des services..."
docker-compose build

echo "🐳 Démarrage des conteneurs..."
docker-compose up -d

# Attente du démarrage
echo "⏳ Attente du démarrage des services..."
sleep 30

# Vérification des services
echo "🔍 Vérification des services..."
services=(
    "http://localhost/health:API Gateway"
    "http://localhost:8001/api/machines/health:Service Machines"
    "http://localhost:8002/health:Service Interventions"
    "http://localhost:8003/health:Service Techniciens"
    "http://localhost:8004/health:Service Stock"
    "http://localhost:8005/health:Service Prédiction"
)

for service in "${services[@]}"; do
    url="${service%%:*}"
    name="${service##*:}"
    if curl -f -s "$url" > /dev/null; then
        echo "✅ $name: OK"
    else
        echo "❌ $name: Échec"
    fi
done

echo ""
echo "🎉 Système GMAO démarré!"
echo "🌐 API Gateway: http://localhost"
echo "📚 Documentation:"
echo "   - Machines: http://localhost/machines-docs"
echo "   - Interventions: http://localhost/interventions-docs"
echo "   - Techniciens: http://localhost/techniciens-docs"
echo "   - Stock: http://localhost/stock-docs"
echo "🖥️  Portainer: http://localhost:9000"
echo ""
echo "👤 Comptes de test:"
echo "   - Admin: admin / admin123"
echo "   - Manager: manager / manager123"
echo "   - Technicien: tech1 / tech123"
