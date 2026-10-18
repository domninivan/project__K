package ru.mirea.production.order.dto;

public record ProductionOrderItemResponseDTO(
        Integer id,
        Integer productId,
        String productTitle,
        Integer quantity
) {}