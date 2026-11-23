package ru.mirea.production.planning.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.HttpServerErrorException;
import org.springframework.web.client.ResourceAccessException;

import ru.mirea.production.planning.client.GoPlanningClient;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.planning.dto.response.*;
import ru.mirea.production.planning.exception.*;

import java.time.LocalDate;
import java.time.OffsetDateTime;

import java.util.List;
import java.util.function.Supplier;

@Service
public class PlanningServiceImpl implements PlanningService {

    @Autowired
    private GoPlanningClient goPlanningClient;

    @Autowired
    private OrderRouteService orderRouteService;

    private <T> T executeGoRequest(Supplier<T> request) {
        try {
            return request.get();

        } catch (HttpClientErrorException.BadRequest ex) {
            throw new GoBadRequestException(
                    "Сервис планирования отклонил переданные данные"
            );

        } catch (HttpClientErrorException.NotFound ex) {
            throw new GoResourceNotFoundException(
                    "Запрашиваемые данные в сервисе планирования не найдены"
            );

        } catch (HttpClientErrorException.Conflict ex) {
            throw new GoConflictException(
                    "Операция не может быть выполнена в текущем состоянии"
            );

        } catch (ResourceAccessException ex) {
            throw new GoServiceUnavailableException();

        } catch (HttpServerErrorException ex) {
            throw new GoServiceException(
                    "Внутренняя ошибка сервиса планирования"
            );

        } catch (HttpClientErrorException ex) {
            throw new GoServiceException(
                    "Сервис планирования вернул неожиданную ошибку"
            );
        }
    }

    @Override
    public ProductionPlanResponse planOrder(int orderId) {
        PlanningRequest request = orderRouteService.preparePlan(orderId);
        return executeGoRequest(() -> goPlanningClient.createPlan(request));
    }

    @Override
    public ProductionPlanResponse getOrderPlan(int orderId) {
        try {
            return goPlanningClient.getOrderPlan(orderId);

        } catch (HttpClientErrorException.NotFound ex) {
            throw new ProductionPlanNotFoundException(orderId);

        } catch (HttpClientErrorException.BadRequest ex) {
            throw new GoBadRequestException(
                    "Некорректный идентификатор заказа"
            );

        } catch (HttpClientErrorException.Conflict ex) {
            throw new GoConflictException(
                    "Невозможно получить производственный план"
            );

        } catch (ResourceAccessException ex) {
            throw new GoServiceUnavailableException();

        } catch (HttpServerErrorException ex) {
            throw new GoServiceException(
                    "Внутренняя ошибка сервиса планирования"
            );

        } catch (HttpClientErrorException ex) {
            throw new GoServiceException(
                    "Сервис планирования вернул неожиданную ошибку"
            );
        }
    }

    @Override
    public List<WorkCenterResponse> getWorkCenters() {
        return executeGoRequest(
                goPlanningClient::getWorkCenters
        );
    }

    @Override
    public WorkCenterAvailabilityResponse getWorkCenterAvailability(
            int workCenterId,
            LocalDate from,
            LocalDate to
    ) {
        return executeGoRequest(
                () -> goPlanningClient.getWorkCenterAvailability(
                        workCenterId,
                        from,
                        to
                )
        );
    }

    @Override
    public ProductionScheduleResponse getSchedule(
            OffsetDateTime from,
            OffsetDateTime to
    ) {
        return executeGoRequest(
                () -> goPlanningClient.getSchedule(
                        from,
                        to
                )
        );
    }

    @Override
    public PlanStatusResponse acceptPlan(int planId) {
        return executeGoRequest(
                () -> goPlanningClient.acceptPlan(planId)
        );
    }

    @Override
    public PlanStatusResponse rejectPlan(int planId) {
        return executeGoRequest(
                () -> goPlanningClient.rejectPlan(planId)
        );
    }

    @Override
    public PlanStatusResponse cancelOrderPlan(int orderId) {
        return executeGoRequest(
                () -> goPlanningClient.cancelOrderPlan(orderId)
        );
    }
}