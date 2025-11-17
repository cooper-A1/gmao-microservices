/**
 * Configuration et gestion de la connexion MySQL
 */

const mysql = require('mysql2/promise');
const logger = require('../utils/logger');

class Database {
    constructor() {
        this.pool = null;
        this.config = {
            host: process.env.DB_HOST || 'mysql-techniciens',
            port: process.env.DB_PORT || 3306,
            user: process.env.DB_USER || 'gmao_user',
            password: process.env.DB_PASSWORD || 'gmao_password',
            database: process.env.DB_NAME || 'techniciens_db',
            charset: 'utf8mb4',
            connectionLimit: 10,
            queueLimit: 0,
            acquireTimeout: 60000,
            timeout: 60000,
        };
    }

    async connect() {
        try {
            this.pool = mysql.createPool(this.config);
            logger.info('Pool de connexions MySQL créé');
        } catch (error) {
            logger.error('Erreur lors de la création du pool MySQL:', error);
            throw error;
        }
    }

    async testConnection() {
        try {
            if (!this.pool) {
                await this.connect();
            }
            
            const connection = await this.pool.getConnection();
            await connection.ping();
            connection.release();
            
            logger.info('✅ Connexion MySQL testée avec succès');
            return true;
        } catch (error) {
            logger.error('❌ Test de connexion MySQL échoué:', error);
            throw error;
        }
    }

    async query(sql, params = []) {
        try {
            if (!this.pool) {
                await this.connect();
            }
            
            const [rows] = await this.pool.execute(sql, params);
            return rows;
        } catch (error) {
            logger.error('Erreur lors de l\'exécution de la requête:', {
                sql,
                params,
                error: error.message
            });
            throw error;
        }
    }

    async transaction(queries) {
        const connection = await this.pool.getConnection();
        
        try {
            await connection.beginTransaction();
            
            const results = [];
            for (const { sql, params } of queries) {
                const [result] = await connection.execute(sql, params);
                results.push(result);
            }
            
            await connection.commit();
            return results;
        } catch (error) {
            await connection.rollback();
            throw error;
        } finally {
            connection.release();
        }
    }

    async close() {
        if (this.pool) {
            await this.pool.end();
            logger.info('Pool de connexions MySQL fermé');
        }
    }
}

module.exports = new Database();
