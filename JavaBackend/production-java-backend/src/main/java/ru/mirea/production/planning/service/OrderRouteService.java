package ru.mirea.production.planning.service;

import java.time.OffsetDateTime;
import java.time.ZoneId;
import java.util.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;
import ru.mirea.production.order.dao.OrderRouteDAO;
import ru.mirea.production.order.entity.*;
import ru.mirea.production.order.entity.enums.StatusEnum;
import ru.mirea.production.order.exception.OrderNotFoundException;
import ru.mirea.production.product.service.ProductService;
import ru.mirea.production.planning.dto.request.*;
import ru.mirea.production.planning.mapper.PlanningMapper;
import ru.mirea.production.security.CurrentActor;
import ru.mirea.production.user.entity.Role;

@Service
public class OrderRouteService {
    private static final Map<String, String> TITLES = Map.of(
            "CUTTING", "Резка", "TURNING", "Точение", "MILLING", "Фрезерование",
            "GRINDING", "Шлифование", "DRILLING", "Сверление", "INSPECTION", "Контроль");
    private final OrderRouteDAO routes;
    private final ProductService products;
    private final PlanningMapper mapper;
    @Value("${planning.start-delay-minutes:60}")
    private long startDelayMinutes;

    public OrderRouteService(OrderRouteDAO routes, ProductService products, PlanningMapper mapper) {
        this.routes = routes;
        this.products = products;
        this.mapper = mapper;
    }


    @Transactional
    public PlanningRequest preparePlan(int orderId) {
        CurrentActor.require(Role.PLANNER);
        ProductionOrder order = routes.lockOrder(orderId);
        if (order == null) throw new OrderNotFoundException("Заказ не найден");
        if (order.getStatus() == StatusEnum.CANCELLED || order.getStatus() == StatusEnum.COMPLETED
                || order.getStatus() == StatusEnum.IN_PRODUCTION) {
            throw conflict("Расчёт недоступен в текущем состоянии заказа");
        }
        List<OrderItemOperation> frozen = freezeRoute(order);
        List<PlanningItemRequest> items = new ArrayList<>();
        for (ProductionOrderItem item : order.getItems()) {
            var operations = frozen.stream().filter(o -> o.getOrderItem().getId().equals(item.getId()))
                    .sorted(Comparator.comparing(OrderItemOperation::getSequenceNumber)).toList();
            if (operations.isEmpty()) throw conflict("Неполный снимок маршрута заказа");
            items.add(mapper.toPlanningItemRequest(item, operations));
        }
        if (startDelayMinutes < 0) throw new IllegalStateException("planning.start-delay-minutes must be nonnegative");
        return mapper.toPlanningRequest(order, items,
                OffsetDateTime.now(ZoneId.of("Europe/Moscow")).plusMinutes(startDelayMinutes));
    }

    private List<OrderItemOperation> freezeRoute(ProductionOrder order) {
        List<OrderItemOperation> existing = routes.find(order.getId());
        if (!existing.isEmpty()) return existing;
        if (order.getStatus() != StatusEnum.NEW && order.getStatus() != StatusEnum.READY_FOR_PLANNING) {
            throw conflict("Для старого рассчитанного заказа нельзя восстановить исходную технологию из текущего справочника");
        }
        if (order.getItems().isEmpty()) throw conflict("Заказ не содержит позиций");
        List<OrderItemOperation> snapshot = new ArrayList<>();
        for (ProductionOrderItem item : order.getItems()) {
            var source = products.getProductOperations(item.getProduct().getId());
            if (source.isEmpty()) throw conflict("У изделия нет технологического маршрута");
            Set<Integer> sequences = new HashSet<>();
            Set<Integer> ids = new HashSet<>();
            for (var operation : source) {
                if (operation.getId() == null || !ids.add(operation.getId())
                        || operation.getSequenceNumber() == null || operation.getSequenceNumber() <= 0
                        || !sequences.add(operation.getSequenceNumber())
                        || operation.getTimePerUnitSeconds() == null || operation.getTimePerUnitSeconds() <= 0
                        || operation.getOperationType() == null || !TITLES.containsKey(operation.getOperationType())) {
                    throw conflict("Некорректный технологический маршрут изделия");
                }
                snapshot.add(new OrderItemOperation(item, operation, TITLES.get(operation.getOperationType())));
            }
        }
        snapshot.forEach(routes::save);
        return snapshot;
    }

    private ResponseStatusException conflict(String message) {
        return new ResponseStatusException(HttpStatus.CONFLICT, message);
    }
}
