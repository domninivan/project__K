package ru.mirea.production.planning.dto.request;

public record PlanningOperationRequest(
        Integer operationId,
        Integer sequenceNumber,
        String operationType,
        Integer timePerUnitSeconds
) {}
