package ru.mirea.production.product.dto;

public record ProductResponseDTO(
        Integer id,
        String title,
        String description
) {}