package ru.mirea.production.planning.client;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.ProductionPlanResponse;
import ru.mirea.production.planning.dto.response.WorkCenterResponse;

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
                        new ParameterizedTypeReference<List<WorkCenterResponse>>() {}
                );
    }
}