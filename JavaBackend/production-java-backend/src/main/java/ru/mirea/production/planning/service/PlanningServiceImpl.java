package ru.mirea.production.planning.service;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.ProductionOrderItem;
import ru.mirea.production.order.service.ProductionOrderService;
import ru.mirea.production.planning.client.GoPlanningClient;
import ru.mirea.production.planning.dto.request.PlanningItemRequest;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.ProductionPlanResponse;
import ru.mirea.production.planning.dto.response.WorkCenterResponse;
import ru.mirea.production.planning.exception.GoServiceException;
import ru.mirea.production.planning.exception.GoServiceUnavailableException;
import ru.mirea.production.planning.exception.ProductionPlanNotFoundException;
import ru.mirea.production.planning.mapper.PlanningMapper;
import ru.mirea.production.product.entity.ProductOperation;
import ru.mirea.production.product.service.ProductService;

import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;

@Service
public class PlanningServiceImpl implements PlanningService {

    @Autowired
    private GoPlanningClient goPlanningClient;

    @Autowired
    private ProductionOrderService productionOrderService;

    @Autowired
    private ProductService productService;

    @Autowired
    private PlanningMapper planningMapper;

    private <T> T executeGoRequest(Supplier<T> request) {
        try {
            return request.get();

        } catch (ResourceAccessException ex) {
            throw new GoServiceUnavailableException();

        } catch (HttpServerErrorException ex) {
            throw new GoServiceException(
                    "Ошибка сервиса планирования"
            );

        } catch (HttpClientErrorException ex) {
            throw new GoServiceException(
                    "Сервис планирования отклонил запрос"
            );
        }
    }


    @Override
    public ProductionPlanResponse planOrder(int orderId) {

        ProductionOrder order =
                productionOrderService.getOrderById(orderId);

        List<PlanningItemRequest> planningItems = new ArrayList<>();

        for (ProductionOrderItem item : order.getItems()) {

            List<ProductOperation> operations =
                    productService.getProductOperations(
                            item.getProduct().getId()
                    );

            PlanningItemRequest planningItem =
                    planningMapper.toPlanningItemRequest(
                            item,
                            operations
                    );

            planningItems.add(planningItem);
        }

        PlanningRequest request =
                planningMapper.toPlanningRequest(
                        order,
                        planningItems,
                        OffsetDateTime.now(ZoneId.of("Europe/Moscow"))
                );

        return executeGoRequest(
                () -> goPlanningClient.createPlan(request)
        );
    }

    @Override
    public ProductionPlanResponse getOrderPlan(int orderId) {
        try {
            return goPlanningClient.getOrderPlan(orderId);

        } catch (HttpClientErrorException.NotFound ex) {
            throw new ProductionPlanNotFoundException(orderId);

        } catch (ResourceAccessException ex) {
            throw new GoServiceUnavailableException();

        } catch (HttpServerErrorException ex) {
            throw new GoServiceException(
                    "Ошибка сервиса планирования"
            );

        } catch (HttpClientErrorException ex) {
            throw new GoServiceException(
                    "Сервис планирования отклонил запрос"
            );
        }
    }

    @Override
    public List<WorkCenterResponse> getWorkCenters() {
        return executeGoRequest(
                goPlanningClient::getWorkCenters
        );
    }
}