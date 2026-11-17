package ru.mirea.production.planning.client;

import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.ProductionPlanResponse;
import ru.mirea.production.planning.dto.response.WorkCenterResponse;

import java.util.List;

public interface GoPlanningClient {

    ProductionPlanResponse createPlan(PlanningRequest request);

    ProductionPlanResponse getOrderPlan(int orderId);

    List<WorkCenterResponse> getWorkCenters();
}