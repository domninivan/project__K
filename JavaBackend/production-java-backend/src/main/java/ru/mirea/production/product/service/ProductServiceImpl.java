package ru.mirea.production.product.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.mirea.production.product.dao.ProductDAO;
import ru.mirea.production.product.dao.ProductMaterialRequirementDAO;
import ru.mirea.production.product.dao.ProductOperationDAO;
import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.entity.ProductMaterialRequirement;
import ru.mirea.production.product.entity.ProductOperation;
import ru.mirea.production.product.exception.ProductNotFoundException;

import java.util.List;

@Service
public class ProductServiceImpl implements ProductService{

    @Autowired
    private ProductDAO productDAO;

    @Autowired
    private ProductOperationDAO productOperationDAO;

    @Autowired
    private ProductMaterialRequirementDAO productMaterialRequirementDAO;

    private void requireProduct(Product product, int productId){
        if (product == null){
            throw new ProductNotFoundException("Деталь с ID = " + productId + " не найдена!");
        }
    }

    @Override
    @Transactional(readOnly = true)
    public List<Product> getAllProducts() {
        return productDAO.getAllProducts();
    }

    @Override
    @Transactional(readOnly = true)
    public Product getProduct(int id) {
        Product product = productDAO.getProduct(id);
        requireProduct(product, id);
        return product;
    }

    @Override
    @Transactional(readOnly = true)
    public List<ProductOperation> getProductOperationsByProductId(int productId) {
        Product product = productDAO.getProduct(productId);
        requireProduct(product, productId);
        return productOperationDAO.getProductOperationByProductId(productId);
    }

    @Override
    @Transactional(readOnly = true)
    public List<ProductMaterialRequirement> getProductMaterialRequirementsByProductId(int productId) {
        Product product = productDAO.getProduct(productId);
        requireProduct(product, productId);
        return productMaterialRequirementDAO.getProductMaterialRequirementsByProductId(productId);
    }
}
