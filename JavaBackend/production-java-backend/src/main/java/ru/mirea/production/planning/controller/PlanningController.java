package ru.mirea.production.planning.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.web.bind.annotation.*;
import ru.mirea.production.planning.dto.response.*;
import ru.mirea.production.planning.service.PlanningService;

import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;

@RestController
@RequestMapping("/api/planner")
public class PlanningController {

    @Autowired
    private PlanningService planningService;

    @PostMapping("/orders/{orderId}")
    public ProductionPlanResponse planOrder(
            @PathVariable int orderId
    ) {
        return planningService.planOrder(orderId);
    }

    @GetMapping("/orders/{orderId}")
    public ProductionPlanResponse getOrderPlan(
            @PathVariable int orderId
    ) {
        return planningService.getOrderPlan(orderId);
    }

    @GetMapping("/work-centers")
    public List<WorkCenterResponse> getWorkCenters() {
        return planningService.getWorkCenters();
    }

    @GetMapping("/work-centers/{id}/availability")
    public WorkCenterAvailabilityResponse getWorkCenterAvailability(
            @PathVariable("id") int workCenterId,
            @DateTimeFormat(iso = DateTimeFormat.ISO.DATE)
            @RequestParam
            LocalDate from,
            @DateTimeFormat(iso = DateTimeFormat.ISO.DATE)
            @RequestParam LocalDate to
    ) {
        return planningService.getWorkCenterAvailability(
                workCenterId,
                from,
                to
        );
    }

    @GetMapping("/schedule")
    public ProductionScheduleResponse getSchedule(
            @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)
            @RequestParam OffsetDateTime from,
            @DateTimeFormat(iso = DateTimeFormat.ISO.DATE_TIME)
            @RequestParam OffsetDateTime to
    ) {
        return planningService.getSchedule(from, to);
    }

    @PostMapping("/plans/{planId}/accept")
    public PlanStatusResponse acceptPlan(
            @PathVariable("planId") int planId
    ) {
        return planningService.acceptPlan(planId);
    }

    @PostMapping("/plans/{planId}/reject")
    public PlanStatusResponse rejectPlan(
            @PathVariable("planId") int planId
    ) {
        return planningService.rejectPlan(planId);
    }

    @PostMapping("/orders/{orderId}/cancel-plan")
    public PlanStatusResponse cancelOrderPlan(
            @PathVariable("orderId") int orderId
    ) {
        return planningService.cancelOrderPlan(orderId);
    }


}