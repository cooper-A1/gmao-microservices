package com.ics.gmao.machines.service;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;

/**
 * Client pour communiquer avec le service interventions
 */
@Service
public class InterventionServiceClient {
    
    private final WebClient webClient;
    
    @Value("${services.interventions.url:http://interventions-service:8002}")
    private String interventionsServiceUrl;
    
    public InterventionServiceClient(WebClient.Builder webClientBuilder) {
        this.webClient = webClientBuilder.build();
    }
    
    /**
     * Récupère les interventions pour une machine donnée
     */
    public Object getInterventionsByMachine(Long machineId) {
        try {
            return webClient.get()
                    .uri(interventionsServiceUrl + "/api/interventions/machine/" + machineId)
                    .retrieve()
                    .bodyToMono(Object.class)
                    .block();
        } catch (Exception e) {
            // Log l'erreur et retourner une réponse par défaut
            System.err.println("Erreur lors de l'appel au service interventions: " + e.getMessage());
            return "Service interventions indisponible";
        }
    }
}