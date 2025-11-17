// Script d'initialisation pour MongoDB (Interventions)
db = db.getSiblingDB('interventions_db');

// Insertion de données de test
db.interventions.insertMany([
    {
        machine_id: 1,
        type_intervention: "preventive",
        titre: "Maintenance préventive compresseur",
        description: "Vidange huile et changement filtres",
        technicien_id: 1,
        statut: "terminee",
        date_planifiee: new Date("2024-01-15T09:00:00Z"),
        date_creation: new Date("2024-01-10T10:00:00Z"),
        date_debut: new Date("2024-01-15T09:00:00Z"),
        date_fin: new Date("2024-01-15T11:30:00Z"),
        duree_estimee: 120,
        duree_reelle: 150,
        priorite: 2,
        compte_rendu: "Maintenance effectuée selon planning. Huile changée, filtres remplacés.",
        pieces_utilisees: [
            {
                piece_id: "piece-003",
                nom: "Huile hydraulique ISO 68",
                quantite: 5,
                prix_unitaire: 8.90
            }
        ],
        cout_total: 44.50
    },
    {
        machine_id: 2,
        type_intervention: "corrective", 
        titre: "Réparation fraiseuse CNC",
        description: "Panne moteur broche principale",
        technicien_id: 2,
        statut: "en_cours",
        date_planifiee: new Date("2024-01-20T08:00:00Z"),
        date_creation: new Date("2024-01-18T14:30:00Z"),
        date_debut: new Date("2024-01-20T08:15:00Z"),
        duree_estimee: 240,
        priorite: 4
    }
]);

console.log("✅ Données de test insérées dans MongoDB");