package ru.mirea.production.product.dto;

import ru.mirea.production.product.entity.enums.UnitOfMeasure;
import java.math.BigDecimal;

public record ProductMaterialRequirementResponseDTO(
        Integer id,
        String materialCode,
        BigDecimal quantityPerUnit,
        UnitOfMeasure unitCode
) {
}