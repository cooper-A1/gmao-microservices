/**
 * Middleware de gestion globale des erreurs
 */

const logger = require('../utils/logger');

const errorHandler = (error, req, res, next) => {
    logger.error('Erreur non gérée:', {
        message: error.message,
        stack: error.stack,
        url: req.url,
        method: req.method,
        ip: req.ip
    });

    // Erreurs de validation Joi
    if (error.name === 'ValidationError') {
        return res.status(400).json({
            error: 'Données invalides',
            details: error.details?.map(d => d.message) || [error.message]
        });
    }

    // Erreurs MySQL
    if (error.code === 'ER_DUP_ENTRY') {
        return res.status(409).json({
            error: 'Données déjà existantes',
            message: 'Cette ressource existe déjà'
        });
    }

    if (error.code === 'ER_NO_REFERENCED_ROW_2') {
        return res.status(400).json({
            error: 'Référence invalide',
            message: 'Une des références spécifiées n\'existe pas'
        });
    }

    // Erreurs personnalisées métier
    if (error.message === 'Un technicien avec cet email existe déjà') {
        return res.status(409).json({
            error: error.message
        });
    }

    if (error.message === 'Technicien non trouvé') {
        return res.status(404).json({
            error: error.message
        });
    }

    // Erreur générique
    res.status(500).json({
        error: 'Erreur interne du serveur',
        message: process.env.NODE_ENV === 'development' ? error.message : undefined
    });
};

module.exports = {
    errorHandler
};