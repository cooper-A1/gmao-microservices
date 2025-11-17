/**
 * Configuration Swagger/OpenAPI
 */

const swaggerJsdoc = require('swagger-jsdoc');

const options = {
    definition: {
        openapi: '3.0.0',
        info: {
            title: 'Service Techniciens GMAO',
            version: '1.0.0',
            description: 'API de gestion des techniciens pour le système GMAO d\'ICS',
            contact: {
                name: 'ICS GMAO Team',
                email: 'gmao@ics.sn'
            }
        },
        servers: [
            {
                url: process.env.API_URL || 'http://localhost:8003',
                description: 'Serveur de développement'
            }
        ],
        components: {
            securitySchemes: {
                bearerAuth: {
                    type: 'http',
                    scheme: 'bearer',
                    bearerFormat: 'JWT'
                }
            }
        },
        security: [
            {
                bearerAuth: []
            }
        ]
    },
    apis: ['./src/routes/*.js'], // Chemin vers les fichiers contenant les annotations
};

const specs = swaggerJsdoc(options);

module.exports = specs;