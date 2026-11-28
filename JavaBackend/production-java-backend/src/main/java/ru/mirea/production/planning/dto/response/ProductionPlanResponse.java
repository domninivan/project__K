package ru.mirea.production.planning.dto.response;

import ru.mirea.production.planning.dto.enums.PlanStatus;
import java.time.OffsetDateTime;
import java.time.LocalDate;
import java.util.List;

public record ProductionPlanResponse(
        Integer planId,
        Integer orderId,
        OffsetDateTime plannedStartAt,
        OffsetDateTime plannedFinishAt,
        LocalDate dueDate,
        Boolean onTime,
        PlanStatus status,
        List<ProductionPlanItemResponse> items
) {}
