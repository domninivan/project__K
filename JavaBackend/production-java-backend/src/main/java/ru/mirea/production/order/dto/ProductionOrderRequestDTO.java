package ru.mirea.production.order.dto;

import jakarta.validation.Valid;
import jakarta.validation.constraints.*;
import java.time.LocalDate;
import java.util.List;

public record ProductionOrderRequestDTO(

        @NotNull
        LocalDate dueDate,

        @Size(max=200)
        String comment,

        @NotNull
        @NotEmpty
        @Valid
        List<ProductionOrderItemRequestDTO> items

) {}
