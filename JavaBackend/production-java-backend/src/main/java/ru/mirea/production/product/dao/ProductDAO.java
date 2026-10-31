package ru.mirea.production.product.dao;

import ru.mirea.production.product.entity.Product;
import java.util.List;

public interface ProductDAO {
    List<Product> getAllProducts();
    Product getProduct(int id);
}
