package ru.mirea.production.planning.service;

import ru.mirea.production.planning.dto.response.ProductionPlanResponse;
import ru.mirea.production.planning.dto.response.WorkCenterResponse;

import java.util.List;

public interface PlanningService {

    ProductionPlanResponse planOrder(int orderId);

    ProductionPlanResponse getOrderPlan(int orderId);

    List<WorkCenterResponse> getWorkCenters();
}