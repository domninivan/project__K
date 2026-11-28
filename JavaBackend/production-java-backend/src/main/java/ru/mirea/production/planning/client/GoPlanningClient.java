package ru.mirea.production.planning.client;

import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.*;


import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;

public interface GoPlanningClient {

    ProductionPlanResponse createPlan(PlanningRequest request);

    ProductionPlanResponse getOrderPlan(int orderId);

    List<WorkCenterResponse> getWorkCenters();

    WorkCenterAvailabilityResponse getWorkCenterAvailability(
            int workCenterId,
            LocalDate from,
            LocalDate to
    );

    ProductionScheduleResponse getSchedule(
            OffsetDateTime from,
            OffsetDateTime to
    );

    PlanStatusResponse acceptPlan(int planId);

    PlanStatusResponse rejectPlan(int planId);

    PlanStatusResponse cancelOrderPlan(int orderId);


}