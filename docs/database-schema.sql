-- ============================================================================
-- SISTEM INFORMASI MANAJEMEN BENGKEL RESMI HONDA (AHASS) & KASIR POS BILLING
-- PT NUSANTARA SAKTI GROUP (NSS) - DEALER RESMI & BENGKEL AHASS HONDA TERBESAR
-- Database: PostgreSQL 14+ / 16
-- File: docs/database-schema.sql
-- ============================================================================

-- Create Extensions (if needed)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Drop tables if exists in reverse dependency order
DROP TABLE IF EXISTS work_order_items CASCADE;
DROP TABLE IF EXISTS work_orders CASCADE;
DROP TABLE IF EXISTS spareparts CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS mechanics CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ----------------------------------------------------------------------------
-- 1. TABEL PENGGUNA & STAF BENGKEL (USERS)
-- Role: admin (Kepala Bengkel), service_advisor, mechanic, cashier (Kasir POS)
-- ----------------------------------------------------------------------------
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role VARCHAR(30) NOT NULL DEFAULT 'service_advisor',
    phone VARCHAR(25),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_users_username UNIQUE (username),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT chk_users_role CHECK (role IN ('admin', 'service_advisor', 'mechanic', 'cashier'))
);

CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

COMMENT ON TABLE users IS 'Data staf operasional bengkel resmi Honda AHASS PT Nusantara Sakti';

-- ----------------------------------------------------------------------------
-- 2. TABEL PELANGGAN & KENDARAAN HONDA (CUSTOMERS)
-- Menyimpan identitas pemilik dan spesifikasi motor (Beat, Vario 160, PCX 160, dll)
-- ----------------------------------------------------------------------------
CREATE TABLE customers (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(25) NOT NULL,
    email VARCHAR(100),
    address TEXT,
    plate_number VARCHAR(20) NOT NULL, -- Contoh: B 3456 KRF
    motor_model VARCHAR(60) NOT NULL,  -- Contoh: Honda Vario 160 ABS, Honda BeAT CBS
    year INT,                          -- Tahun pembuatan
    frame_number VARCHAR(50),          -- Nomor Rangka (VIN)
    engine_number VARCHAR(50),         -- Nomor Mesin
    current_km INT DEFAULT 0,          -- Odometer terakhir
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_customers_name ON customers(name);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_plate ON customers(plate_number);
CREATE INDEX idx_customers_motor_model ON customers(motor_model);
CREATE INDEX idx_customers_deleted_at ON customers(deleted_at);

COMMENT ON TABLE customers IS 'Data master pelanggan dan unit motor Honda yang masuk bengkel AHASS';

-- ----------------------------------------------------------------------------
-- 3. TABEL KATEGORI SUKU CADANG & JASA (CATEGORIES)
-- ----------------------------------------------------------------------------
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(30) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_categories_name UNIQUE (name),
    CONSTRAINT uq_categories_code UNIQUE (code)
);

CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

COMMENT ON TABLE categories IS 'Kategori suku cadang resmi AHM dan paket layanan jasa servis';

-- ----------------------------------------------------------------------------
-- 4. TABEL SUKU CADANG (HONDA GENUINE PARTS) & JASA (SPAREPARTS)
-- Presisi decimal(12,2) untuk mencegah floating-point rounding drift
-- ----------------------------------------------------------------------------
CREATE TABLE spareparts (
    id BIGSERIAL PRIMARY KEY,
    part_number VARCHAR(50) NOT NULL, -- Contoh: 08232-2MB-K0NZ1 (AHM Oil MPX2)
    name VARCHAR(150) NOT NULL,
    category_id BIGINT NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'part', -- 'part' atau 'service'
    cost_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    selling_price DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    stock INT NOT NULL DEFAULT 0,
    min_stock INT NOT NULL DEFAULT 5, -- Alert threshold untuk stok menipis
    unit VARCHAR(20) DEFAULT 'Pcs',   -- Pcs, Botol, Set, Paket
    location VARCHAR(60),             -- Rak A-01, Area Servis
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_spareparts_part_number UNIQUE (part_number),
    CONSTRAINT fk_spareparts_category FOREIGN KEY (category_id) REFERENCES categories(id) ON UPDATE CASCADE,
    CONSTRAINT chk_spareparts_type CHECK (type IN ('part', 'service'))
);

CREATE INDEX idx_spareparts_category_id ON spareparts(category_id);
CREATE INDEX idx_spareparts_type ON spareparts(type);
CREATE INDEX idx_spareparts_stock ON spareparts(stock);
CREATE INDEX idx_spareparts_deleted_at ON spareparts(deleted_at);

COMMENT ON TABLE spareparts IS 'Katalog Honda Genuine Parts (HGP), Oli AHM, dan Paket Jasa Servis';

