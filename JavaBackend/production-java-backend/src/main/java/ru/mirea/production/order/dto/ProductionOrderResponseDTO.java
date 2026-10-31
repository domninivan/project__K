package ru.mirea.production.order.dto;

import ru.mirea.production.order.entity.enums.PriorityEnum;
import ru.mirea.production.order.entity.enums.StatusEnum;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.List;

public record ProductionOrderResponseDTO(
        Integer id,
        LocalDate dueDate,
        PriorityEnum priority,
        StatusEnum status,
        LocalDateTime createdAt,
        String comment,
        List<ProductionOrderItemResponseDTO> items
) {}