/**
 * Routes pour la gestion des techniciens
 */

const express = require('express');
const technicienController = require('../controllers/technicienController');
const { authMiddleware, requireRole } = require('../middleware/auth');
const { validateId } = require('../middleware/validation');

const router = express.Router();

/**
 * @swagger
 * components:
 *   schemas:
 *     Technicien:
 *       type: object
 *       required:
 *         - nom
 *         - prenom
 *         - email
 *         - telephone
 *         - competences
 *       properties:
 *         id:
 *           type: integer
 *           description: ID unique du technicien
 *         nom:
 *           type: string
 *           description: Nom de famille du technicien
 *         prenom:
 *           type: string
 *           description: Prénom du technicien
 *         email:
 *           type: string
 *           format: email
 *           description: Adresse email
 *         telephone:
 *           type: string
 *           description: Numéro de téléphone
 *         competences:
 *           type: array
 *           items:
 *             type: string
 *           description: Liste des compétences
 *         niveau_experience:
 *           type: string
 *           enum: [junior, senior, expert]
 *           description: Niveau d'expérience
 *         disponibilite:
 *           type: boolean
 *           description: Disponibilité du technicien
 *         salaire:
 *           type: number
 *           description: Salaire (optionnel)
 *         date_embauche:
 *           type: string
 *           format: date
 *           description: Date d'embauche
 *         notes:
 *           type: string
 *           description: Notes additionnelles
 */

/**
 * @swagger
 * /api/techniciens:
 *   get:
 *     summary: Récupère la liste des techniciens
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: query
 *         name: disponibilite
 *         schema:
 *           type: boolean
 *         description: Filtrer par disponibilité
 *       - in: query
 *         name: niveau_experience
 *         schema:
 *           type: string
 *           enum: [junior, senior, expert]
 *         description: Filtrer par niveau d'expérience
 *       - in: query
 *         name: competence
 *         schema:
 *           type: string
 *         description: Filtrer par compétence
 *       - in: query
 *         name: search
 *         schema:
 *           type: string
 *         description: Recherche dans nom, prénom, email
 *       - in: query
 *         name: limit
 *         schema:
 *           type: integer
 *         description: Nombre limite de résultats
 *       - in: query
 *         name: offset
 *         schema:
 *           type: integer
 *         description: Décalage pour la pagination
 *     responses:
 *       200:
 *         description: Liste des techniciens
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 message:
 *                   type: string
 *                 data:
 *                   type: array
 *                   items:
 *                     $ref: '#/components/schemas/Technicien'
 *                 count:
 *                   type: integer
 */
router.get('/', authMiddleware, technicienController.getAllTechniciens);

/**
 * @swagger
 * /api/techniciens:
 *   post:
 *     summary: Crée un nouveau technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Technicien'
 *     responses:
 *       201:
 *         description: Technicien créé avec succès
 *       400:
 *         description: Données invalides
 */
router.post('/', authMiddleware, requireRole(['admin', 'manager']), technicienController.createTechnicien);

/**
 * @swagger
 * /api/techniciens/{id}:
 *   get:
 *     summary: Récupère un technicien par ID
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *         description: ID du technicien
 *     responses:
 *       200:
 *         description: Technicien trouvé
 *       404:
 *         description: Technicien non trouvé
 */
router.get('/:id', authMiddleware, validateId, technicienController.getTechnicienById);

/**
 * @swagger
 * /api/techniciens/{id}:
 *   put:
 *     summary: Met à jour un technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/Technicien'
 *     responses:
 *       200:
 *         description: Technicien mis à jour
 *       404:
 *         description: Technicien non trouvé
 */
router.put('/:id', authMiddleware, requireRole(['admin', 'manager']), validateId, technicienController.updateTechnicien);

/**
 * @swagger
 * /api/techniciens/{id}:
 *   delete:
 *     summary: Supprime un technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     responses:
 *       204:
 *         description: Technicien supprimé
 *       404:
 *         description: Technicien non trouvé
 */
router.delete('/:id', authMiddleware, requireRole(['admin']), validateId, technicienController.deleteTechnicien);

/**
 * @swagger
 * /api/techniciens/{id}/availability:
 *   get:
 *     summary: Vérifie la disponibilité d'un technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *       - in: query
 *         name: date
 *         required: true
 *         schema:
 *           type: string
 *           format: date-time
 *         description: Date et heure de l'intervention
 *     responses:
 *       200:
 *         description: Statut de disponibilité
 */
router.get('/:id/availability', authMiddleware, validateId, technicienController.checkAvailability);

/**
 * @swagger
 * /api/techniciens/{id}/assign:
 *   post:
 *     summary: Assigne un technicien à une intervention
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - intervention_id
 *             properties:
 *               intervention_id:
 *                 type: string
 *                 description: ID de l'intervention
 *     responses:
 *       200:
 *         description: Assignation réussie
 */
router.post('/:id/assign', authMiddleware, requireRole(['admin', 'manager']), validateId, technicienController.assignToIntervention);

/**
 * @swagger
 * /api/techniciens/{id}/interventions:
 *   get:
 *     summary: Récupère les interventions d'un technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     responses:
 *       200:
 *         description: Liste des interventions
 */
router.get('/:id/interventions', authMiddleware, validateId, technicienController.getTechnicienInterventions);

/**
 * @swagger
 * /api/techniciens/{id}/stats:
 *   get:
 *     summary: Récupère les statistiques d'un technicien
 *     tags: [Techniciens]
 *     security:
 *       - bearerAuth: []
 *     parameters:
 *       - in: path
 *         name: id
 *         required: true
 *         schema:
 *           type: integer
 *     responses:
 *       200:
 *         description: Statistiques du technicien
 */
router.get('/:id/stats', authMiddleware, validateId, technicienController.getTechnicienStats);

module.exports = router;