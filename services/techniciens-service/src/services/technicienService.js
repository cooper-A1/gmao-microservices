/**
 * Service métier pour la gestion des techniciens
 */

const database = require('../database/connection');
const logger = require('../utils/logger');
const moment = require('moment');

class TechnicienService {
    
    /**
     * Crée un nouveau technicien
     */
    async createTechnicien(technicienData) {
        try {
            const {
                nom, prenom, email, telephone, competences,
                niveau_experience = 'junior', disponibilite = true,
                salaire, date_embauche, notes
            } = technicienData;

            // Vérification de l'unicité de l'email
            const existingTechnicien = await this.findByEmail(email);
            if (existingTechnicien) {
                throw new Error('Un technicien avec cet email existe déjà');
            }

            const sql = `
                INSERT INTO techniciens 
                (nom, prenom, email, telephone, competences, niveau_experience, 
                 disponibilite, salaire, date_embauche, notes, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
            `;

            const params = [
                nom, prenom, email, telephone, 
                JSON.stringify(competences), niveau_experience, 
                disponibilite, salaire, date_embauche, notes
            ];

            const result = await database.query(sql, params);
            
            // Récupération du technicien créé
            return await this.findById(result.insertId);
            
        } catch (error) {
            logger.error('Erreur lors de la création du technicien:', error);
            throw error;
        }
    }

    /**
     * Récupère tous les techniciens avec filtres optionnels
     */
    async getAllTechniciens(filters = {}) {
        try {
            let sql = 'SELECT * FROM techniciens WHERE 1=1';
            const params = [];

            // Application des filtres
            if (filters.disponibilite !== undefined) {
                sql += ' AND disponibilite = ?';
                params.push(filters.disponibilite);
            }

            if (filters.niveau_experience) {
                sql += ' AND niveau_experience = ?';
                params.push(filters.niveau_experience);
            }

            if (filters.competence) {
                sql += ' AND JSON_CONTAINS(competences, ?)';
                params.push(`"${filters.competence}"`);
            }

            if (filters.search) {
                sql += ' AND (nom LIKE ? OR prenom LIKE ? OR email LIKE ?)';
                const searchTerm = `%${filters.search}%`;
                params.push(searchTerm, searchTerm, searchTerm);
            }

            sql += ' ORDER BY nom ASC, prenom ASC';

            // Pagination
            if (filters.limit) {
                sql += ' LIMIT ?';
                params.push(parseInt(filters.limit));
                
                if (filters.offset) {
                    sql += ' OFFSET ?';
                    params.push(parseInt(filters.offset));
                }
            }

            const techniciens = await database.query(sql, params);
            
            // Parsing des compétences JSON
            return techniciens.map(this.formatTechnicien);
            
        } catch (error) {
            logger.error('Erreur lors de la récupération des techniciens:', error);
            throw error;
        }
    }

    /**
     * Récupère un technicien par ID
     */
    async findById(id) {
        try {
            const sql = 'SELECT * FROM techniciens WHERE id = ?';
            const [technicien] = await database.query(sql, [id]);
            
            return technicien ? this.formatTechnicien(technicien) : null;
            
        } catch (error) {
            logger.error('Erreur lors de la récupération du technicien:', error);
            throw error;
        }
    }

    /**
     * Récupère un technicien par email
     */
    async findByEmail(email) {
        try {
            const sql = 'SELECT * FROM techniciens WHERE email = ?';
            const [technicien] = await database.query(sql, [email]);
            
            return technicien ? this.formatTechnicien(technicien) : null;
            
        } catch (error) {
            logger.error('Erreur lors de la recherche par email:', error);
            throw error;
        }
    }

    /**
     * Met à jour un technicien
     */
    async updateTechnicien(id, updateData) {
        try {
            // Vérification de l'existence
            const existingTechnicien = await this.findById(id);
            if (!existingTechnicien) {
                throw new Error('Technicien non trouvé');
            }

            // Vérification de l'unicité de l'email si modifié
            if (updateData.email && updateData.email !== existingTechnicien.email) {
                const emailExists = await this.findByEmail(updateData.email);
                if (emailExists) {
                    throw new Error('Un technicien avec cet email existe déjà');
                }
            }

            const fields = [];
            const params = [];

            // Construction dynamique de la requête
            Object.keys(updateData).forEach(key => {
                if (updateData[key] !== undefined) {
                    fields.push(`${key} = ?`);
                    
                    // Sérialisation JSON pour les compétences
                    if (key === 'competences') {
                        params.push(JSON.stringify(updateData[key]));
                    } else {
                        params.push(updateData[key]);
                    }
                }
            });

            if (fields.length === 0) {
                throw new Error('Aucune donnée à mettre à jour');
            }

            fields.push('updated_at = NOW()');
            params.push(id);

            const sql = `UPDATE techniciens SET ${fields.join(', ')} WHERE id = ?`;
            
            await database.query(sql, params);
            
            return await this.findById(id);
            
        } catch (error) {
            logger.error('Erreur lors de la mise à jour du technicien:', error);
            throw error;
        }
    }

