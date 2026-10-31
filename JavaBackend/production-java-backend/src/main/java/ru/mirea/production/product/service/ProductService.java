package ru.mirea.production.product.service;

import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.entity.ProductMaterialRequirement;
import ru.mirea.production.product.entity.ProductOperation;

import java.util.List;

public interface ProductService {
    List<Product> getAllProducts();
    Product getProduct(int id);
    List<ProductOperation> getProductOperationsByProductId(int productId);
    List<ProductMaterialRequirement> getProductMaterialRequirementsByProductId(int productId);
}
