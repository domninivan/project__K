package ru.mirea.production.planning.dto.response;

import java.time.OffsetDateTime;

public record ScheduledOperationResponse(
        Integer scheduledOperationId,
        Integer operationId,
        Integer sequenceNumber,
        String operationType,
        Integer workCenterId,
        String workCenterCode,
        String workCenterTitle,
        OffsetDateTime startAt,
        OffsetDateTime finishAt
) {}