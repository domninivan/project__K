package ru.mirea.production.order.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.mirea.production.order.dao.ProductionOrderDAO;
import ru.mirea.production.order.dto.ProductionOrderRequestDTO;
import ru.mirea.production.order.entity.ProductionOrderItem;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.enums.StatusEnum;
import ru.mirea.production.order.exception.OrderNotFoundException;
import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.service.ProductService;
import java.util.List;

@Service
public class ProductionOrderServiceImpl implements ProductionOrderService{

    @Autowired
    private ProductionOrderDAO productionOrderDAO;

    @Autowired
    private ProductService productService;

    @Override
    @Transactional
    public ProductionOrder createOrder(ProductionOrderRequestDTO orderDTO){
        ProductionOrder order =
                new ProductionOrder(orderDTO.dueDate(), orderDTO.comment());

        for(var item: orderDTO.items()){
            Product product = productService.getProduct(item.productId());
            order.addItem(new ProductionOrderItem(item.quantity(), product));
        }

        return productionOrderDAO.saveOrder(order);
    }

    @Override
    @Transactional(readOnly = true)
    public ProductionOrder getOrderById(int id) {
        ProductionOrder productionOrder = productionOrderDAO.getOrderById(id);
        if (productionOrder == null) {
            throw new OrderNotFoundException(
                    "Заказ с ID = " + id + " не найден"
            );
        }
        return productionOrder;
    }

    @Override
    @Transactional(readOnly = true)
    public List<ProductionOrder> getAllOrders() {
        return productionOrderDAO.getAllOrders();
    }

    @Override
    @Transactional
    public ProductionOrder cancelOrder(int id) {
        ProductionOrder productionOrder = getOrderById(id);
        productionOrder.setStatus(StatusEnum.CANCELLED);
        return productionOrderDAO.saveOrder(productionOrder);
    }

    @Transactional
    @Override
    public List<ProductionOrderItem> getOrderItemsByOrderId(int order_id) {
        ProductionOrder productionOrder = getOrderById(order_id);
        return productionOrder.getItems();
    }
}
