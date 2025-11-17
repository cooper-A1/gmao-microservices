package com.ics.gmao.machines.service;

import com.ics.gmao.machines.dto.MachineDTO;
import com.ics.gmao.machines.model.EtatMachine;
import com.ics.gmao.machines.model.Machine;
import com.ics.gmao.machines.repository.MachineRepository;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * Service métier pour la gestion des machines
 */
@Service
public class MachineService {
    
    @Autowired
    private MachineRepository machineRepository;
    
    @Autowired
    private InterventionServiceClient interventionServiceClient;
    
    /**
     * Récupère toutes les machines
     */
    public List<MachineDTO> getAllMachines() {
        return machineRepository.findAll()
                .stream()
                .map(this::convertToDTO)
                .collect(Collectors.toList());
    }
    
    /**
     * Récupère une machine par ID
     */
    public Optional<MachineDTO> getMachineById(Long id) {
        return machineRepository.findById(id)
                .map(this::convertToDTO);
    }
    
    /**
     * Crée une nouvelle machine
     */
    public MachineDTO createMachine(MachineDTO machineDTO) {
        Machine machine = convertToEntity(machineDTO);
        machine = machineRepository.save(machine);
        return convertToDTO(machine);
    }
    
    /**
     * Met à jour une machine existante
     */
    public Optional<MachineDTO> updateMachine(Long id, MachineDTO machineDTO) {
        return machineRepository.findById(id)
                .map(machine -> {
                    // Mise à jour des champs (sauf ID, createdAt)
                    machine.setNom(machineDTO.getNom());
                    machine.setSite(machineDTO.getSite());
                    machine.setDateInstallation(machineDTO.getDateInstallation());
                    machine.setEtat(machineDTO.getEtat());
                    machine.setDescription(machineDTO.getDescription());
                    machine.setModele(machineDTO.getModele());
                    machine.setFabricant(machineDTO.getFabricant());
                    machine.setNumeroSerie(machineDTO.getNumeroSerie());
                    
                    Machine updatedMachine = machineRepository.save(machine);
                    return convertToDTO(updatedMachine);
                });
    }
    
    /**
     * Supprime une machine
     */
    public boolean deleteMachine(Long id) {
        if (machineRepository.existsById(id)) {
            machineRepository.deleteById(id);
            return true;
        }
        return false;
    }
    
    /**
     * Recherche des machines avec filtres
     */
    public Page<MachineDTO> searchMachines(String site, EtatMachine etat, String nom, Pageable pageable) {
        return machineRepository.findWithFilters(site, etat, nom, pageable)
                .map(this::convertToDTO);
    }
    
    /**
     * Récupère l'historique des interventions pour une machine
     */
    public Object getInterventionHistory(Long machineId) {
        // Vérifier que la machine existe
        if (!machineRepository.existsById(machineId)) {
            return null;
        }
        
        // Appel au service interventions via client HTTP
        return interventionServiceClient.getInterventionsByMachine(machineId);
    }
    
    /**
     * Conversion entité vers DTO
     */
    private MachineDTO convertToDTO(Machine machine) {
        MachineDTO dto = new MachineDTO();
        BeanUtils.copyProperties(machine, dto);
        return dto;
    }
    
    /**
     * Conversion DTO vers entité
     */
    private Machine convertToEntity(MachineDTO dto) {
        Machine machine = new Machine();
        BeanUtils.copyProperties(dto, machine);
        return machine;
    }
}
