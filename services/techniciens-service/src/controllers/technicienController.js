
/**
 * Contrôleur pour la gestion des techniciens
 */

const technicienService = require('../services/technicienService');
const { createTechnicienSchema, updateTechnicienSchema } = require('../models/Technicien');
const logger = require('../utils/logger');

class TechnicienController {

    /**
     * Crée un nouveau technicien
     */
    async createTechnicien(req, res, next) {
        try {
            // Validation des données
            const { error, value } = createTechnicienSchema.validate(req.body);
            if (error) {
                return res.status(400).json({
                    error: 'Données invalides',
                    details: error.details.map(d => d.message)
                });
            }

            const technicien = await technicienService.createTechnicien(value);
            
            logger.info(`Technicien créé: ${technicien.nom} ${technicien.prenom}`);
            
            res.status(201).json({
                message: 'Technicien créé avec succès',
                data: technicien
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Récupère tous les techniciens
     */
    async getAllTechniciens(req, res, next) {
        try {
            const filters = {
                disponibilite: req.query.disponibilite === 'true' ? true : 
                             req.query.disponibilite === 'false' ? false : undefined,
                niveau_experience: req.query.niveau_experience,
                competence: req.query.competence,
                search: req.query.search,
                limit: req.query.limit,
                offset: req.query.offset
            };

            const techniciens = await technicienService.getAllTechniciens(filters);
            
            res.json({
                message: 'Liste des techniciens récupérée',
                data: techniciens,
                count: techniciens.length
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Récupère un technicien par ID
     */
    async getTechnicienById(req, res, next) {
        try {
            const { id } = req.params;
            
            const technicien = await technicienService.findById(id);
            if (!technicien) {
                return res.status(404).json({
                    error: 'Technicien non trouvé'
                });
            }
            
            res.json({
                message: 'Technicien trouvé',
                data: technicien
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Met à jour un technicien
     */
    async updateTechnicien(req, res, next) {
        try {
            const { id } = req.params;
            
            // Validation des données
            const { error, value } = updateTechnicienSchema.validate(req.body);
            if (error) {
                return res.status(400).json({
                    error: 'Données invalides',
                    details: error.details.map(d => d.message)
                });
            }

            const technicien = await technicienService.updateTechnicien(id, value);
            
            logger.info(`Technicien mis à jour: ${id}`);
            
            res.json({
                message: 'Technicien mis à jour avec succès',
                data: technicien
            });
            
        } catch (error) {
            if (error.message === 'Technicien non trouvé') {
                return res.status(404).json({ error: error.message });
            }
            next(error);
        }
    }

    /**
     * Supprime un technicien
     */
    async deleteTechnicien(req, res, next) {
        try {
            const { id } = req.params;
            
            const deleted = await technicienService.deleteTechnicien(id);
            if (!deleted) {
                return res.status(404).json({
                    error: 'Technicien non trouvé'
                });
            }
            
            logger.info(`Technicien supprimé: ${id}`);
            
            res.status(204).send();
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Vérifie la disponibilité d'un technicien
     */
    async checkAvailability(req, res, next) {
        try {
            const { id } = req.params;
            const { date } = req.query;
            
            if (!date) {
                return res.status(400).json({
                    error: 'La date est requise'
                });
            }

            const available = await technicienService.checkAvailability(id, new Date(date));
            
            res.json({
                technicien_id: parseInt(id),
                date: date,
                available: available
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Assigne un technicien à une intervention
     */
    async assignToIntervention(req, res, next) {
        try {
            const { id } = req.params;
            const { intervention_id } = req.body;
            
            if (!intervention_id) {
                return res.status(400).json({
                    error: 'L\'ID de l\'intervention est requis'
                });
            }

            await technicienService.assignToIntervention(id, intervention_id);
            
            res.json({
                message: 'Technicien assigné avec succès',
                technicien_id: parseInt(id),
                intervention_id: intervention_id
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Récupère les interventions d'un technicien
     */
    async getTechnicienInterventions(req, res, next) {
        try {
            const { id } = req.params;
            
            const interventions = await technicienService.getTechnicienInterventions(id);
            
            res.json({
                message: 'Interventions du technicien récupérées',
                technicien_id: parseInt(id),
                data: interventions
            });
            
        } catch (error) {
            next(error);
        }
    }

    /**
     * Récupère les statistiques d'un technicien
     */
    async getTechnicienStats(req, res, next) {
        try {
            const { id } = req.params;
            
            const stats = await technicienService.getTechnicienStats(id);
            
            res.json({
                message: 'Statistiques du technicien récupérées',
                technicien_id: parseInt(id),
                data: stats
            });
            
        } catch (error) {
            next(error);
        }
    }
}

module.exports = new TechnicienController();