
# README.md - Documentation complète
# Système GMAO - Architecture Microservices

## 🏗️ Architecture

Ce système GMAO (Gestion de Maintenance Assistée par Ordinateur) est conçu selon une architecture microservices polyglotte pour l'entreprise industrielle ICS au Sénégal.

### Services Microservices

| Service | Technologie | Base de données | Port | Description |
|---------|-------------|-----------------|------|-------------|
| **Machines** | Java Spring Boot | PostgreSQL | 8001 | Gestion du parc machines |
| **Interventions** | Python FastAPI | MongoDB | 8002 | Gestion des interventions |
| **Techniciens** | Node.js Express | MySQL | 8003 | Gestion des techniciens |
| **Stock** | Go Gin | Redis | 8004 | Gestion des pièces détachées |
| **Prédiction IA** | Python | - | 8005 | Analyse prédictive (bonus) |
| **API Gateway** | Nginx | - | 80 | Point d'entrée unique |

## 🚀 Démarrage rapide

### Prérequis
- Docker et Docker Compose
- Au moins 4GB de RAM libre
- Ports 80, 8001-8005, 3306, 5432, 6379, 27017, 9000 disponibles

### Installation

```bash
# Cloner le projet
git clone <repo-url>
cd gmao-microservices

# Démarrer tous les services
chmod +x start.sh
./start.sh
```

### Accès aux services

- **API Gateway**: http://localhost
- **Documentation Swagger**: 
  - Machines: http://localhost/machines-docs
  - Interventions: http://localhost/interventions-docs  
  - Techniciens: http://localhost/techniciens-docs
  - Stock: http://localhost/stock-docs
- **Portainer**: http://localhost:9000

## 🔐 Authentification

Le système utilise JWT. Comptes de test:

```json
{
  "admin": { "username": "admin", "password": "admin123" },
  "manager": { "username": "manager", "password": "manager123" },
  "technicien": { "username": "tech1", "password": "tech123" }
}
```

### Obtenir un token

```bash
curl -X POST http://localhost/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

## 📖 Exemples d'API

### 1. Créer une machine

```bash
curl -X POST http://localhost/api/machines \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "nom": "Presse hydraulique",
    "site": "Atelier",
    "dateInstallation": "2024-01-01T10:00:00",
    "etat": "OPERATIONNELLE"
  }'
```

### 2. Créer une intervention

```bash
curl -X POST http://localhost/api/interventions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "machine_id": 1,
    "type_intervention": "preventive",
    "titre": "Maintenance mensuelle",
    "date_planifiee": "2024-02-01T09:00:00",
    "priorite": 3
  }'
```

### 3. Gérer le stock

```bash
# Décrémenter une pièce
curl -X POST http://localhost/api/stock/piece-001/decrement \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "quantite": 2,
    "motif": "Intervention maintenance"
  }'
```

### 4. Prédiction IA

```bash
curl http://localhost/api/prediction/123 \
  -H "Authorization: Bearer <token>"
```

## 🛠️ Développement

### Structure du projet

```
gmao-microservices/
├── api-gateway/           # Nginx reverse proxy
├── services/
│   ├── machines-service/     # Java Spring Boot
│   ├── interventions-service/ # Python FastAPI
│   ├── techniciens-service/   # Node.js Express  
│   ├── stock-service/         # Go Gin
│   └── prediction-service/    # Python IA
├── database/
│   └── init-scripts/      # Scripts d'initialisation
├── docker-compose.yml     # Orchestration
└── README.md
```

### Commandes utiles

```bash
# Voir les logs
./logs.sh [service-name]

# Arrêter le système  
./stop.sh

# Reset complet (⚠️ supprime les données)
./reset.sh

# Status des services
docker-compose ps
```

## 🏭 Spécifique à ICS Sénégal

- Données de test avec du matériel industriel sénégalais
- Compétences techniques adaptées au contexte local
- Fournisseurs basés au Sénégal (SKF, Total, Schneider Electric)
- Support multilingue (français/anglais)

## 📊 Monitoring et observabilité

- **Logs centralisés**: Chaque service log dans sa technologie
- **Health checks**: Endpoints `/health` pour chaque service
- **Métriques**: Actuator (Spring), built-in pour autres
- **Portainer**: Interface graphique de gestion des conteneurs

## 🔒 Sécurité

- Authentification JWT centralisée
- Autorisation par rôles (admin, manager, technicien)
- CORS configuré
- Variables d'environnement pour les secrets
- Base de données avec authentification

## 🚢 Déploiement Production

### Docker Swarm ou Kubernetes

Le système est conçu pour être facilement déployé en production avec:

- Health checks pour rolling updates
- Variables d'environnement externalisées
- Volumes persistants pour les données
- Services stateless (sauf bases de données)

### Configuration production

1. Modifier les secrets dans `.env`
2. Configurer HTTPS sur l'API Gateway
3. Utiliser des bases de données externes managées
4. Configurer la sauvegarde des volumes
5. Mettre en place le monitoring (Prometheus/Grafana)

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature
3. Commiter les changements
4. Pousser vers la branche
5. Ouvrir une Pull Request

## 📄 Licence

MIT License - Voir le fichier LICENSE