package ru.mirea.production.order.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import ru.mirea.production.order.dao.ProductionOrderDAO;
import ru.mirea.production.order.entity.enums.PriorityEnum;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.enums.StatusEnum;
import ru.mirea.production.order.exception.OrderNotFoundException;
import ru.mirea.production.product.entity.Product;
import ru.mirea.production.product.service.ProductService;

import java.time.LocalDate;
import java.util.List;

@Service
public class ProductionOrderServiceImpl implements ProductionOrderService{

    @Autowired
    private ProductionOrderDAO productionOrderDAO;

    @Autowired
    private ProductService productService;


    @Override
    @Transactional
    public ProductionOrder createOrder(int productId, int quantity, LocalDate dueDate, PriorityEnum priority, String comment){
        Product product = productService.getProduct(productId);
        ProductionOrder productionOrder = new ProductionOrder(quantity,
                                                              dueDate,
                                                              priority,
                                                              comment,
                                                              product);
        return productionOrderDAO.saveOrder(productionOrder);
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
}
