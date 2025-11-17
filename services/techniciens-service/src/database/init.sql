-- Script d'initialisation de la base de données MySQL pour les techniciens

CREATE DATABASE IF NOT EXISTS techniciens_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE techniciens_db;

-- Table des techniciens
CREATE TABLE IF NOT EXISTS techniciens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nom VARCHAR(100) NOT NULL,
    prenom VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telephone VARCHAR(20) NOT NULL,
    competences JSON NOT NULL COMMENT 'Liste des compétences au format JSON',
    niveau_experience ENUM('junior', 'senior', 'expert') DEFAULT 'junior',
    disponibilite BOOLEAN DEFAULT TRUE,
    salaire DECIMAL(10, 2) NULL,
    date_embauche DATE NULL,
    notes TEXT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_nom_prenom (nom, prenom),
    INDEX idx_email (email),
    INDEX idx_disponibilite (disponibilite),
    INDEX idx_niveau_experience (niveau_experience),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB;

-- Table de liaison technicien-interventions (pour l'historique)
CREATE TABLE IF NOT EXISTS technicien_interventions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    technicien_id INT NOT NULL,
    intervention_id VARCHAR(50) NOT NULL COMMENT 'ID de l\'intervention du service interventions',
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (technicien_id) REFERENCES techniciens(id) ON DELETE CASCADE,
    UNIQUE KEY unique_assignment (technicien_id, intervention_id),
    INDEX idx_technicien (technicien_id),
    INDEX idx_intervention (intervention_id),
    INDEX idx_assigned_at (assigned_at)
) ENGINE=InnoDB;

-- Insertion de données de test
INSERT INTO techniciens (nom, prenom, email, telephone, competences, niveau_experience, disponibilite, salaire, date_embauche, notes) VALUES
('DIOP', 'Amadou', 'amadou.diop@ics.sn', '+221 77 123 45 67', '["Mécanique générale", "Hydraulique", "Maintenance préventive"]', 'senior', TRUE, 450000, '2020-03-15', 'Technicien expérimenté en maintenance hydraulique'),

('FALL', 'Fatou', 'fatou.fall@ics.sn', '+221 78 987 65 43', '["Électricité industrielle", "Automatisme", "Électronique"]', 'expert', TRUE, 550000, '2018-01-10', 'Spécialiste en systèmes électriques et automatisme'),

('NDIAYE', 'Moussa', 'moussa.ndiaye@ics.sn', '+221 76 555 12 34', '["Pneumatique", "Mécanique générale", "Soudure"]', 'junior', TRUE, 320000, '2023-06-01', 'Nouveau technicien, en formation continue'),

('SECK', 'Aissatou', 'aissatou.seck@ics.sn', '+221 77 888 99 00', '["Informatique industrielle", "Régulation", "Maintenance préventive"]', 'senior', TRUE, 480000, '2019-09-20', 'Experte en systèmes informatiques industriels'),

('KANE', 'Ibrahima', 'ibrahima.kane@ics.sn', '+221 78 111 22 33', '["Usinage", "Mécanique générale", "Diagnostic de pannes"]', 'senior', FALSE, 420000, '2021-02-28', 'Actuellement en formation spécialisée');

-- Insertion de quelques assignations de test
INSERT INTO technicien_interventions (technicien_id, intervention_id, assigned_at) VALUES
(1, 'test-intervention-1', '2024-01-15 09:00:00'),
(1, 'test-intervention-2', '2024-01-20 14:30:00'),
(2, 'test-intervention-3', '2024-01-18 11:00:00'),
(3, 'test-intervention-4', '2024-01-22 08:15:00'),
(2, 'test-intervention-5', '2024-01-25 16:45:00');

-- Procédure pour récupérer les statistiques d'un technicien
DELIMITER $$
CREATE PROCEDURE GetTechnicienStats(IN tech_id INT)
BEGIN
    SELECT 
        t.id,
        t.nom,
        t.prenom,
        t.niveau_experience,
        COUNT(ti.id) as total_interventions,
        COUNT(CASE WHEN ti.assigned_at >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 1 END) as interventions_30j,
        COUNT(CASE WHEN ti.assigned_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 END) as interventions_7j
    FROM techniciens t
    LEFT JOIN technicien_interventions ti ON t.id = ti.technicien_id
    WHERE t.id = tech_id
    GROUP BY t.id;
END$$
DELIMITER ;

-- Vue pour les techniciens disponibles
CREATE VIEW techniciens_disponibles AS
SELECT 
    id,
    nom,
    prenom,
    email,
    telephone,
    competences,
    niveau_experience,
    date_embauche,
    (SELECT COUNT(*) FROM technicien_interventions ti WHERE ti.technicien_id = techniciens.id) as nb_interventions
FROM techniciens
WHERE disponibilite = TRUE;

-- Index full-text pour la recherche
-- ALTER TABLE techniciens ADD FULLTEXT(nom, prenom, notes);

-- Création de l'utilisateur pour l'application
CREATE USER IF NOT EXISTS 'gmao_user'@'%' IDENTIFIED BY 'gmao_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON techniciens_db.* TO 'gmao_user'@'%';
FLUSH PRIVILEGES;

-- Fin du script d'initialisation