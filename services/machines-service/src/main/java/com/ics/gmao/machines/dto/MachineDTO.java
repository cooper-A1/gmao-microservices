package com.ics.gmao.machines.dto;

import com.ics.gmao.machines.model.EtatMachine;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

import java.time.LocalDateTime;

/**
 * DTO pour les opérations sur les machines
 */
public class MachineDTO {
    private Long id;
    
    @NotBlank(message = "Le nom de la machine est obligatoire")
    private String nom;
    
    @NotBlank(message = "Le site est obligatoire")
    private String site;
    
    @NotNull(message = "La date d'installation est obligatoire")
    private LocalDateTime dateInstallation;
    
    private EtatMachine etat;
    private String description;
    private String modele;
    private String fabricant;
    private String numeroSerie;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
    
    // Constructeurs
    public MachineDTO() {}
    
    // Getters et Setters (identiques au modèle)
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    
    public String getNom() { return nom; }
    public void setNom(String nom) { this.nom = nom; }
    
    public String getSite() { return site; }
    public void setSite(String site) { this.site = site; }
    
    public LocalDateTime getDateInstallation() { return dateInstallation; }
    public void setDateInstallation(LocalDateTime dateInstallation) { this.dateInstallation = dateInstallation; }
    
    public EtatMachine getEtat() { return etat; }
    public void setEtat(EtatMachine etat) { this.etat = etat; }
    
    public String getDescription() { return description; }
    public void setDescription(String description) { this.description = description; }
    
    public String getModele() { return modele; }
    public void setModele(String modele) { this.modele = modele; }
    
    public String getFabricant() { return fabricant; }
    public void setFabricant(String fabricant) { this.fabricant = fabricant; }
    
    public String getNumeroSerie() { return numeroSerie; }
    public void setNumeroSerie(String numeroSerie) { this.numeroSerie = numeroSerie; }
    
    public LocalDateTime getCreatedAt() { return createdAt; }
    public void setCreatedAt(LocalDateTime createdAt) { this.createdAt = createdAt; }
    
    public LocalDateTime getUpdatedAt() { return updatedAt; }
    public void setUpdatedAt(LocalDateTime updatedAt) { this.updatedAt = updatedAt; }
}