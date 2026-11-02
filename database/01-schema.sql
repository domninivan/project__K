-- =========================================================
-- Case 22: Production planning
-- Java backend schema
-- PostgreSQL 17+
-- File: 01-schema.sql
-- =========================================================

-- 1. Product catalog
CREATE TABLE IF NOT EXISTS product (
                                       id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                       title VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(300)
    );

-- 2. Product manufacturing operations
CREATE TABLE IF NOT EXISTS product_operation (
                                                 id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                                 product_id INT NOT NULL,
                                                 sequence_number INT NOT NULL,
                                                 operation_type VARCHAR(20) NOT NULL,
    time_per_unit_seconds INT NOT NULL,

    CONSTRAINT fk_product_operation_product
    FOREIGN KEY (product_id)
    REFERENCES product(id),

    CONSTRAINT positive_operation_time
    CHECK (time_per_unit_seconds > 0),

    CONSTRAINT positive_operation_sequence
    CHECK (sequence_number > 0),

    CONSTRAINT uq_product_operation_sequence
    UNIQUE (product_id, sequence_number)
    );

-- 3. Material requirements per one product unit
-- Material catalog / warehouse / prices are planned to be stored in 1C.
CREATE TABLE IF NOT EXISTS product_material_requirement (
                                                            id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                                            product_id INT NOT NULL,
                                                            material_code VARCHAR(50) NOT NULL,
    quantity_per_unit NUMERIC(12, 3) NOT NULL,
    unit_code VARCHAR(10) NOT NULL,

    CONSTRAINT fk_product_material_requirement_product
    FOREIGN KEY (product_id)
    REFERENCES product(id),

    CONSTRAINT positive_material_quantity
    CHECK (quantity_per_unit > 0),

    CONSTRAINT uq_product_material
    UNIQUE (product_id, material_code)
    );

-- 4. Production order header
CREATE TABLE IF NOT EXISTS production_order (
                                                id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                                due_date DATE NOT NULL,
                                                priority VARCHAR(15) NOT NULL DEFAULT 'MEDIUM',
    status VARCHAR(30) NOT NULL DEFAULT 'NEW',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    comment VARCHAR(200),

    CONSTRAINT priority_check
    CHECK (priority IN (
           'HIGH',
           'MEDIUM',
           'LOW'
                       )),

    CONSTRAINT status_check
    CHECK (status IN (
           'NEW',
           'READY_FOR_PLANNING',
           'PLAN_PROPOSED',
           'PLANNED',
           'IN_PRODUCTION',
           'COMPLETED',
           'CANCELLED'
                     ))
    );

-- 5. Production order items
-- One order can contain several different products.
CREATE TABLE IF NOT EXISTS production_order_item (
                                                     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
                                                     order_id INT NOT NULL,
                                                     product_id INT NOT NULL,
                                                     quantity INT NOT NULL,

                                                     CONSTRAINT fk_order_item_order
                                                     FOREIGN KEY (order_id)
    REFERENCES production_order(id)
    ON DELETE CASCADE,

    CONSTRAINT fk_order_item_product
    FOREIGN KEY (product_id)
    REFERENCES product(id),

    CONSTRAINT positive_order_item_quantity
    CHECK (quantity > 0),

    CONSTRAINT uq_order_product
    UNIQUE (order_id, product_id)
    );

-- =========================================================
-- Indexes
-- =========================================================

CREATE INDEX IF NOT EXISTS idx_product_operation_product_id
    ON product_operation(product_id);

CREATE INDEX IF NOT EXISTS idx_product_material_requirement_product_id
    ON product_material_requirement(product_id);

CREATE INDEX IF NOT EXISTS idx_production_order_due_date
    ON production_order(due_date);

CREATE INDEX IF NOT EXISTS idx_production_order_status
    ON production_order(status);

CREATE INDEX IF NOT EXISTS idx_production_order_item_order_id
    ON production_order_item(order_id);

CREATE INDEX IF NOT EXISTS idx_production_order_item_product_id
    ON production_order_item(product_id);
