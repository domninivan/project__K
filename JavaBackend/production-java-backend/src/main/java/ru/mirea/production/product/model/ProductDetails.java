package ru.mirea.production.product.model;

import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.entity.ProductMaterialRequirement;
import ru.mirea.production.product.entity.ProductOperation;

import java.util.List;

public record ProductDetails(
        Product product,
        List<ProductOperation> operations,
        List<ProductMaterialRequirement> materials
) {
}