package ru.mirea.production.product.entity;

import jakarta.persistence.*;
import ru.mirea.production.product.entity.enums.UnitOfMeasure;
import java.math.BigDecimal;

@Entity
@Table(name="product_material_requirement")
public class ProductMaterialRequirement {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    @Column(name="id")
    private Integer id;

    @Column(name="material_code", nullable = false, length = 50)
    private String materialCode;

    @Column(name="quantity_per_unit", nullable = false, precision = 12, scale = 3)
    private BigDecimal quantityPerUnit;

    @Enumerated(EnumType.STRING)
    @Column(name="unit_code", nullable = false)
    private UnitOfMeasure unitCode;

    @ManyToOne(fetch = FetchType.LAZY, optional = false)
    @JoinColumn(name = "product_id", nullable = false)
    private Product product;

    public ProductMaterialRequirement(){
    }

    public ProductMaterialRequirement(String materialCode, BigDecimal quantityPerUnit, UnitOfMeasure unitCode, Product product) {
        this.materialCode = materialCode;
        this.quantityPerUnit = quantityPerUnit;
        this.unitCode = unitCode;
        this.product = product;
    }

    public Integer getId() {
        return id;
    }

    public String getMaterialCode() {
        return materialCode;
    }

    public void setMaterialCode(String materialCode) {
        this.materialCode = materialCode;
    }

    public BigDecimal getQuantityPerUnit() {
        return quantityPerUnit;
    }

    public void setQuantityPerUnit(BigDecimal quantityPerUnit) {
        this.quantityPerUnit = quantityPerUnit;
    }

    public UnitOfMeasure getUnitCode() {
        return unitCode;
    }

    public void setUnitCode(UnitOfMeasure unitCode) {
        this.unitCode = unitCode;
    }

    public Product getProduct() {
        return product;
    }

    public void setProduct(Product product) {
        this.product = product;
    }

    @Override
    public String toString() {
        return "ProductMaterialRequirement{" +
                "unitCode=" + unitCode +
                ", quantityPerUnit=" + quantityPerUnit +
                ", materialCode='" + materialCode + '\'' +
                ", id=" + id +
                '}';
    }
}
