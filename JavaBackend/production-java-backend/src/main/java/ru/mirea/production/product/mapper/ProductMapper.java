package ru.mirea.production.product.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import ru.mirea.production.product.dto.*;
import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.model.ProductDetails;
import ru.mirea.production.product.entity.ProductMaterialRequirement;
import ru.mirea.production.product.entity.ProductOperation;

import java.util.List;

@Mapper(componentModel = "spring")
public interface ProductMapper {

    // Для GET /api/products
    ProductResponseDTO toResponse(Product product);

    List<ProductResponseDTO> toResponses(List<Product> products);


    // Для GET /api/products/{id}
    @Mapping(source = "product.id", target = "id")
    @Mapping(source = "product.title", target = "title")
    @Mapping(source = "product.description", target = "description")
    @Mapping(source = "operations", target = "operations")
    @Mapping(source = "materials", target = "materials")
    ProductDetailResponseDTO toDetailResponse(ProductDetails details);


    // Вложенные операции
    ProductOperationResponseDTO toOperationResponse(
            ProductOperation operation
    );

    List<ProductOperationResponseDTO> toOperationResponses(
            List<ProductOperation> operations
    );


    // Вложенные материалы
    ProductMaterialRequirementResponseDTO toMaterialResponse(
            ProductMaterialRequirement material
    );

    List<ProductMaterialRequirementResponseDTO> toMaterialResponses(
            List<ProductMaterialRequirement> materials
    );
}