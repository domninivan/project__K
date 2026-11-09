-- Тестовые рабочие центры
INSERT INTO work_center (code, title, active)
VALUES
    ('CUT-001',   'Отрезной станок №1',     TRUE),
    ('TURN-001',  'Токарный станок №1',     TRUE),
    ('MILL-001',  'Фрезерный станок №1',    TRUE),
    ('GRIND-001', 'Шлифовальный станок №1', TRUE),
    ('DRILL-001', 'Сверлильный станок №1',  TRUE),
    ('QC-001',    'Пост контроля №1',       TRUE)
ON CONFLICT (code) DO NOTHING;

-- Операции для этих центров (id берём по code, чтобы не зависеть от автоинкремента)
INSERT INTO work_center_operation_type (work_center_id, operation_type)
SELECT wc.id, v.operation_type
FROM (VALUES
    ('CUT-001',   'CUTTING'),
    ('TURN-001',  'TURNING'),
    ('MILL-001',  'MILLING'),
    ('GRIND-001', 'GRINDING'),
    ('DRILL-001', 'DRILLING'),
    ('QC-001',    'INSPECTION')
) AS v(code, operation_type)
JOIN work_center wc ON wc.code = v.code
ON CONFLICT (work_center_id, operation_type) DO NOTHING;