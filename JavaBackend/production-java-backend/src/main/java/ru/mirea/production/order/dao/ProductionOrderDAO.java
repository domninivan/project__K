package ru.mirea.production.order.dao;

import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.ProductionOrderItem;

import java.util.List;

public interface ProductionOrderDAO {
    List<ProductionOrder> getAllOrders();
    ProductionOrder getOrderById(int id);
    ProductionOrder saveOrder(ProductionOrder order);

}
