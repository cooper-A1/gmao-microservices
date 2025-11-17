/**
 * Middlewares de validation
 */

/**
 * Valide que l'ID est un entier positif
 */
const validateId = (req, res, next) => {
    const { id } = req.params;
    
    if (!id || isNaN(id) || parseInt(id) <= 0) {
        return res.status(400).json({
            error: 'ID invalide',
            message: 'L\'ID doit être un entier positif'
        });
    }
    
    req.params.id = parseInt(id);
    next();
};

module.exports = {
    validateId
};
