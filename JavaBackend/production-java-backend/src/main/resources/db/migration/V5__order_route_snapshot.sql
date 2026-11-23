CREATE TABLE order_item_operation (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_item_id INTEGER NOT NULL REFERENCES production_order_item(id),
    source_operation_id BIGINT NOT NULL,
    sequence_number INTEGER NOT NULL CHECK (sequence_number > 0),
    operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN
        ('CUTTING', 'TURNING', 'MILLING', 'GRINDING', 'DRILLING', 'INSPECTION')),
    title VARCHAR(200) NOT NULL,
    time_per_unit_seconds BIGINT NOT NULL CHECK (time_per_unit_seconds > 0),
    UNIQUE (order_item_id, source_operation_id),
    UNIQUE (order_item_id, sequence_number)
);
