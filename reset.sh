
# reset.sh - Script de remise à zéro
#!/bin/bash

echo "⚠️  ATTENTION: Ceci va supprimer toutes les données!"
read -p "Êtes-vous sûr? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🗑️  Suppression des conteneurs et volumes..."
    docker-compose down -v
    docker system prune -f
    echo "✅ Système remis à zéro!"
else
    echo "❌ Opération annulée"
fi
