package ru.mirea.production.planning.dto.response;

import java.time.OffsetDateTime;

public record ScheduleOperationResponse(
        Integer scheduledOperationId,
        Integer orderId,
        Integer orderItemId,
        Integer productId,
        String operationType,
        Integer sequenceNumber,
        OffsetDateTime startAt,
        OffsetDateTime finishAt
) {
}