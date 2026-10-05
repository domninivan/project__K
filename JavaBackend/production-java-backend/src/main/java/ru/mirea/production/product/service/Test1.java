package ru.mirea.production.product.service;

import org.springframework.context.support.ClassPathXmlApplicationContext;
import ru.mirea.production.order.service.ProductionOrderService;

import java.time.LocalDate;

public class Test1 {

    public static void main(String[] args) {
        ClassPathXmlApplicationContext context =
                new ClassPathXmlApplicationContext("applicationContext.xml");
        ProductionOrderService service = context.getBean("productionOrderServiceImpl", ProductionOrderService.class);

        System.out.println(service.getOrderById(1));




    }
}
