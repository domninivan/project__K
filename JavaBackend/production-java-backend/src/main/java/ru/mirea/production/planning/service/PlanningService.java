package ru.mirea.production.planning.service;

import ru.mirea.production.planning.dto.response.*;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;

public interface PlanningService {

    ProductionPlanResponse planOrder(int orderId);

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