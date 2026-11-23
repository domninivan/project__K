package ru.mirea.production.order.entity;

import jakarta.persistence.*;
import org.hibernate.annotations.Immutable;
import ru.mirea.production.product.entity.ProductOperation;

@Entity
@Immutable
@Table(name = "order_item_operation")
public class OrderItemOperation {
    @Id @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "order_item_id", nullable = false)
    private ProductionOrderItem orderItem;
    @Column(name = "source_operation_id", nullable = false)
    private Long sourceOperationId;
    @Column(name = "sequence_number", nullable = false)
    private Integer sequenceNumber;
    @Column(name = "operation_type", nullable = false, length = 20)
    private String operationType;
    @Column(nullable = false, length = 200)
    private String title;
    @Column(name = "time_per_unit_seconds", nullable = false)
    private Long timePerUnitSeconds;

    protected OrderItemOperation() {}
    public OrderItemOperation(ProductionOrderItem item, ProductOperation source, String title) {
        this.orderItem = item;
        this.sourceOperationId = source.getId().longValue();
        this.sequenceNumber = source.getSequenceNumber();
        this.operationType = source.getOperationType();
        this.title = title;
        this.timePerUnitSeconds = source.getTimePerUnitSeconds().longValue();
    }
    public Long getId() { return id; }
    public ProductionOrderItem getOrderItem() { return orderItem; }
    public Long getSourceOperationId() { return sourceOperationId; }
    public Integer getSequenceNumber() { return sequenceNumber; }
    public String getOperationType() { return operationType; }
    public String getTitle() { return title; }
    public Long getTimePerUnitSeconds() { return timePerUnitSeconds; }
}
