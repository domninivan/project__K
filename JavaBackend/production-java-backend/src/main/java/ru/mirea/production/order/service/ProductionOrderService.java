package ru.mirea.production.order.service;

import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.enums.PriorityEnum;

import java.time.LocalDate;
import java.util.List;

public interface ProductionOrderService {
    ProductionOrder createOrder(int productId, int quantity, LocalDate dueDate, PriorityEnum priority, String comment);
    ProductionOrder getOrderById(int id);
    List<ProductionOrder> getAllOrders();
    ProductionOrder cancelOrder(int id);
}
