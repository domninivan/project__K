package ru.mirea.production.product.dao;

import ru.mirea.production.product.entity.ProductMaterialRequirement;

import java.util.List;

public interface ProductMaterialRequirementDAO {
    List<ProductMaterialRequirement> getProductMaterialRequirementsByProductId(int productId);
}
