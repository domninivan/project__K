package ru.mirea.production.planning.dto.response;

import java.time.OffsetDateTime;
import java.util.List;

public record ProductionPlanItemResponse(
        Integer orderItemId,
        Integer productId,
        Integer quantity,
        OffsetDateTime plannedStartAt,
        OffsetDateTime plannedFinishAt,
        List<ScheduledOperationResponse> operations
) {}