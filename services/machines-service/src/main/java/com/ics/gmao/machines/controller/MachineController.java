package com.ics.gmao.machines.controller;

import com.ics.gmao.machines.dto.MachineDTO;
import com.ics.gmao.machines.model.EtatMachine;
import com.ics.gmao.machines.service.MachineService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

/**
 * Contrôleur REST pour la gestion des machines
 */
@RestController
@RequestMapping("/api/machines")
@Tag(name = "Machines", description = "API de gestion des machines industrielles")
@CrossOrigin(origins = "*")
public class MachineController {
    
    @Autowired
    private MachineService machineService;
    
    /**
     * Récupère toutes les machines
     */
    @GetMapping
    @Operation(summary = "Récupérer toutes les machines", description = "Retourne la liste de toutes les machines")
    @ApiResponse(responseCode = "200", description = "Liste des machines récupérée avec succès")
    public ResponseEntity<List<MachineDTO>> getAllMachines() {
        List<MachineDTO> machines = machineService.getAllMachines();
        return ResponseEntity.ok(machines);
    }
    
    /**
     * Récupère une machine par ID
     */
    @GetMapping("/{id}")
    @Operation(summary = "Récupérer une machine par ID", description = "Retourne les détails d'une machine")
    @ApiResponse(responseCode = "200", description = "Machine trouvée")
    @ApiResponse(responseCode = "404", description = "Machine non trouvée")
    public ResponseEntity<MachineDTO> getMachineById(@PathVariable Long id) {
        return machineService.getMachineById(id)
                .map(machine -> ResponseEntity.ok(machine))
                .orElse(ResponseEntity.notFound().build());
    }
    
    /**
     * Crée une nouvelle machine
     */
    @PostMapping
    @Operation(summary = "Créer une nouvelle machine", description = "Ajoute une nouvelle machine au parc")
    @ApiResponse(responseCode = "201", description = "Machine créée avec succès")
    @ApiResponse(responseCode = "400", description = "Données invalides")
    public ResponseEntity<MachineDTO> createMachine(@Valid @RequestBody MachineDTO machineDTO) {
        MachineDTO createdMachine = machineService.createMachine(machineDTO);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdMachine);
    }
    
    /**
     * Met à jour une machine existante
     */
    @PutMapping("/{id}")
    @Operation(summary = "Mettre à jour une machine", description = "Met à jour les informations d'une machine existante")
    @ApiResponse(responseCode = "200", description = "Machine mise à jour avec succès")
    @ApiResponse(responseCode = "404", description = "Machine non trouvée")
    @ApiResponse(responseCode = "400", description = "Données invalides")
    public ResponseEntity<MachineDTO> updateMachine(@PathVariable Long id, 
                                                   @Valid @RequestBody MachineDTO machineDTO) {
        return machineService.updateMachine(id, machineDTO)
                .map(machine -> ResponseEntity.ok(machine))
                .orElse(ResponseEntity.notFound().build());
    }
    
    /**
     * Supprime une machine
     */
    @DeleteMapping("/{id}")
    @Operation(summary = "Supprimer une machine", description = "Supprime une machine du parc")
    @ApiResponse(responseCode = "204", description = "Machine supprimée avec succès")
    @ApiResponse(responseCode = "404", description = "Machine non trouvée")
    public ResponseEntity<Void> deleteMachine(@PathVariable Long id) {
        if (machineService.deleteMachine(id)) {
            return ResponseEntity.noContent().build();
        }
        return ResponseEntity.notFound().build();
    }
    
    /**
     * Recherche des machines avec filtres et pagination
     */
    @GetMapping("/search")
    @Operation(summary = "Rechercher des machines", description = "Recherche des machines avec filtres et pagination")
    public ResponseEntity<Page<MachineDTO>> searchMachines(
            @RequestParam(required = false) String site,
            @RequestParam(required = false) EtatMachine etat,
            @RequestParam(required = false) String nom,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "id") String sortBy,
            @RequestParam(defaultValue = "ASC") Sort.Direction direction) {
        
        Pageable pageable = PageRequest.of(page, size, Sort.by(direction, sortBy));
        Page<MachineDTO> machines = machineService.searchMachines(site, etat, nom, pageable);
        return ResponseEntity.ok(machines);
    }
    
    /**
     * Récupère l'historique des interventions pour une machine
     */
    @GetMapping("/{id}/interventions")
    @Operation(summary = "Historique des interventions", description = "Récupère l'historique des interventions pour une machine")
    @ApiResponse(responseCode = "200", description = "Historique récupéré avec succès")
    @ApiResponse(responseCode = "404", description = "Machine non trouvée")
    public ResponseEntity<Object> getInterventionHistory(@PathVariable Long id) {
        Object history = machineService.getInterventionHistory(id);
        if (history != null) {
            return ResponseEntity.ok(history);
        }
        return ResponseEntity.notFound().build();
    }
    
    /**
     * Endpoint de santé du service
     */
    @GetMapping("/health")
    @Operation(summary = "Vérification de santé", description = "Vérifie que le service est opérationnel")
    public ResponseEntity<String> health() {
        return ResponseEntity.ok("Service Machines opérationnel");
    }
}