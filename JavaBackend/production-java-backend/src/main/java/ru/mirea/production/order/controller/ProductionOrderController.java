package ru.mirea.production.order.controller;

import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import ru.mirea.production.order.dto.ProductionOrderRequestDTO;
import ru.mirea.production.order.dto.ProductionOrderResponseDTO;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.mapper.ProductionOrderMapper;
import ru.mirea.production.order.service.ProductionOrderService;

import java.util.List;

@RestController
@RequestMapping("/api/orders")
public class ProductionOrderController {

    @Autowired
    private ProductionOrderService productionOrderService;

    @Autowired
    private ProductionOrderMapper productionOrderMapper;

    @PostMapping
    public ProductionOrderResponseDTO createOrder(@Valid @RequestBody ProductionOrderRequestDTO orderDTO){
        ProductionOrder order = productionOrderService.createOrder(orderDTO);
        return productionOrderMapper.toResponse(order);
    }

    @GetMapping
    public List<ProductionOrderResponseDTO> getAllOrders(){
        List<ProductionOrder> orders = productionOrderService.getAllOrders();
        return productionOrderMapper.toOrdersResponses(orders);
    }

    @GetMapping("/{id}")
    public ProductionOrderResponseDTO getOrderById(@PathVariable("id") int id){
        ProductionOrder order = productionOrderService.getOrderById(id);
        return productionOrderMapper.toResponse(order);
    }

    @DeleteMapping("/{id}")
    public ProductionOrderResponseDTO cancelOrderById(@PathVariable("id") int id){
        ProductionOrder order = productionOrderService.cancelOrder(id);
        return productionOrderMapper.toResponse(order);
    }

    @GetMapping("/hello")
    public String hello(){
        return "Hello from Case 22";
    }
}
