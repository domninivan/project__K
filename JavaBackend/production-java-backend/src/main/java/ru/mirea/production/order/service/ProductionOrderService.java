package ru.mirea.production.order.service;

import ru.mirea.production.order.dto.ProductionOrderRequestDTO;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.ProductionOrderItem;
import ru.mirea.production.order.entity.enums.PriorityEnum;

import java.time.LocalDate;
import java.util.List;

public interface ProductionOrderService {
    ProductionOrder createOrder(ProductionOrderRequestDTO orderDTO);
    ProductionOrder getOrderById(int id);
    List<ProductionOrder> getAllOrders();
    List<ProductionOrderItem> getOrderItemsByOrderId(int order_id);
    ProductionOrder cancelOrder(int id);
}
