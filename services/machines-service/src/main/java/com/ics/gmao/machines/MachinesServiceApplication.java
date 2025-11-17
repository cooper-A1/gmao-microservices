package com.ics.gmao.machines;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Application principale du service de gestion des machines
 * Service responsable de la gestion du parc machines de l'entreprise ICS
 */
@SpringBootApplication
public class MachinesServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(MachinesServiceApplication.class, args);
    }
}