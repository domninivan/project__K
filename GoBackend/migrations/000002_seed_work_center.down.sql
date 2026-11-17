-- Удаляем только тестовые центры; операции уйдут каскадом (ON DELETE CASCADE)
DELETE FROM work_center
WHERE code IN (
    'CUT-001', 'TURN-001', 'MILL-001',
    'GRIND-001', 'DRILL-001', 'QC-001'
);.