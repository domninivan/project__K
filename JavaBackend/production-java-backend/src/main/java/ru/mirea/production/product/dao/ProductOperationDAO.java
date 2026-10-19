package ru.mirea.production.product.dao;

import ru.mirea.production.product.entity.ProductOperation;
import java.util.List;

public interface ProductOperationDAO {
    List<ProductOperation> getProductOperationByProductId(int productId);
}
