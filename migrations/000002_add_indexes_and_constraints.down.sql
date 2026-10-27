-- Migration: 000002_add_indexes_and_constraints.down.sql
-- Description: Drop indexes and constraints added in migration 000002

DROP INDEX IF EXISTS idx_customers_search;
DROP INDEX IF EXISTS idx_spareparts_lookup;
DROP INDEX IF EXISTS idx_work_orders_lookup;
DROP INDEX IF EXISTS idx_work_order_items_order;

ALTER TABLE work_orders DROP CONSTRAINT IF EXISTS chk_wo_totals;
ALTER TABLE spareparts DROP CONSTRAINT IF EXISTS chk_spareparts_selling_price;

