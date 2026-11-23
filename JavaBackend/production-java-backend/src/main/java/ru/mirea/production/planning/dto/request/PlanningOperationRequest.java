package ru.mirea.production.planning.dto.request;

public record PlanningOperationRequest(
        Long operationId,
        Integer sequenceNumber,
        String operationType,
        Long timePerUnitSeconds
) {}
