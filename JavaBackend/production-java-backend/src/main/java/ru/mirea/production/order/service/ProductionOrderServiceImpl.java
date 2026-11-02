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

import java.util.*;

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

        Map<Integer, Integer> productQuantities = new LinkedHashMap<>();
        for(var item: orderDTO.items()){
            productQuantities.merge(item.productId(), item.quantity(), Integer::sum);
        }

        for(var entry: productQuantities.entrySet()){
            Product product = productService.getProduct(entry.getKey());
            order.addItem(new ProductionOrderItem(entry.getValue(), product));
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

        if (productionOrder.getStatus() == StatusEnum.CANCELLED) {
            return productionOrder;
        }

        productionOrder.setStatus(StatusEnum.CANCELLED);
        return productionOrderDAO.saveOrder(productionOrder);
    }

}
