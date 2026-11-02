package ru.mirea.production.product.dto;

public record ProductOperationResponseDTO(
        Integer id,
        Integer sequenceNumber,
        String operationType,
        Integer timePerUnitSeconds
) {
}