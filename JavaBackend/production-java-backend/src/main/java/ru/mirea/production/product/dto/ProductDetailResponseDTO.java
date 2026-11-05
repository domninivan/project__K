package ru.mirea.production.product.dto;

import java.util.List;

public record ProductDetailResponseDTO(
        Integer id,
        String title,
        String description,
        List<ProductOperationResponseDTO> operations,
        List<ProductMaterialRequirementResponseDTO> materials
) {
}