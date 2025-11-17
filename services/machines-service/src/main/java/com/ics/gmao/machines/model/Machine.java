package com.ics.gmao.machines.model;

import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;

/**
 * Entité représentant une machine industrielle
 */
@Entity
@Table(name = "machines")
public class Machine {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @NotBlank(message = "Le nom de la machine est obligatoire")
    @Column(nullable = false, length = 100)
    private String nom;
    
    @NotBlank(message = "Le site est obligatoire")
    @Column(nullable = false, length = 50)
    private String site;
    
    @NotNull(message = "La date d'installation est obligatoire")
    @Column(name = "date_installation", nullable = false)
    private LocalDateTime dateInstallation;
    
    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private EtatMachine etat = EtatMachine.OPERATIONNELLE;
    
    @Column(length = 500)
    private String description;
    
    @Column(length = 100)
    private String modele;
    
    @Column(length = 100)
    private String fabricant;
    
    @Column(name = "numero_serie", length = 100)
    private String numeroSerie;
    
    @CreationTimestamp
    @Column(name = "created_at")
    private LocalDateTime createdAt;
    
    @UpdateTimestamp
    @Column(name = "updated_at")
    private LocalDateTime updatedAt;
    
    // Constructeurs
    public Machine() {}
    
    public Machine(String nom, String site, LocalDateTime dateInstallation) {
        this.nom = nom;
        this.site = site;
        this.dateInstallation = dateInstallation;
    }
    
    // Getters et Setters
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
    public LocalDateTime getUpdatedAt() { return updatedAt; }
}
