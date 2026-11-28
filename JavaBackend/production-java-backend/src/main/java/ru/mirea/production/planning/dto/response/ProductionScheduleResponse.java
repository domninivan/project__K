package ru.mirea.production.planning.dto.response;

import java.time.OffsetDateTime;
import java.util.List;

public record ProductionScheduleResponse(
        OffsetDateTime from,
        OffsetDateTime to,
        List<WorkCenterProductionScheduleResponse> workCenters
) {
}