/**
 * Routes d'authentification
 */

const express = require('express');
const jwt = require('jsonwebtoken');
const bcrypt = require('bcryptjs');
const Joi = require('joi');

const router = express.Router();

// Utilisateurs de démo (en production, je vais utiliser une vraie base de données)
const demoUsers = [
    {
        id: 1,
        username: 'admin',
        email: 'admin@ics.sn',
        password: bcrypt.hashSync('admin123', 10),
        role: 'admin'
    },
    {
        id: 2,
        username: 'manager',
        email: 'manager@ics.sn',
        password: bcrypt.hashSync('manager123', 10),
        role: 'manager'
    },
    {
        id: 3,
        username: 'tech1',
        email: 'tech1@ics.sn',
        password: bcrypt.hashSync('tech123', 10),
        role: 'technicien'
    }
];

// Schéma de validation pour la connexion
const loginSchema = Joi.object({
    username: Joi.string().required(),
    password: Joi.string().required()
});

/**
 * @swagger
 * /api/auth/login:
 *   post:
 *     summary: Connexion utilisateur
 *     tags: [Authentication]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             type: object
 *             required:
 *               - username
 *               - password
 *             properties:
 *               username:
 *                 type: string
 *               password:
 *                 type: string
 *     responses:
 *       200:
 *         description: Connexion réussie
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 access_token:
 *                   type: string
 *                 token_type:
 *                   type: string
 *                 user_info:
 *                   type: object
 *       401:
 *         description: Identifiants incorrects
 */
router.post('/login', async (req, res) => {
    try {
        // Validation des données
        const { error, value } = loginSchema.validate(req.body);
        if (error) {
            return res.status(400).json({
                error: 'Données invalides',
                details: error.details.map(d => d.message)
            });
        }

        const { username, password } = value;

        // Recherche de l'utilisateur
        const user = demoUsers.find(u => u.username === username);
        if (!user) {
            return res.status(401).json({
                error: 'Nom d\'utilisateur ou mot de passe incorrect'
            });
        }

        // Vérification du mot de passe
        const isValidPassword = bcrypt.compareSync(password, user.password);
        if (!isValidPassword) {
            return res.status(401).json({
                error: 'Nom d\'utilisateur ou mot de passe incorrect'
            });
        }

        // Génération du token JWT
        const token = jwt.sign(
            {
                sub: user.username,
                user_id: user.id,
                role: user.role
            },
            process.env.JWT_SECRET || 'your-secret-key',
            { expiresIn: '24h' }
        );

        res.json({
            access_token: token,
            token_type: 'bearer',
            user_info: {
                id: user.id,
                username: user.username,
                email: user.email,
                role: user.role
            }
        });

    } catch (error) {
        res.status(500).json({
            error: 'Erreur interne du serveur'
        });
    }
});

module.exports = router;