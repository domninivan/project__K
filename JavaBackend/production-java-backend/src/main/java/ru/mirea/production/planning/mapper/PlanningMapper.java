package ru.mirea.production.planning.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import ru.mirea.production.order.entity.ProductionOrder;
import ru.mirea.production.order.entity.ProductionOrderItem;
import ru.mirea.production.planning.dto.request.PlanningItemRequest;
import ru.mirea.production.planning.dto.request.PlanningOperationRequest;
import ru.mirea.production.planning.dto.request.PlanningRequest;
import ru.mirea.production.order.entity.OrderItemOperation;

import java.time.OffsetDateTime;
import java.util.List;

@Mapper(componentModel = "spring")
public interface PlanningMapper {

    @Mapping(source = "item.id", target = "orderItemId")
    @Mapping(source = "item.product.id", target = "productId")
    @Mapping(source = "item.quantity", target = "quantity")
    @Mapping(source = "operations", target = "operations")
    PlanningItemRequest toPlanningItemRequest(
            ProductionOrderItem item,
            List<OrderItemOperation> operations
    );

    @Mapping(source = "sourceOperationId", target = "operationId")
    PlanningOperationRequest toPlanningOperationRequest(
            OrderItemOperation operation
    );

    @Mapping(source = "order.id", target = "orderId")
    @Mapping(source = "earliestStartAt", target = "earliestStartAt")
    @Mapping(source = "order.dueDate", target = "dueDate")
    @Mapping(source = "order.priority", target = "priority")
    @Mapping(source = "items", target = "items")
    PlanningRequest toPlanningRequest(
            ProductionOrder order,
            List<PlanningItemRequest> items,
            OffsetDateTime earliestStartAt
    );
}
