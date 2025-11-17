/**
 * Middleware d'authentification JWT
 */

const jwt = require('jsonwebtoken');
const logger = require('../utils/logger');

/**
 * Middleware de vérification du token JWT
 */
const authMiddleware = (req, res, next) => {
    try {
        const authHeader = req.headers.authorization;
        
        if (!authHeader || !authHeader.startsWith('Bearer ')) {
            return res.status(401).json({
                error: 'Token d\'authentification requis'
            });
        }

        const token = authHeader.substring(7); // Supprime "Bearer "
        
        const decoded = jwt.verify(token, process.env.JWT_SECRET || 'your-secret-key');
        
        // Ajout des informations utilisateur à la requête
        req.user = {
            username: decoded.sub,
            user_id: decoded.user_id,
            role: decoded.role
        };
        
        next();
        
    } catch (error) {
        logger.error('Erreur d\'authentification:', error);
        
        if (error.name === 'TokenExpiredError') {
            return res.status(401).json({
                error: 'Token expiré'
            });
        }
        
        return res.status(401).json({
            error: 'Token invalide'
        });
    }
};

/**
 * Middleware de vérification des rôles
 */
const requireRole = (allowedRoles) => {
    return (req, res, next) => {
        if (!req.user) {
            return res.status(401).json({
                error: 'Authentification requise'
            });
        }

        const userRole = req.user.role;
        
        // L'admin a tous les droits
        if (userRole === 'admin') {
            return next();
        }
        
        // Vérification du rôle spécifique
        if (!allowedRoles.includes(userRole)) {
            return res.status(403).json({
                error: 'Permissions insuffisantes',
                required_roles: allowedRoles,
                user_role: userRole
            });
        }
        
        next();
    };
};

module.exports = {
    authMiddleware,
    requireRole
};
