package ru.mirea.production.planning.dto.response;

import java.util.List;

public record WorkCenterProductionScheduleResponse(
        Integer workCenterId,
        String workCenterCode,
        String workCenterTitle,
        List<ScheduleOperationResponse> operations
) {
}