-- ----------------------------------------------------------------------------
-- 5. TABEL MEKANIK & PIT BENGKEL (MECHANICS)
-- ----------------------------------------------------------------------------
CREATE TABLE mechanics (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL, -- Contoh: MEK-01
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(25),
    specialization VARCHAR(100), -- Contoh: Injeksi PGM-FI & Kelistrikan
    pit_number INT NOT NULL DEFAULT 1,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_mechanics_code UNIQUE (code)
);

CREATE INDEX idx_mechanics_pit_number ON mechanics(pit_number);
CREATE INDEX idx_mechanics_available ON mechanics(is_available);
CREATE INDEX idx_mechanics_deleted_at ON mechanics(deleted_at);

COMMENT ON TABLE mechanics IS 'Daftar mekanik bersertifikat Honda AHASS beserta alokasi Pit pengerjaan';

-- ----------------------------------------------------------------------------
-- 6. TABEL PERINTAH KERJA & FAKTUR KASIR (WORK_ORDERS)
-- Concurrency lock anti-race: WO-AHASS-YYYYMMDD-XXXX
-- ----------------------------------------------------------------------------
CREATE TABLE work_orders (
    id BIGSERIAL PRIMARY KEY,
    wo_number VARCHAR(50) NOT NULL,
    customer_id BIGINT NOT NULL,
    mechanic_id BIGINT,
    service_advisor_id BIGINT,
    entry_km INT DEFAULT 0,
    complaints TEXT,
    service_diagnosis TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'draft',
    payment_method VARCHAR(30) DEFAULT 'cash',
    payment_status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
    parts_subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    services_subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    total_discount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    tax_rate DECIMAL(5,2) NOT NULL DEFAULT 11.00, -- PPN 11%
    tax_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    grand_total DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    change_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    completed_at TIMESTAMP WITH TIME ZONE,
    paid_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_work_orders_wo_number UNIQUE (wo_number),
    CONSTRAINT fk_work_orders_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON UPDATE CASCADE,
    CONSTRAINT fk_work_orders_mechanic FOREIGN KEY (mechanic_id) REFERENCES mechanics(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_work_orders_service_advisor FOREIGN KEY (service_advisor_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT chk_wo_status CHECK (status IN ('draft', 'in_progress', 'completed', 'paid', 'cancelled')),
    CONSTRAINT chk_wo_payment_status CHECK (payment_status IN ('unpaid', 'paid'))
);

CREATE INDEX idx_work_orders_customer_id ON work_orders(customer_id);
CREATE INDEX idx_work_orders_mechanic_id ON work_orders(mechanic_id);
CREATE INDEX idx_work_orders_status ON work_orders(status);
CREATE INDEX idx_work_orders_payment_status ON work_orders(payment_status);
CREATE INDEX idx_work_orders_created_at ON work_orders(created_at);
CREATE INDEX idx_work_orders_paid_at ON work_orders(paid_at);
CREATE INDEX idx_work_orders_deleted_at ON work_orders(deleted_at);

COMMENT ON TABLE work_orders IS 'Surat Perintah Kerja (SPK) servis motor dan faktur tagihan kasir';

-- ----------------------------------------------------------------------------
-- 7. TABEL RINCIAN ITEM SERVIS & SPAREPART (WORK_ORDER_ITEMS)
-- Memiliki relasi Cascade Delete ke work_orders dan snapshot harga historis
-- ----------------------------------------------------------------------------
CREATE TABLE work_order_items (
    id BIGSERIAL PRIMARY KEY,
    work_order_id BIGINT NOT NULL,
    sparepart_id BIGINT,
    item_type VARCHAR(20) NOT NULL DEFAULT 'part',
    part_number_snapshot VARCHAR(50) NOT NULL,
    item_name_snapshot VARCHAR(150) NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    price_snapshot DECIMAL(12,2) NOT NULL DEFAULT 0.00, -- Snapshot harga transaksi
    discount_percent DECIMAL(5,2) DEFAULT 0.00,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    mechanic_notes VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_wo_items_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id) ON DELETE CASCADE,
    CONSTRAINT fk_wo_items_sparepart FOREIGN KEY (sparepart_id) REFERENCES spareparts(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT chk_wo_items_type CHECK (item_type IN ('part', 'service')),
    CONSTRAINT chk_wo_items_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_wo_items_work_order_id ON work_order_items(work_order_id);
CREATE INDEX idx_wo_items_sparepart_id ON work_order_items(sparepart_id);
CREATE INDEX idx_wo_items_item_type ON work_order_items(item_type);
CREATE INDEX idx_wo_items_deleted_at ON work_order_items(deleted_at);

COMMENT ON TABLE work_order_items IS 'Rincian suku cadang dan jasa mekanik per work order (Cascade Delete)';

