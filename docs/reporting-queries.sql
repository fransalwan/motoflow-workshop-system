-- ============================================================================
-- SHOWCASE QUERY SQL AGREGASI & PELAPORAN EKSEKUTIF AHASS
-- PT NUSANTARA SAKTI GROUP (NSS) - BENGKEL RESMI HONDA
-- Kualifikasi Loker: PostgreSQL (SELECT, JOIN, GROUP BY, SUM, COUNT, AVG, WINDOW)
-- File: docs/reporting-queries.sql
-- ============================================================================

-- ============================================================================
-- QUERY 1: LAPORAN OMSET HARIAN & BULANAN BENGKEL
-- Tujuan: Menghitung total omset kasir, unit motor masuk, rata-rata tiket transaksi,
--         serta akumulasi pendapatan bulanan secara agregatif.
-- Klausa: SELECT, SUM, COUNT, AVG, GROUP BY, ORDER BY, DATE_TRUNC
-- ============================================================================

SELECT 
    DATE(wo.paid_at) AS tanggal_transaksi,
    COUNT(wo.id) AS total_unit_servis,
    SUM(wo.parts_subtotal) AS omset_suku_cadang,
    SUM(wo.services_subtotal) AS omset_jasa_servis,
    SUM(wo.tax_amount) AS total_ppn_11,
    SUM(wo.grand_total) AS total_omset_bruto,
    ROUND(AVG(wo.grand_total), 2) AS rata_rata_per_motor
FROM 
    work_orders wo
WHERE 
    wo.payment_status = 'paid'
    AND wo.deleted_at IS NULL
    AND wo.paid_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY 
    DATE(wo.paid_at)
ORDER BY 
    tanggal_transaksi DESC;

-- Rekapitulasi Tingkat Bulanan (Executive MTD / YTD):
SELECT 
    TO_CHAR(DATE_TRUNC('month', wo.paid_at), 'Month YYYY') AS periode_bulan,
    COUNT(wo.id) AS total_motor_diservis,
    SUM(wo.grand_total) AS total_omset_bulanan,
    ROUND(AVG(wo.grand_total), 2) AS rata_rata_omset_harian
FROM 
    work_orders wo
WHERE 
    wo.payment_status = 'paid'
    AND wo.deleted_at IS NULL
GROUP BY 
    DATE_TRUNC('month', wo.paid_at)
ORDER BY 
    DATE_TRUNC('month', wo.paid_at) DESC;


-- ============================================================================
-- QUERY 2: TOP 5 HONDA GENUINE PARTS (HGP) TERLARIS
-- Tujuan: Mengidentifikasi fast-moving suku cadang dengan kontribusi pendapatan tertinggi
--         guna optimasi supply chain dan restock suku cadang dari PT Astra Honda Motor (AHM).
-- Klausa: SELECT, JOIN, WHERE, GROUP BY, SUM, ORDER BY, LIMIT
-- ============================================================================

SELECT 
    sp.part_number AS nomor_part_hgp,
    sp.name AS nama_suku_cadang,
    cat.name AS kategori,
    SUM(woi.quantity) AS total_kuantitas_terjual,
    sp.stock AS sisa_stok_gudang,
    SUM(woi.subtotal) AS total_omset_penjualan,
    ROUND(
        SUM(woi.subtotal) / NULLIF(SUM(SUM(woi.subtotal)) OVER(), 0) * 100, 
        2
    ) AS persentase_kontribusi_omset
FROM 
    work_order_items woi
INNER JOIN 
    work_orders wo ON wo.id = woi.work_order_id
INNER JOIN 
    spareparts sp ON sp.id = woi.sparepart_id
INNER JOIN 
    categories cat ON cat.id = sp.category_id
WHERE 
    wo.payment_status = 'paid'
    AND woi.item_type = 'part'
    AND woi.deleted_at IS NULL
    AND wo.deleted_at IS NULL
GROUP BY 
    sp.part_number, 
    sp.name, 
    cat.name, 
    sp.stock
ORDER BY 
    total_omset_penjualan DESC
LIMIT 5;


