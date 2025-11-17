/**
 * Service Techniciens - GMAO ICS
 * Technologie: Node.js Express + MySQL
 * Port: 8003
 */

require('dotenv').config();
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const rateLimit = require('express-rate-limit');
const swaggerUi = require('swagger-ui-express');
const swaggerSpec = require('./config/swagger');
const logger = require('./utils/logger');
const database = require('./database/connection');

// Import des routes
const techniciensRoutes = require('./routes/techniciens');
const authRoutes = require('./routes/auth');
const { errorHandler } = require('./middleware/errorHandler');
const { authMiddleware } = require('./middleware/auth');

const app = express();
const PORT = process.env.PORT || 8003;

// Configuration des middlewares de sécurité
app.use(helmet());
app.use(cors({
    origin: process.env.ALLOWED_ORIGINS?.split(',') || '*',
    credentials: true
}));

// Limitation du taux de requêtes
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // Limite à 100 requêtes par fenêtre de 15 min
    message: 'Trop de requêtes, réessayez plus tard.',
    standardHeaders: true,
    legacyHeaders: false,
});
app.use('/api/', limiter);

// Middlewares pour parsing JSON
app.use(express.json({ limit: '10mb' }));
app.use(express.urlencoded({ extended: true }));

// Logging des requêtes
app.use((req, res, next) => {
    logger.info(`${req.method} ${req.url}`, {
        ip: req.ip,
        userAgent: req.get('User-Agent')
    });
    next();
});

// Routes de l'API
app.use('/api/auth', authRoutes);
app.use('/api/techniciens', techniciensRoutes);

// Documentation Swagger
app.use('/docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));

// Route racine
app.get('/', (req, res) => {
    res.json({
        service: 'Service Techniciens GMAO',
        version: '1.0.0',
        status: 'operational',
        technology: 'Node.js Express + MySQL',
        documentation: '/docs',
        health: '/health'
    });
});

// Endpoint de santé
app.get('/health', async (req, res) => {
    try {
        // Vérification de la connexion à la base de données
        await database.query('SELECT 1');
        
        res.status(200).json({
            status: 'healthy',
            service: 'techniciens-service',
            timestamp: new Date().toISOString(),
            database: 'connected'
        });
    } catch (error) {
        logger.error('Health check failed:', error);
        res.status(503).json({
            status: 'unhealthy',
            service: 'techniciens-service',
            timestamp: new Date().toISOString(),
            database: 'disconnected',
            error: error.message
        });
    }
});

// Middleware de gestion des erreurs
app.use(errorHandler);

// Gestion des routes non trouvées
app.use('*', (req, res) => {
    res.status(404).json({
        error: 'Route non trouvée',
        method: req.method,
        url: req.originalUrl
    });
});

// Initialisation de la base de données et démarrage du serveur
async function startServer() {
    try {
        // Test de connexion à la base de données
        await database.testConnection();
        logger.info('Connexion MySQL établie');
        
        // Démarrage du serveur
        app.listen(PORT, '0.0.0.0', () => {
            logger.info(`🚀 Service Techniciens démarré sur le port ${PORT}`);
            logger.info(`📚 Documentation disponible sur http://localhost:${PORT}/docs`);
        });
        
    } catch (error) {
        logger.error('Erreur lors du démarrage:', error);
        process.exit(1);
    }
}

// Gestion propre de l'arrêt
process.on('SIGTERM', async () => {
    logger.info('SIGTERM reçu, arrêt en cours...');
    await database.close();
    process.exit(0);
});

process.on('SIGINT', async () => {
    logger.info('SIGINT reçu, arrêt en cours...');
    await database.close();
    process.exit(0);
});

startServer();