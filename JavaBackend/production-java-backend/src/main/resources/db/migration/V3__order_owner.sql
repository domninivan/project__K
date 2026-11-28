-- Legacy orders remain unassigned until an explicit owner mapping is supplied.
ALTER TABLE production_order ADD COLUMN client_id BIGINT REFERENCES app_user(id);
CREATE INDEX idx_production_order_client_id ON production_order(client_id);
