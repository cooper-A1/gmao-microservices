-- Script d'initialisation pour PostgreSQL (Machines)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insertion de données de test
INSERT INTO machines (nom, site, date_installation, etat, description, modele, fabricant, numero_serie) VALUES
('Compresseur Atlas Copco', 'Usine Dakar', '2020-05-15 10:00:00', 'OPERATIONNELLE', 'Compresseur d''air principal', 'GA 22', 'Atlas Copco', 'AC2020001'),
('Fraiseuse CNC Haas', 'Atelier Mécanique', '2019-03-10 14:30:00', 'OPERATIONNELLE', 'Fraiseuse à commande numérique', 'VF-2', 'Haas Automation', 'HAAS2019001'),
('Pompe hydraulique Bosch', 'Station Hydraulique', '2021-08-20 09:15:00', 'EN_MAINTENANCE', 'Pompe principale circuit hydraulique', 'A10VSO', 'Bosch Rexroth', 'BR2021001'),
('Convoyeur transporteur', 'Zone Expédition', '2022-01-05 16:45:00', 'OPERATIONNELLE', 'Convoyeur à bande principale', 'CB-1200', 'Flexco', 'FLX2022001');

---