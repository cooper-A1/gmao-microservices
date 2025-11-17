# stop.sh - Script d'arrêt
#!/bin/bash

echo "🛑 Arrêt du système GMAO..."

docker-compose down

echo "✅ Système arrêté!"