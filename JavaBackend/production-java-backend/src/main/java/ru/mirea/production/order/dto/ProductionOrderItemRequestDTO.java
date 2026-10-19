package ru.mirea.production.order.dto;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

public record ProductionOrderItemRequestDTO(

        @NotNull
        @Positive
        Integer productId,

        @NotNull
        @Positive
        Integer quantity
) {}
