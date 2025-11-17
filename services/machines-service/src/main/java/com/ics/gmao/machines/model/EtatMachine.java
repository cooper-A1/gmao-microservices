package com.ics.gmao.machines.model;

/**
 * Énumération des états possibles d'une machine
 */
public enum EtatMachine {
    OPERATIONNELLE("Opérationnelle"),
    EN_PANNE("En panne"),
    EN_MAINTENANCE("En maintenance"),
    ARRETEE("Arrêtée"),
    HORS_SERVICE("Hors service");
    
    private final String libelle;
    
    EtatMachine(String libelle) {
        this.libelle = libelle;
    }
    
    public String getLibelle() {
        return libelle;
    }
}