    /**
     * Supprime un technicien
     */
    async deleteTechnicien(id) {
        try {
            const sql = 'DELETE FROM techniciens WHERE id = ?';
            const result = await database.query(sql, [id]);
            
            return result.affectedRows > 0;
            
        } catch (error) {
            logger.error('Erreur lors de la suppression du technicien:', error);
            throw error;
        }
    }

    /**
     * Vérifie la disponibilité d'un technicien à une date donnée
     */
    async checkAvailability(technicienId, dateIntervention) {
        try {
            const technicien = await this.findById(technicienId);
            if (!technicien) {
                return false;
            }

            // Vérification de la disponibilité générale
            if (!technicien.disponibilite) {
                return false;
            }

            // TODO: Vérifier les interventions planifiées
            // Pour l'instant, on suppose que le technicien est disponible
            // En production, on vérifierait les conflits d'horaires
            
            return true;
            
        } catch (error) {
            logger.error('Erreur lors de la vérification de disponibilité:', error);
            return false;
        }
    }

    /**
     * Assigne un technicien à une intervention
     */
    async assignToIntervention(technicienId, interventionId) {
        try {
            // Enregistrement de l'assignation (pour historique)
            const sql = `
                INSERT INTO technicien_interventions (technicien_id, intervention_id, assigned_at)
                VALUES (?, ?, NOW())
                ON DUPLICATE KEY UPDATE assigned_at = NOW()
            `;
            
            await database.query(sql, [technicienId, interventionId]);
            
            logger.info(`Technicien ${technicienId} assigné à l'intervention ${interventionId}`);
            return true;
            
        } catch (error) {
            logger.error('Erreur lors de l\'assignation:', error);
            throw error;
        }
    }

    /**
     * Récupère les interventions d'un technicien
     */
    async getTechnicienInterventions(technicienId) {
        try {
            const sql = `
                SELECT intervention_id, assigned_at 
                FROM technicien_interventions 
                WHERE technicien_id = ? 
                ORDER BY assigned_at DESC
            `;
            
            return await database.query(sql, [technicienId]);
            
        } catch (error) {
            logger.error('Erreur lors de la récupération des interventions:', error);
            throw error;
        }
    }

    /**
     * Récupère les statistiques d'un technicien
     */
    async getTechnicienStats(technicienId) {
        try {
            const sql = `
                SELECT 
                    COUNT(*) as total_interventions,
                    COUNT(CASE WHEN assigned_at >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 END) as interventions_30j
                FROM technicien_interventions 
                WHERE technicien_id = ?
            `;
            
            const [stats] = await database.query(sql, [technicienId]);
            
            return {
                total_interventions: stats?.total_interventions || 0,
                interventions_30_jours: stats?.interventions_30j || 0,
                taux_activite: this.calculateActivityRate(stats?.interventions_30j || 0)
            };
            
        } catch (error) {
            logger.error('Erreur lors du calcul des statistiques:', error);
            throw error;
        }
    }

    /**
     * Calcule le taux d'activité
     */
    calculateActivityRate(interventions30j) {
        // Logique simple: plus de 10 interventions = actif
        if (interventions30j >= 10) return 'élevé';
        if (interventions30j >= 5) return 'moyen';
        return 'faible';
    }

    /**
     * Formate un objet technicien (parse JSON, etc.)
     */
    formatTechnicien(technicien) {
        return {
            ...technicien,
            competences: typeof technicien.competences === 'string' 
                ? JSON.parse(technicien.competences) 
                : technicien.competences,
            created_at: moment(technicien.created_at).format(),
            updated_at: moment(technicien.updated_at).format()
        };
    }
}

module.exports = new TechnicienService();
