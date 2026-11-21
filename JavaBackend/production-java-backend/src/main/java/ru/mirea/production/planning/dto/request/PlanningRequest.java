package ru.mirea.production.planning.dto.request;

import ru.mirea.production.order.entity.enums.PriorityEnum;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;

public record PlanningRequest(
        Integer orderId,
        OffsetDateTime earliestStartAt,
        LocalDate dueDate,
        PriorityEnum priority,
        List<PlanningItemRequest> items
) {}