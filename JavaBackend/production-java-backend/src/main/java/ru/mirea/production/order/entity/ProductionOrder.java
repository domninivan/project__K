package ru.mirea.production.order.entity;
import jakarta.persistence.*;
import ru.mirea.production.order.entity.enums.PriorityEnum;
import ru.mirea.production.order.entity.enums.StatusEnum;
import ru.mirea.production.product.entity.Product;

import java.time.LocalDate;
import java.time.LocalDateTime;


@Entity
@Table(name = "production_order")
public class ProductionOrder {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Integer id;

    @Column(name = "quantity", nullable = false)
    private Integer quantity;

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

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "product_id", nullable = false)
    private Product product;

    public ProductionOrder() {}

    @Override
    public String toString() {
        return "ProductionOrder{" +
                "id=" + id +
                ", quantity=" + quantity +
                ", dueDate=" + dueDate +
                ", priority=" + priority +
                ", status=" + status +
                ", createdAt=" + createdAt +
                ", comment='" + comment + '\'' +
                ", productId=" + product.getId() +
                '}';
    }

    public ProductionOrder(Integer quantity, LocalDate dueDate, PriorityEnum priority, String comment, Product product) {
        this.quantity = quantity;
        this.dueDate = dueDate;
        this.priority = priority;
        this.status = StatusEnum.NEW;
        this.createdAt = LocalDateTime.now();
        this.comment = comment;
        this.product = product;
    }

    public Integer getId() {
        return id;
    }

    public Integer getQuantity() {
        return quantity;
    }

    public void setQuantity(Integer quantity) {
        this.quantity = quantity;
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

    public String getComment() {
        return comment;
    }

    public void setComment(String comment) {
        this.comment = comment;
    }

    public Product getProduct() {
        return product;
    }

    public void setProduct(Product product) {
        this.product = product;
    }
}
