package ru.mirea.production.product.entity;
import jakarta.persistence.*;

@Entity
@Table(name = "product_operation")
public class ProductOperation {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name = "id")
    private Integer id;

    @Column(name = "sequence_number", nullable = false)
    private Integer sequenceNumber;

    @Column(name = "operation_type", nullable = false, length = 20)
    private String operationType;

    @Column(name = "time_per_unit_seconds", nullable = false)
    private Integer timePerUnitSeconds;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "product_id", nullable = false)
    private Product product;

    public ProductOperation() {}

    public ProductOperation(Integer sequenceNumber, String operationType, Integer timePerUnitSeconds, Product product) {
        this.sequenceNumber = sequenceNumber;
        this.operationType = operationType;
        this.timePerUnitSeconds = timePerUnitSeconds;
        this.product = product;
    }

    public Integer getId() {
        return id;
    }

    public Integer getSequenceNumber() {
        return sequenceNumber;
    }

    public void setSequenceNumber(Integer sequenceNumber) {
        this.sequenceNumber = sequenceNumber;
    }

    public String getOperationType() {
        return operationType;
    }

    public void setOperationType(String operationType) {
        this.operationType = operationType;
    }

    public Integer getTimePerUnitSeconds() {
        return timePerUnitSeconds;
    }

    public void setTimePerUnitSeconds(Integer timePerUnitSeconds) {
        this.timePerUnitSeconds = timePerUnitSeconds;
    }

    public Product getProduct() {
        return product;
    }

    public void setProduct(Product product) {
        this.product = product;
    }

    @Override
    public String toString() {
        return "ProductOperation{" +
                "id=" + id +
                ", sequenceNumber=" + sequenceNumber +
                ", operationType='" + operationType + '\'' +
                ", timePerUnitSeconds=" + timePerUnitSeconds +
                '}';
    }
}
