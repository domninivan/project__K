package ru.mirea.production.order.entity;
import jakarta.persistence.*;
import ru.mirea.production.order.entity.enums.PriorityEnum;
import ru.mirea.production.order.entity.enums.StatusEnum;

import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;


@Entity
@Table(name = "production_order")
public class ProductionOrder {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Integer id;

    @Column(name = "due_date", nullable = false)
    private LocalDate dueDate;

    @Enumerated(EnumType.STRING)
    @Column(name = "priority", nullable = false)
    private PriorityEnum priority;

    @Enumerated(EnumType.STRING)
    @Column(name = "status", nullable = false)
    private StatusEnum status;

    @Column(name = "created_at", nullable = false)
    private LocalDateTime createdAt;

    @Column(name = "comment", length = 200)
    private String comment;

    @OneToMany(mappedBy = "order", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<ProductionOrderItem> items = new ArrayList<>();

    public ProductionOrder() {
    }

    public ProductionOrder(LocalDate dueDate, String comment) {
        this.dueDate = dueDate;
        this.comment = comment;

        this.priority = PriorityEnum.MEDIUM;
        this.status = StatusEnum.NEW;
        this.createdAt = LocalDateTime.now();
    }

    public void addItem(ProductionOrderItem item) {
        items.add(item);
        item.setOrder(this);
    }

    public Integer getId() {
        return id;
    }

    public LocalDate getDueDate() {
        return dueDate;
    }

    public void setDueDate(LocalDate dueDate) {
        this.dueDate = dueDate;
    }

    public PriorityEnum getPriority() {
        return priority;
    }

    public void setPriority(PriorityEnum priority) {
        this.priority = priority;
    }

    public StatusEnum getStatus() {
        return status;
    }

    public void setStatus(StatusEnum status) {
        this.status = status;
    }

    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    public void setCreatedAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
    }

    public String getComment() {
        return comment;
    }

    public void setComment(String comment) {
        this.comment = comment;
    }

    public List<ProductionOrderItem> getItems() {
        return items;
    }

}
