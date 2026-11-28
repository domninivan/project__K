-- Дев-фикстура для воспроизведения сценария из 03_Go_сценарий_планирования_заказа.docx.
-- НЕ прод-данные: занятость «прочих заказов» выдумана для теста алгоритма.
-- Идентификаторы заказов/позиций/операций совпадают с примером getSchedule
-- из 02_Go_API_контракт_Java_Go.docx, чтобы вывод совпал с примером контракта.

-- 0. Станки из сценария доки 03 (п.3), отсутствующие в сиде 000002.
INSERT INTO work_center (code, title, active)
VALUES
    ('CUT-002',  'Отрезной станок №2', TRUE),
    ('TURN-002', 'Токарный станок №2', TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO work_center_operation_type (work_center_id, operation_type)
SELECT wc.id, v.operation_type
FROM (VALUES
    ('CUT-002',  'CUTTING'),
    ('TURN-002', 'TURNING')
) AS v(code, operation_type)
JOIN work_center wc ON wc.code = v.code
ON CONFLICT (work_center_id, operation_type) DO NOTHING;

-- 1. Все станки работают Пн–Пт 09:00–18:00 (дока 03, п.3).
INSERT INTO work_center_schedule (work_center_id, day_of_week, start_time, end_time)
SELECT wc.id, d.day_of_week, TIME '09:00', TIME '18:00'
FROM work_center wc
CROSS JOIN generate_series(1, 5) AS d(day_of_week)
ON CONFLICT (work_center_id, day_of_week) DO NOTHING;

-- 2. Занятость активных планов на 12.11.2026 (дока 03, п.4).
-- Границы планов = MIN/MAX их операций.
INSERT INTO production_plan (
    external_order_id, due_date, priority,
    planned_start_at, planned_finish_at, on_time, status
) VALUES
    (7,  '2026-11-20', 'MEDIUM',
         '2026-11-12 09:00:00+03', '2026-11-12 16:00:00+03', TRUE, 'ACCEPTED'),
    (8,  '2026-11-20', 'MEDIUM',
         '2026-11-12 09:00:00+03', '2026-11-12 18:00:00+03', TRUE, 'ACCEPTED'),
    (12, '2026-11-20', 'LOW',
         '2026-11-12 09:00:00+03', '2026-11-12 14:00:00+03', TRUE, 'ACCEPTED');

INSERT INTO production_plan_item (
    production_plan_id, external_order_item_id, external_product_id, quantity,
    planned_start_at, planned_finish_at
)
SELECT pp.id, v.order_item_id, v.product_id, 10, pp.planned_start_at, pp.planned_finish_at
FROM (VALUES
    (7,  270, 1),
    (8,  301, 2),
    (12, 410, 4)
) AS v(external_order_id, order_item_id, product_id)
JOIN production_plan pp ON pp.external_order_id = v.external_order_id;

-- Планировщик читает из этих строк только (work_center_id, start_at, finish_at)
-- и статус плана; sequence_number=1 для всех — занятость синтетическая.
INSERT INTO scheduled_operation (
    production_plan_item_id, external_operation_id,
    sequence_number, operation_type, work_center_id, start_at, finish_at
)
SELECT ppi.id, v.op_id, 1, v.operation_type, wc.id,
       v.start_at::timestamptz, v.finish_at::timestamptz
FROM (VALUES
    -- op  | владелец | тип       | станок     | занято (дока 03, п.4)
    (1801, 8,  'CUTTING',  'CUT-001',   '2026-11-12 09:00:00+03', '2026-11-12 12:00:00+03'),
    (1802, 8,  'CUTTING',  'CUT-001',   '2026-11-12 14:00:00+03', '2026-11-12 18:00:00+03'),
    (1901, 12, 'CUTTING',  'CUT-002',   '2026-11-12 09:00:00+03', '2026-11-12 11:00:00+03'),
    (1902, 12, 'CUTTING',  'CUT-002',   '2026-11-12 13:00:00+03', '2026-11-12 14:00:00+03'),
    (1701, 7,  'TURNING',  'TURN-001',  '2026-11-12 09:00:00+03', '2026-11-12 12:00:00+03'),
    (1702, 7,  'TURNING',  'TURN-002',  '2026-11-12 09:00:00+03', '2026-11-12 13:00:00+03'),
    (1703, 7,  'TURNING',  'TURN-001',  '2026-11-12 12:30:00+03', '2026-11-12 16:00:00+03'),
    (1803, 8,  'MILLING',  'MILL-001',  '2026-11-12 09:00:00+03', '2026-11-12 12:00:00+03'),
    (1804, 8,  'MILLING',  'MILL-001',  '2026-11-12 15:00:00+03', '2026-11-12 18:00:00+03'),
    (1805, 8,  'GRINDING', 'GRIND-001', '2026-11-12 09:00:00+03', '2026-11-12 10:00:00+03')
) AS v(op_id, owner_order_id, operation_type, code, start_at, finish_at)
JOIN work_center wc ON wc.code = v.code
JOIN production_plan pp ON pp.external_order_id = v.owner_order_id
JOIN production_plan_item ppi ON ppi.production_plan_id = pp.id;