package ru.mirea.production.planning.exception;

public class ProductionPlanNotFoundException extends RuntimeException {

    public ProductionPlanNotFoundException(int orderId) {
        super("Производственный план для заказа с id " + orderId + " не найден");
    }
}


