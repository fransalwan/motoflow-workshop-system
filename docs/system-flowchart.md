# Diagram Alur Sistem & Arsitektur MotoFlow AHASS
PT Nusantara Sakti Group (NSS) — Bengkel Resmi Honda

Dokumen ini memuat diagram alur proses bisnis operasional bengkel AHASS, sequence diagram interaksi komponen, serta diagram relasi entitas basis data (ERD).

---

## 1. Flowchart Alur Operasional Bengkel AHASS

```mermaid
flowchart TD
    Start(["Motor Honda Datang ke AHASS"]) --> SA["Pendaftaran oleh Service Advisor"]
    SA --> RecordCustomer["Catat Nama, No. HP, No. Polisi Plat, Model Motor & KM Odometer"]
    RecordCustomer --> Diagnosis["Diagnosa Kerusakan & Pemilihan Paket Servis / Part HGP"]
    
    Diagnosis --> CreateWO["Terbitkan Work Order (WO-AHASS-YYYYMMDD-XXXX)<br/>Status: Draft"]
    CreateWO --> AssignPit["Alokasikan ke Mekanik & Pit Servis (Pit 1 - 4)"]
    
    AssignPit --> StartWork["Mekanik Mulai Pengerjaan<br/>Status: In Progress"]
    StartWork --> PickParts["Pengambilan Suku Cadang dari Gudang (HGP)<br/>Stok Otomatis Terpotong (ACID Transaction)"]
    PickParts --> MechanicWork["Pengerjaan Servis & Penggantian Suku Cadang"]
    
    MechanicWork --> QualityCheck{"Inspeksi Akhir & Safety Check?"}
    QualityCheck -- Belum Lolos --> MechanicWork
    QualityCheck -- Lolos Inspeksi --> CompleteWO["Work Order Selesai Dikerjakan<br/>Status: Completed"]
    
    CompleteWO --> Cashier["Pelanggan ke Meja Kasir POS Billing"]
    Cashier --> SelectPayment["Pilih Metode Pembayaran<br/>(Tunai / Transfer / QRIS / Debit)"]
    SelectPayment --> ProcessPay["Input Nominal Bayar & Hitung Kembalian Bebas Floating Drift"]
    
    ProcessPay --> Finalize["Status: Paid (Lunas) & Set Timestamp"]
    Finalize --> PrintThermal["Cetak Struk Kasir Thermal (58mm / 80mm)"]
    Finalize --> PrintA4["Cetak Faktur Resmi Servis A4 (Tanda Tangan Ganda)"]
    
    Finalize --> LiveDashboard["Perbarui Dashboard Eksekutif & Omset Harian Secara Real-time"]
    LiveDashboard --> MotorKeluar(["Penyerahan Motor & Garansi Servis 500 KM"])
```

---

## 2. Sequence Diagram Interaksi Komponen Sistem

```mermaid
sequenceDiagram
    autonumber
    actor Pelanggan
    actor SA as Service Advisor
    actor Mekanik
    actor Kasir
    participant Frontend as Vue 3 SPA
    participant Backend as Golang Gin API
    participant DB as PostgreSQL Database

    Pelanggan->>SA: Bawa Motor Honda & Sampaikan Keluhan
    SA->>Frontend: Input Data Pelanggan, Plat Motor & Keluhan
    SA->>Frontend: Pilih Part HGP & Paket Jasa Servis
    Frontend->>Backend: POST /api/v1/work-orders (Payload Items)
    Backend->>DB: BEGIN Transaction & Advisory Lock
    Backend->>DB: Deduct Stock Suku Cadang
    Backend->>DB: INSERT WorkOrder & Items (Price Snapshot)
    Backend->>DB: COMMIT Transaction
    Backend-->>Frontend: 201 Created (WO-AHASS-YYYYMMDD-XXXX)
    
    SA->>Mekanik: Serahkan SPK Servis & Arahkan ke Pit
    Mekanik->>Frontend: Update Status (In Progress)
    Mekanik->>Mekanik: Kerjakan Servis & Ganti Part
    Mekanik->>Frontend: Update Status (Completed - Selesai Inspeksi)
    
    Pelanggan->>Kasir: Menuju Kasir untuk Pembayaran
    Kasir->>Frontend: Buka Rincian WO & Input Nominal Bayar
    Frontend->>Backend: POST /api/v1/work-orders/:id/pay
    Backend->>DB: UPDATE status='paid', payment_status='paid'
    Backend-->>Frontend: 200 OK (Payment Verified)
    
    Kasir->>Frontend: Klik Cetak Struk / Faktur
    Frontend->>Backend: GET /api/v1/work-orders/:id/receipt (Thermal)
    Frontend->>Backend: GET /api/v1/work-orders/:id/invoice-pdf (PDF Stream)
    Kasir-->>Pelanggan: Serahkan Kunci, Struk Thermal & Garansi Servis
```

---

## 3. Entity-Relationship Diagram (ERD)

```mermaid
erDiagram
    USERS ||--o{ WORK_ORDERS : "creates as advisor"
    CUSTOMERS ||--|{ WORK_ORDERS : "owns"
    CATEGORIES ||--|{ SPAREPARTS : "classifies"
    MECHANICS ||--o{ WORK_ORDERS : "executes"
    WORK_ORDERS ||--|{ WORK_ORDER_ITEMS : "contains (cascade delete)"
    SPAREPARTS ||--o{ WORK_ORDER_ITEMS : "referenced by snapshot"

    USERS {
        bigint id PK
        string username UK
        string email UK
        string password_hash
        string full_name
        string role
        string phone
        boolean is_active
        timestamp created_at
    }

    CUSTOMERS {
        bigint id PK
        string name
        string phone
        string plate_number
        string motor_model
        int year
        string frame_number
        string engine_number
        int current_km
        timestamp created_at
    }

    CATEGORIES {
        bigint id PK
        string name UK
        string code UK
        string description
    }

    SPAREPARTS {
        bigint id PK
        string part_number UK
        string name
        bigint category_id FK
        string type
        decimal cost_price
        decimal selling_price
        int stock
        int min_stock
        string unit
        string location
        boolean is_active
    }

    MECHANICS {
        bigint id PK
        string code UK
        string name
        string phone
        string specialization
        int pit_number
        boolean is_available
    }

    WORK_ORDERS {
        bigint id PK
        string wo_number UK
        bigint customer_id FK
        bigint mechanic_id FK
        bigint service_advisor_id FK
        int entry_km
        text complaints
        text service_diagnosis
        string status
        string payment_method
        string payment_status
        decimal parts_subtotal
        decimal services_subtotal
        decimal total_discount
        decimal tax_rate
        decimal tax_amount
        decimal grand_total
        decimal paid_amount
        decimal change_amount
        timestamp completed_at
        timestamp paid_at
    }

    WORK_ORDER_ITEMS {
        bigint id PK
        bigint work_order_id FK
        bigint sparepart_id FK
        string item_type
        string part_number_snapshot
        string item_name_snapshot
        int quantity
        decimal price_snapshot
        decimal discount_percent
        decimal subtotal
        string mechanic_notes
    }
```

