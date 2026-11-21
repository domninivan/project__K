package ru.mirea.production.planning.dto.request;

import java.util.List;

public record PlanningItemRequest(
        Integer orderItemId,
        Integer productId,
        Integer quantity,
        List<PlanningOperationRequest> operations
) {}
