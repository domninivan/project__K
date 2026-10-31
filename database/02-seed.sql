-- =========================================================
-- Case 22: Production planning
-- Java backend demo / seed data
-- File: 02-seed.sql
-- =========================================================

-- =========================================================
-- 1. Products
-- =========================================================

INSERT INTO product (title, description)
VALUES
    ('Вал приводной', 'Деталь для передачи крутящего момента'),
    ('Втулка', 'Опорная цилиндрическая деталь'),
    ('Шестерня', 'Зубчатая деталь механической передачи'),
    ('Кронштейн', 'Крепёжная деталь для фиксации узлов')
    ON CONFLICT (title) DO NOTHING;


-- =========================================================
-- 2. Product operations
-- time_per_unit_seconds = normative time for one product unit
-- =========================================================

-- Вал приводной
INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 1, 'CUTTING', 45
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 2, 'TURNING', 120
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 3, 'MILLING', 90
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 4, 'GRINDING', 60
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;


-- Втулка
INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 1, 'CUTTING', 30
FROM product
WHERE title = 'Втулка'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 2, 'TURNING', 75
FROM product
WHERE title = 'Втулка'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 3, 'DRILLING', 40
FROM product
WHERE title = 'Втулка'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;


-- Шестерня
INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 1, 'CUTTING', 50
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 2, 'TURNING', 100
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 3, 'MILLING', 180
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 4, 'GRINDING', 90
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 5, 'INSPECTION', 30
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;


-- Кронштейн
INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 1, 'CUTTING', 35
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 2, 'MILLING', 110
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 3, 'DRILLING', 55
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;

INSERT INTO product_operation
(product_id, sequence_number, operation_type, time_per_unit_seconds)
SELECT id, 4, 'INSPECTION', 25
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, sequence_number) DO NOTHING;


-- =========================================================
-- 3. Material requirements
-- =========================================================

-- Вал приводной
INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'STEEL_40X', 2.400, 'KG'
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, material_code) DO NOTHING;

INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'LUBRICANT_M01', 0.030, 'L'
FROM product
WHERE title = 'Вал приводной'
    ON CONFLICT (product_id, material_code) DO NOTHING;


-- Втулка
INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'BRONZE_BR02', 0.450, 'KG'
FROM product
WHERE title = 'Втулка'
    ON CONFLICT (product_id, material_code) DO NOTHING;

INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'LUBRICANT_M01', 0.010, 'L'
FROM product
WHERE title = 'Втулка'
    ON CONFLICT (product_id, material_code) DO NOTHING;


-- Шестерня
INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'STEEL_20X', 1.800, 'KG'
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, material_code) DO NOTHING;

INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'QUENCHING_OIL', 0.050, 'L'
FROM product
WHERE title = 'Шестерня'
    ON CONFLICT (product_id, material_code) DO NOTHING;


-- Кронштейн
INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'STEEL_SHEET_3MM', 1.250, 'KG'
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, material_code) DO NOTHING;

INSERT INTO product_material_requirement
(product_id, material_code, quantity_per_unit, unit_code)
SELECT id, 'PAINT_BLACK', 0.060, 'L'
FROM product
WHERE title = 'Кронштейн'
    ON CONFLICT (product_id, material_code) DO NOTHING;


-- =========================================================
-- 4. Demo production orders
-- Demo comments are used to avoid duplicate insertion.
-- =========================================================

INSERT INTO production_order
(due_date, priority, status, comment)
SELECT
    DATE '2026-11-15',
    'MEDIUM',
    'NEW',
    'DEMO_ORDER_001'
    WHERE NOT EXISTS (
    SELECT 1
    FROM production_order
    WHERE comment = 'DEMO_ORDER_001'
);

INSERT INTO production_order
(due_date, priority, status, comment)
SELECT
    DATE '2026-11-20',
    'HIGH',
    'NEW',
    'DEMO_ORDER_002'
    WHERE NOT EXISTS (
    SELECT 1
    FROM production_order
    WHERE comment = 'DEMO_ORDER_002'
);


-- =========================================================
-- 5. Demo production order items
-- =========================================================

-- DEMO_ORDER_001:
-- 100 drive shafts + 50 bushings
INSERT INTO production_order_item
(order_id, product_id, quantity)
SELECT
    o.id,
    p.id,
    100
FROM production_order o
         JOIN product p ON p.title = 'Вал приводной'
WHERE o.comment = 'DEMO_ORDER_001'
    ON CONFLICT (order_id, product_id) DO NOTHING;

INSERT INTO production_order_item
(order_id, product_id, quantity)
SELECT
    o.id,
    p.id,
    50
FROM production_order o
         JOIN product p ON p.title = 'Втулка'
WHERE o.comment = 'DEMO_ORDER_001'
    ON CONFLICT (order_id, product_id) DO NOTHING;


-- DEMO_ORDER_002:
-- 30 gears + 20 brackets
INSERT INTO production_order_item
(order_id, product_id, quantity)
SELECT
    o.id,
    p.id,
    30
FROM production_order o
         JOIN product p ON p.title = 'Шестерня'
WHERE o.comment = 'DEMO_ORDER_002'
    ON CONFLICT (order_id, product_id) DO NOTHING;

INSERT INTO production_order_item
(order_id, product_id, quantity)
SELECT
    o.id,
    p.id,
    20
FROM production_order o
         JOIN product p ON p.title = 'Кронштейн'
WHERE o.comment = 'DEMO_ORDER_002'
    ON CONFLICT (order_id, product_id) DO NOTHING;


-- =========================================================
-- 6. Optional verification queries
-- =========================================================
-- SELECT * FROM product ORDER BY id;
-- SELECT * FROM product_operation ORDER BY product_id, sequence_number;
-- SELECT * FROM product_material_requirement ORDER BY product_id, id;
-- SELECT * FROM production_order ORDER BY id;
-- SELECT * FROM production_order_item ORDER BY order_id, id;
