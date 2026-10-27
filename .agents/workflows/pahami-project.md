# Workflow: Pahami Proyek MotoFlow AHASS
### Panduan Orientasi Komprehensif Arsitektur, Alur Bisnis, dan Basis Kode

Workflow ini wajib dibaca saat pertama kali membuka repositori untuk memahami seluruh konteks bisnis, arsitektur teknis, dan alur operasional sistem.

---

## 1. Profil & Tujuan Proyek

* **Nama Proyek**: **MotoFlow AHASS Workshop & Retail System**
* **Target Profil**: Portofolio posisi **Programmer Junior** di **PT Nusantara Sakti Group (NSS)**
* **Spesialisasi Domain**:
  1. *Dealer Resmi Motor Honda*: Showroom fisik unit baru, Nomor Rangka pabrikan Astra (*VIN*), Nomor Mesin.
  2. *Bengkel Resmi Honda AHASS*: Tiket servis Work Order 4 tahap, alokasi Pit mekanik 1–4, struk thermal ESC/POS, faktur resmi A4.
  3. *Retail Suku Cadang Over-The-Counter (HGP)*: Kasir langsung oli AHM MPX/SPX & suku cadang tanpa antrean servis mekanik.
  4. *Pembiayaan NSC Multifinance (PT Nusa Surya Ciptadana)*: Simulasi kredit cicilan motor, DP OJK 10–50%, tenor 11–35 bulan, bunga flat 18% p.a., biaya fidusia & asuransi TLO, hingga penerbitan SPK resmi.

---

## 2. Peta Struktur Monorepo

```
motoflow-workshop-system/
├── .agents/                           # Instruksi Agen AI & Aturan Rekayasa
│   ├── rules/                         # Aturan Arsitektur & Coding Standard (scrapflow-core.md)
│   └── workflows/                     # SOP Alur Kerja (mulai-task.md, pahami-project.md)
│
├── .github/workflows/ci.yml           # Automated CI/CD (Linter, Unit Tests, Go Build, npm build)
│
├── motoflow-backend/                  # RESTful API Service (Golang 1.24 + Gin Framework)
│   ├── config/                        # Database PostgreSQL 16 & SQLite In-Memory Engine
│   ├── controllers/                   # Handler HTTP (Auth, WO, Retail, Showroom, Print, Dashboard)
│   ├── dto/                           # Data Transfer Object & Skema Validasi Request/Response
│   ├── middleware/                    # JWT Authentication, Role Permissions & CORS
│   ├── models/                        # GORM Entities (User, Customer, Sparepart, WorkOrder, SPK, dll)
│   ├── routes/                        # Registrasi Route Terpusat Gin (/api/v1/*)
│   ├── seeders/                       # Seeder Data Realistis Pelanggan, Motor, HGP & Mekanik
│   ├── services/                      # Finansial Decimal Engine, Cache Service, & Table-Driven Tests
│   ├── e2e_test.go                    # Automated End-to-End Test Suite 6 Skenario Bisnis Riil
│   └── main.go                        # Entry Point Server (Port 8080)
│
├── motoflow-frontend/                 # Single Page Application (Vue 3.4 + TypeScript + Vite)
│   ├── src/
│   │   ├── components/                # ThermalReceiptModal, A4InvoiceModal, QRISModal, Sidebar, Navbar
│   │   ├── layouts/                   # MainLayout (Safe-area aware Bottom Navigation Bar)
│   │   ├── pages/                     # Dashboard, WorkOrders, Retail POS, Showroom, Inventory, dll
│   │   ├── router/                    # Guard Autentikasi & Dynamic Role-Based Route Filtering
│   │   ├── stores/                    # Pinia Stores (Auth, WorkOrder, Retail, Showroom, Toast, dll)
│   │   ├── types/                     # Interface TypeScript Terpusat
│   │   └── utils/                     # printer.ts (Isolated Print), format.ts, sound.ts (Web Audio)
│   ├── index.html                     # Mobile Viewport-fit cover & PWA theme color #CC0000
│   └── package.json
│
└── docs/                              # Dokumen Showcase untuk HRD / Tech Lead NSS
    ├── api-documentation.md           # Referensi Endpoint RESTful API & JSON Payload
    ├── business-architecture.md       # Spesifikasi Domain NSS, Finansial & Formula Leasing NSC
    ├── database-schema.sql            # Skema Relasi DDL, Foreign Key, & Index PostgreSQL
    ├── reporting-queries.sql          # Showcase Query SQL: JOIN, GROUP BY, SUM, COUNT, AVG
    └── system-flowchart.md            # Diagram Alur Bisnis Servis AHASS (Mermaid)
```