-- ============================================================================
-- QUERY 3: KINERJA & PRODUKTIVITAS MEKANIK AHASS
-- Tujuan: Mengukur efektivitas dan produktivitas mekanik per Pit pengerjaan,
--         total unit motor yang diselesaikan, dan omset jasa servis yang dihasilkan.
-- Klausa: SELECT, LEFT JOIN, COUNT, SUM, AVG, GROUP BY, ORDER BY
-- ============================================================================

SELECT 
    m.code AS kode_mekanik,
    m.name AS nama_mekanik,
    m.pit_number AS nomor_pit,
    m.specialization AS spesialisasi,
    COUNT(wo.id) AS total_motor_diselesaikan,
    COALESCE(SUM(wo.services_subtotal), 0) AS total_pendapatan_jasa,
    COALESCE(ROUND(AVG(wo.services_subtotal), 2), 0) AS rata_rata_pendapatan_per_motor,
    ROUND(
        AVG(
            EXTRACT(EPOCH FROM (wo.completed_at - wo.created_at)) / 60
        ), 
        1
    ) AS rata_rata_menit_pengerjaan
FROM 
    mechanics m
LEFT JOIN 
    work_orders wo ON wo.mechanic_id = m.id 
    AND wo.status IN ('completed', 'paid')
    AND wo.deleted_at IS NULL
WHERE 
    m.deleted_at IS NULL
GROUP BY 
    m.id, 
    m.code, 
    m.name, 
    m.pit_number, 
    m.specialization
ORDER BY 
    total_pendapatan_jasa DESC;


-- ============================================================================
-- QUERY 4: NOTIFIKASI PERINGATAN SUKU CADANG MENIPIS (LOW STOCK ALERT)
-- Tujuan: Memberikan early warning trigger pada dashboard kepala bengkel saat
--         stok fisik berada pada atau di bawah batas minimum pemesanan ulang (Reorder Point).
-- Klausa: SELECT, JOIN, WHERE (stock <= min_stock), ORDER BY
-- ============================================================================

SELECT 
    sp.part_number AS nomor_part_hgp,
    sp.name AS nama_suku_cadang,
    cat.name AS kategori,
    sp.stock AS stok_saat_ini,
    sp.min_stock AS batas_minimum_stok,
    (sp.min_stock - sp.stock) AS defisit_kekurangan_stok,
    sp.unit AS satuan,
    sp.location AS lokasi_rak,
    sp.cost_price AS estimasi_harga_beli_ahm,
    ROUND((sp.min_stock - sp.stock + 5) * sp.cost_price, 2) AS estimasi_modal_restock
FROM 
    spareparts sp
INNER JOIN 
    categories cat ON cat.id = sp.category_id
WHERE 
    sp.type = 'part'
    AND sp.stock <= sp.min_stock
    AND sp.is_active = TRUE
    AND sp.deleted_at IS NULL
ORDER BY 
    (sp.min_stock - sp.stock) DESC, 
    sp.stock ASC;


-- ============================================================================
-- QUERY 5: ANALISA BREAKDOWN PENDAPATAN JASA VS SUKU CADANG (REVENUE MIX)
-- Tujuan: Analisa komparasi porsi pendapatan antara layanan jasa bengkel dengan
--         penjualan suku cadang untuk evaluasi rasio laba margin AHASS NSS.
-- Klausa: SELECT, CASE WHEN, SUM, ROUND
-- ============================================================================

SELECT 
    COUNT(wo.id) AS total_transaksi,
    SUM(wo.parts_subtotal) AS total_omset_parts,
    ROUND(
        (SUM(wo.parts_subtotal) / NULLIF(SUM(wo.parts_subtotal + wo.services_subtotal), 0)) * 100, 
        2
    ) AS persentase_parts,
    SUM(wo.services_subtotal) AS total_omset_jasa,
    ROUND(
        (SUM(wo.services_subtotal) / NULLIF(SUM(wo.parts_subtotal + wo.services_subtotal), 0)) * 100, 
        2
    ) AS persentase_jasa,
    SUM(wo.tax_amount) AS total_ppn,
    SUM(wo.grand_total) AS grand_total_keseluruhan
FROM 
    work_orders wo
WHERE 
    wo.payment_status = 'paid'
    AND wo.deleted_at IS NULL;

