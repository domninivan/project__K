package ru.mirea.production.planning.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import ru.mirea.production.planning.dto.response.ProductionPlanResponse;
import ru.mirea.production.planning.dto.response.WorkCenterResponse;
import ru.mirea.production.planning.service.PlanningService;

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
}