---

## 3. Empat Peran Pengguna (Role-Based Access Control)

| Peran | Username / Password | Hak Akses Utama | Halaman Default |
| :--- | :--- | :--- | :--- |
| **Admin (Kepala Bengkel)** | `admin` / `password123` | Akses penuh dashboard omset eksekutif, master data, laporan laba rugi | `/dashboard` |
| **Service Advisor (SA)** | `advisor` / `password123` | Penerimaan motor masuk, diagnosis keluhan, penugasan pit mekanik | `/work-orders` |
| **Kasir POS AHASS** | `kasir` / `password123` | Penyelesaian billing servis, kasir retail OTC, promo voucher, cetak struk | `/pos-retail` |
| **Mekanik Senior** | `mekanik` / `password123` | Pemantauan Kanban Pit Board, pengerjaan servis motor, inspeksi QC | `/work-orders` |

---

## 4. Tiga Skenario Alur Bisnis Utama

### A. Alur Servis Motor Honda AHASS
1. **Tahap 1: Pendaftaran SA (`draft`)**: Konsumen datang, SA mencatat nomor plat (`B 3456 KRF`), model motor, kilometer odometer, diagnosa keluhan, dan memilih part HGP / paket servis. Stok gudang otomatis berkurang.
2. **Tahap 2: Pengerjaan di Pit (`in_progress`)**: Mekanik pada Pit 1–4 mengerjakan motor sesuai instruksi pengerjaan.
3. **Tahap 3: Inspeksi Akhir QC (`completed`)**: Pemeriksaan keselamatan selesai, motor dipindahkan ke area penyerahan kasir.
4. **Tahap 4: Kasir POS Lunas (`paid`)**: Kasir memproses pembayaran (Tunai, QRIS dinamis, Debit, Transfer), mencetak struk thermal 58mm/80mm bergaransi 500 KM, atau mencetak faktur resmi A4.

### B. Alur Retail Suku Cadang Over-The-Counter (OTC)
Kasir langsung melayani pembelian oli (AHM Oil MPX/SPX) dan suku cadang fast-moving tanpa registrasi tiket servis. Dilengkapi pemindai barcode fisik, kode promo bundling (*NSSHEMAT10*), dan tombol pecahan uang cepat (*50k, 100k, 200k*).

### C. Alur Showroom Motor & Simulasi Kredit NSC Finance
1. Calon pembeli memilih model motor Honda (*Matic, Sport, EV*).
2. Tenaga penjual membuka **Kalkulator Simulasi Kredit NSC Finance** dan menggeser slider Uang Muka (10%–50%) serta tenor (11–35 bulan).
3. Sistem menghitung pokok hutang, bunga flat 18% p.a., asuransi TLO, biaya fidusia, cicilan bulanan, dan Total Bayar Pertama (TDP).
4. Tenaga penjual menerbitkan **SPK resmi** dengan alokasi fisik Nomor Rangka Astra (*VIN*) unit ready showroom.

---

## 5. Cara Menjalankan & Memverifikasi Sistem

```bash
# 1. Menjalankan Backend (Port 8080)
cd motoflow-backend
go run main.go

# 2. Menjalankan Frontend (Port 5173)
cd motoflow-frontend
npm run dev

# 3. Menjalankan Seluruh Unit & E2E Tests (Wajib 100% PASS)
cd motoflow-backend
go test -v ./...

# 4. Menjalankan Type Check & Frontend Production Build (Wajib 0 Error)
cd motoflow-frontend
npm run build
```
