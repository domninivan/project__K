package ru.mirea.production.planning.dto.response;

import java.time.OffsetDateTime;

public record AvailabilityIntervalResponse(
        OffsetDateTime startAt,
        OffsetDateTime endAt
) {}