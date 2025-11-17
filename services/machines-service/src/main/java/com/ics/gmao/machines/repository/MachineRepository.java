package com.ics.gmao.machines.repository;

import com.ics.gmao.machines.model.EtatMachine;
import com.ics.gmao.machines.model.Machine;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;

/**
 * Repository pour l'accès aux données des machines
 */
@Repository
public interface MachineRepository extends JpaRepository<Machine, Long> {
    
    /**
     * Recherche des machines par site
     */
    List<Machine> findBySite(String site);
    
    /**
     * Recherche des machines par état
     */
    List<Machine> findByEtat(EtatMachine etat);
    
    /**
     * Recherche des machines par site et état
     */
    List<Machine> findBySiteAndEtat(String site, EtatMachine etat);
    
    /**
     * Recherche des machines par nom (contient)
     */
    @Query("SELECT m FROM Machine m WHERE LOWER(m.nom) LIKE LOWER(CONCAT('%', :nom, '%'))")
    List<Machine> findByNomContainingIgnoreCase(@Param("nom") String nom);
    
    /**
     * Recherche paginée avec filtres
     */
    @Query("SELECT m FROM Machine m WHERE " +
           "(:site IS NULL OR m.site = :site) AND " +
           "(:etat IS NULL OR m.etat = :etat) AND " +
           "(:nom IS NULL OR LOWER(m.nom) LIKE LOWER(CONCAT('%', :nom, '%')))")
    Page<Machine> findWithFilters(@Param("site") String site, 
                                 @Param("etat") EtatMachine etat, 
                                 @Param("nom") String nom, 
                                 Pageable pageable);
}