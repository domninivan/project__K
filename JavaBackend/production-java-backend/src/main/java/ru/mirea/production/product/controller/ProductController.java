package ru.mirea.production.product.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import ru.mirea.production.product.dto.ProductDetailResponseDTO;
import ru.mirea.production.product.dto.ProductResponseDTO;
import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.model.ProductDetails;
import ru.mirea.production.product.mapper.ProductMapper;
import ru.mirea.production.product.service.ProductService;

import java.util.List;

@RestController
@RequestMapping("/api/products")
public class ProductController {

    @Autowired
    private ProductService productService;

    @Autowired
    private ProductMapper mapper;

    @GetMapping
    public List<ProductResponseDTO> getAllProducts(){
        List<Product> products = productService.getAllProducts();
        return mapper.toResponses(products);
    }

    @GetMapping("/{id}")
    public ProductDetailResponseDTO getProductById(@PathVariable("id") int id){
        ProductDetails product = productService.getProductDetails(id);
        return mapper.toDetailResponse(product);
    }
}
