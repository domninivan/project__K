package ru.mirea.production.order.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import ru.mirea.production.order.dto.ProductionOrderItemResponseDTO;
import ru.mirea.production.order.dto.ProductionOrderResponseDTO;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.ProductionOrderItem;

import java.util.List;

@Mapper(componentModel = "spring")
public interface ProductionOrderMapper {
    ProductionOrderResponseDTO toResponse(ProductionOrder order);


    @Mapping(source = "product.title", target = "productTitle")
    @Mapping(source = "product.id",   target = "productId")
    ProductionOrderItemResponseDTO toItemResponse(ProductionOrderItem item);

    List<ProductionOrderItemResponseDTO> toItemResponses(List<ProductionOrderItem> items);

    List<ProductionOrderResponseDTO> toOrdersResponses(List<ProductionOrder> orders);
}
