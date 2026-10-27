-- Migration: 000001_create_initial_schema.down.sql
-- Description: Drop all core tables in reverse dependency order

DROP TABLE IF EXISTS work_order_items CASCADE;
DROP TABLE IF EXISTS work_orders CASCADE;
DROP TABLE IF EXISTS spareparts CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS mechanics CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS users CASCADE;

