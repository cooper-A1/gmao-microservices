/**
 * Modèle Technicien avec validation
 */

const Joi = require('joi');

/**
 * Schéma de validation pour la création d'un technicien
 */
const createTechnicienSchema = Joi.object({
    nom: Joi.string().min(2).max(100).required()
        .messages({
            'string.empty': 'Le nom est obligatoire',
            'string.min': 'Le nom doit contenir au moins 2 caractères',
            'string.max': 'Le nom ne peut pas dépasser 100 caractères'
        }),
    
    prenom: Joi.string().min(2).max(100).required()
        .messages({
            'string.empty': 'Le prénom est obligatoire',
            'string.min': 'Le prénom doit contenir au moins 2 caractères'
        }),
    
    email: Joi.string().email().max(255).required()
        .messages({
            'string.email': 'L\'email doit être valide',
            'string.empty': 'L\'email est obligatoire'
        }),
    
    telephone: Joi.string().pattern(/^[+]?[0-9\s\-\(\)]{8,20}$/).required()
        .messages({
            'string.pattern.base': 'Le numéro de téléphone n\'est pas valide'
        }),
    
    competences: Joi.array().items(Joi.string()).min(1).required()
        .messages({
            'array.min': 'Au moins une compétence est requise'
        }),
    
    niveau_experience: Joi.string().valid('junior', 'senior', 'expert').default('junior'),
    
    disponibilite: Joi.boolean().default(true),
    
    salaire: Joi.number().positive().optional(),
    
    date_embauche: Joi.date().optional(),
    
    notes: Joi.string().max(1000).optional()
});

/**
 * Schéma de validation pour la mise à jour d'un technicien
 */
const updateTechnicienSchema = Joi.object({
    nom: Joi.string().min(2).max(100).optional(),
    prenom: Joi.string().min(2).max(100).optional(),
    email: Joi.string().email().max(255).optional(),
    telephone: Joi.string().pattern(/^[+]?[0-9\s\-\(\)]{8,20}$/).optional(),
    competences: Joi.array().items(Joi.string()).min(1).optional(),
    niveau_experience: Joi.string().valid('junior', 'senior', 'expert').optional(),
    disponibilite: Joi.boolean().optional(),
    salaire: Joi.number().positive().optional(),
    date_embauche: Joi.date().optional(),
    notes: Joi.string().max(1000).optional()
});

/**
 * Énumérations pour les niveaux d'expérience
 */
const NIVEAUX_EXPERIENCE = {
    JUNIOR: 'junior',
    SENIOR: 'senior',  
    EXPERT: 'expert'
};

/**
 * Compétences techniques disponibles
 */
const COMPETENCES_DISPONIBLES = [
    'Mécanique générale',
    'Électricité industrielle',
    'Pneumatique',
    'Hydraulique', 
    'Automatisme',
    'Électronique',
    'Soudure',
    'Usinage',
    'Maintenance préventive',
    'Diagnostic de pannes',
    'Informatique industrielle',
    'Régulation',
    'Climatisation',
    'Plomberie industrielle'
];

module.exports = {
    createTechnicienSchema,
    updateTechnicienSchema,
    NIVEAUX_EXPERIENCE,
    COMPETENCES_DISPONIBLES
};
        