package ru.mirea.production.planning.client;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.*;


import java.time.LocalDate;
import java.time.OffsetDateTime;
import java.util.List;

@Component
public class GoPlanningClientImpl implements GoPlanningClient {

    @Autowired
    @Qualifier("goRestClient")
    private RestClient restClient;

    @Override
    public ProductionPlanResponse createPlan(PlanningRequest request) {
        return restClient.post()
                .uri("/api/plans")
                .body(request)
                .retrieve()
                .body(ProductionPlanResponse.class);
    }

    @Override
    public ProductionPlanResponse getOrderPlan(int orderId) {
        return restClient.get()
                .uri("/api/plans/orders/{orderId}", orderId)
                .retrieve()
                .body(ProductionPlanResponse.class);
    }

    @Override
    public List<WorkCenterResponse> getWorkCenters() {
        return restClient.get()
                .uri("/api/work-centers")
                .retrieve()
                .body(
                        new ParameterizedTypeReference<List<WorkCenterResponse>>() {
                        }
                );
    }

    @Override
    public WorkCenterAvailabilityResponse getWorkCenterAvailability(
            int workCenterId,
            LocalDate from,
            LocalDate to
    ) {
        return restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/work-centers/{id}/availability")
                        .queryParam("from", from)
                        .queryParam("to", to)
                        .build(workCenterId))
                .retrieve()
                .body(WorkCenterAvailabilityResponse.class);
    }

    @Override
    public ProductionScheduleResponse getSchedule(
            OffsetDateTime from,
            OffsetDateTime to
    ) {
        return restClient.get()
                .uri(uriBuilder -> uriBuilder
                        .path("/api/schedule")
                        .queryParam("from", from)
                        .queryParam("to", to)
                        .build())
                .retrieve()
                .body(ProductionScheduleResponse.class);
    }

    @Override
    public PlanStatusResponse acceptPlan(int planId) {
        return restClient.post()
                .uri("/api/plans/{planId}/accept", planId)
                .retrieve()
                .body(PlanStatusResponse.class);
    }

    @Override
    public PlanStatusResponse rejectPlan(int planId) {
        return restClient.post()
                .uri("/api/plans/{planId}/reject", planId)
                .retrieve()
                .body(PlanStatusResponse.class);
    }

    @Override
    public PlanStatusResponse cancelOrderPlan(int orderId) {
        return restClient.post()
                .uri("/api/plans/orders/{orderId}/cancel", orderId)
                .retrieve()
                .body(PlanStatusResponse.class);
    }

}

