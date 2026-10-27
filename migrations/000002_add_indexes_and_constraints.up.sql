-- Migration: 000002_add_indexes_and_constraints.up.sql
-- Description: Add composite & covering indexes and check constraints for high-throughput AHASS POS

CREATE INDEX IF NOT EXISTS idx_customers_search ON customers(plate_number, name, phone);
CREATE INDEX IF NOT EXISTS idx_spareparts_lookup ON spareparts(part_number, type, stock);
CREATE INDEX IF NOT EXISTS idx_work_orders_lookup ON work_orders(status, payment_status, created_at);
CREATE INDEX IF NOT EXISTS idx_work_order_items_order ON work_order_items(work_order_id, item_type);

ALTER TABLE work_orders ADD CONSTRAINT chk_wo_totals CHECK (grand_total >= 0);
ALTER TABLE spareparts ADD CONSTRAINT chk_spareparts_selling_price CHECK (selling_price >= 0);

