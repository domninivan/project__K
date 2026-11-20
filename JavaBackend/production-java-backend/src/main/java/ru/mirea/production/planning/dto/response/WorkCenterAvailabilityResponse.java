package ru.mirea.production.planning.dto.response;

import java.util.List;

public record WorkCenterAvailabilityResponse(
        Integer workCenterId,
        String timezone,
        List<AvailabilityIntervalResponse> intervals
) {}