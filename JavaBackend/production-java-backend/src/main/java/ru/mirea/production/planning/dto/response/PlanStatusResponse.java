package ru.mirea.production.planning.dto.response;

import ru.mirea.production.planning.dto.enums.PlanStatus;

import java.time.OffsetDateTime;

public record PlanStatusResponse(
        Integer planId,
        Integer orderId,
        PlanStatus status,
        OffsetDateTime updatedAt
) {
}