# MotoFlow AHASS Workshop System
### Sistem Informasi Manajemen Bengkel Resmi Honda AHASS, Kasir POS Billing, dan Inventaris Suku Cadang

[![Go Version](https://img.shields.io/badge/Golang-1.24%2B-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Vue 3](https://img.shields.io/badge/Vue.js-3.4%2B-4FC08D?style=for-the-badge&logo=vue.js)](https://vuejs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4-38B2AC?style=for-the-badge&logo=tailwind-css)](https://tailwindcss.com)
[![Docker](https://img.shields.io/badge/Docker-Multi--Stage-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com)
[![CI/CD](https://img.shields.io/badge/GitHub_Actions-Automated_CI-2088FF?style=for-the-badge&logo=github-actions)](.github/workflows/ci.yml)
[![Live Demo Tunnel](https://img.shields.io/badge/Live_Tunnel-Cloudflare-F38020?style=for-the-badge&logo=cloudflare)](https://online-thousands-tutorial-king.trycloudflare.com)
[![Mobile Ready](https://img.shields.io/badge/Mobile-First%20Responsive-brightgreen?style=for-the-badge)](https://online-thousands-tutorial-king.trycloudflare.com)

> **Portofolio Aplikasi Web Khusus Posisi "Programmer Junior" di PT Nusantara Sakti Group (NSS)**  
> *Dealer Resmi Honda, Bengkel AHASS, dan Perusahaan Multifinance Terbesar di Indonesia.*  
> 🌐 **Live Public Demo**: [https://online-thousands-tutorial-king.trycloudflare.com](https://online-thousands-tutorial-king.trycloudflare.com)  
> 📺 **Live TV Kiosk Display**: [https://online-thousands-tutorial-king.trycloudflare.com/waiting-display](https://online-thousands-tutorial-king.trycloudflare.com/waiting-display)

---

## Quick Navigation / Panduan Demo Penguji (HRD & User Teknis NSS)

| Modul Fitur | URL / Rute Demo | Fitur Utama yang Dapat Diuji |
| :--- | :--- | :--- |
| 🔐 **Portal Login Demo** | [`/login`](https://online-thousands-tutorial-king.trycloudflare.com/login) | Tombol 1-klik untuk berganti role: `Admin`, `Advisor`, `Kasir`, `Mekanik` |
| 🛒 **POS Retail Counter OTC** | [`/pos-retail`](https://online-thousands-tutorial-king.trycloudflare.com/pos-retail) | Penjualan langsung Oli/Part tanpa tiket servis, Barcode scan, Promo Bundling, Struk POS |
| 🏍️ **Showroom Motor & NSC** | [`/showroom`](https://online-thousands-tutorial-king.trycloudflare.com/showroom) | Katalog unit fisik, alokasi VIN/No. Rangka, Kalkulator Simulasi Kredit NSC Finance, Terbit SPK |
| 📊 **Executive Dashboard** | [`/dashboard`](https://online-thousands-tutorial-king.trycloudflare.com/dashboard) | Omset Harian/MTD, Top 5 Part HGP, Efisiensi Mekanik, Alert Stok Menipis |
| 📋 **Work Order & POS Cashier** | [`/work-orders`](https://online-thousands-tutorial-king.trycloudflare.com/work-orders) | Switcher **Tabel vs Kanban Pit Board**, Hotkeys POS (`F4`, `F8`, `F9`), QRIS Dinamis |
| 📺 **TV Waiting Lounge Monitor** | [`/waiting-display`](https://online-thousands-tutorial-king.trycloudflare.com/waiting-display) | Tampilan TV Kiosk 16:9 fullscreen, auto-polling 5s status Pit 1-4 & Marquee |
| 📦 **Katalog Honda Genuine Parts** | [`/inventory`](https://online-thousands-tutorial-king.trycloudflare.com/inventory) | Penyesuaian stok HGP, filter oli MPX/SPX, kampas rem, busi, dan v-belt |
| 👥 **Direktori Pelanggan & Motor** | [`/customers`](https://online-thousands-tutorial-king.trycloudflare.com/customers) | Master no. plat polisi, model motor Honda (Beat, Vario, PCX), riwayat KM |

---

## 1. Ringkasan & Kualifikasi Lowongan PT Nusantara Sakti Group

Proyek **motoflow-workshop-system** dirancang secara khusus untuk membuktikan kesiapan teknis dan penguasaan arsitektur terhadap kualifikasi posisi Programmer Junior di PT Nusantara Sakti Group:

1. **Database Relasional PostgreSQL Berstandar Industri**:
   - DDL skema relasi lengkap dengan *Foreign Keys* berintegritas tinggi.
   - Kebijakan *Cascade Delete* terisolasi pada `work_order_items`.
   - *Advisory Lock* / transactional locking anti *race condition* untuk penomoran otomatis `WO-AHASS-YYYYMMDD-XXXX`.
   - Showcase query agregasi kompleks (`docs/reporting-queries.sql`): `JOIN`, `GROUP BY`, `SUM`, `COUNT`, `AVG`, dan *Window Functions*.
2. **Modern Backend (Golang & Gin Framework)**:
   - Perhitungan harga, diskon, dan PPN 11% bebas *floating-point drift* via library `shopspring/decimal`.
   - Arsitektur berlapis (*Controller-Service-DTO-Model*) dengan *Table-Driven Unit Tests*.
   - Mesin cetak struk kasir thermal (58mm/80mm ESC/POS) dan faktur resmi A4 via `jung-kurt/gofpdf`.
   - Autentikasi JWT dan Role-Based Access Control (Admin, Service Advisor, Kasir, Mekanik).
3. **Frontend Reaktif (Vue 3, TypeScript, Vite & Tailwind CSS)**:
   - Komponen SPA modern (*Single Page Application*) dengan *Composition API* (`<script setup>`).
   - Pinia state management terstruktur dan modular.
   - Desain antarmuka berstandar visual resmi Honda AHASS Racing Red (`#CC0000`).

---

## 2. Matriks Adaptasi Arsitektur ScrapFlow ke Domain AHASS

| Komponen Asal (ScrapFlow) | Domain Bengkel Resmi Honda AHASS (MotoFlow) | Implementasi Teknis |
| :--- | :--- | :--- |
| **Master Customer** | Pelanggan & Kendaraan Honda | Nama Pemilik, No. HP, No. Plat (`B 3456 KRF`), Model Motor (Vario 160, Beat, PCX), Odometer KM |
| **Master Items** | Honda Genuine Parts (HGP) & Jasa Servis | AHM Oil MPX1/MPX2/SPX2, Kampas Rem Cakram/Tromol, V-Belt, Busi NGK, Jasa Servis Lengkap |
| **Scale Details** | Work Order Items (Rincian Servis) | `price_snapshot`, kuantitas suku cadang, diskon per item, subtotal bebas floating drift |
| **Invoices** | Work Order & Service Invoices | Status alur servis: `Draft` &rarr; `In Progress` &rarr; `Completed` &rarr; `Paid` |
| **Fleet & Driver** | Mekanik & Service Advisor | Alokasi Pit Servis (Pit 1-4), spesialisasi injeksi/CVT, dan Service Advisor penerima motor |
| **Print Engine** | POS Thermal Slip & Faktur A4 | Struk kasir thermal (58mm/80mm) & dokumen formal faktur A4 tanda tangan ganda via GoFPDF |
| **Decimal Engine** | Financial Arithmetic Engine | Perhitungan PPN 11%, diskon, dan kembalian via library `github.com/shopspring/decimal` |
| **Protection** | Concurrency Lock | Penomoran SPK anti-tabrakan: `WO-AHASS-YYYYMMDD-XXXX` |

---

## 3. Struktur Monorepo Proyek

```
motoflow-workshop-system/
├── motoflow-backend/                  # RESTful API Service (Golang)
│   ├── config/                        # Koneksi Database (PostgreSQL & SQLite Fallback)
│   ├── controllers/                   # Handler HTTP (Auth, WO, Sparepart, Dashboard, Print)
│   ├── dto/                           # Request & Response Payload Validation
│   ├── middleware/                    # Autentikasi JWT, Role Permission & CORS
│   ├── models/                        # Entitas GORM (User, Customer, Sparepart, WorkOrder, dll)
│   ├── routes/                        # Registrasi Route Terpusat Gin
│   ├── seeders/                       # Data Seed Realistis Motor Honda & Suku Cadang HGP
│   ├── services/                      # Logika Finansial (Decimal) & Table-Driven Unit Tests
│   ├── docker-compose.yml             # Container PostgreSQL Lokal
│   └── main.go                        # Entry Point Server
│
├── motoflow-frontend/                 # Web SPA (Vue 3 + Vite + TypeScript)
│   ├── src/
│   │   ├── components/                # ThermalReceiptModal, A4InvoiceModal, Sidebar, Navbar
│   │   ├── layouts/                   # MainLayout Navigasi Responsif
│   │   ├── pages/                     # Dashboard, WorkOrderList, CreateWO, POS Billing, Inventory
│   │   ├── router/                    # Guard Autentikasi & Rute Halaman
│   │   ├── stores/                    # Pinia Store (Auth, WorkOrder, Sparepart, Customer)
│   │   ├── types/                     # Interface TypeScript
│   │   └── utils/                     # Formatter Mata Uang Rupiah & Status Badge
│   ├── package.json
│   └── vite.config.ts
│
├── docs/                              # Technical Showcase untuk HRD / User Teknis NSS
│   ├── database-schema.sql            # Skema Relasi DDL, Foreign Key, & Index PostgreSQL
│   ├── reporting-queries.sql          # Showcase Query SQL: JOIN, GROUP BY, SUM, COUNT, AVG
│   └── system-flowchart.md            # Diagram Alur Bisnis Servis AHASS (Mermaid)
└── README.md                          # Dokumentasi Utama
```

### Diagram Topologi Arsitektur Produksi

```mermaid
graph TD
    Client["Peramban Klien (Kasir / SA / HRD NSS)"] --> Tunnel["Cloudflare Secure Tunnel<br/>(HTTPS Public Gateway)"]
    Tunnel --> Frontend["Frontend: Vue 3 SPA + Nginx Alpine<br/>(Port 80 / Reverse Proxy)"]
    Frontend -- "API Proxy: /api/v1/*" --> Backend["Backend: Golang Gin Framework<br/>(Decimal Engine & Concurrency Lock)"]
    Backend --> Cache["Caching Layer: In-Memory TTL & Redis<br/>(Dashboard Metrics & Fast-Moving Parts)"]
    Backend --> DB[("Database: PostgreSQL 16 Enterprise<br/>(DDL, Foreign Keys, ACID & Advisory Locks)")]
    Backend --> PDF["GoFPDF Engine<br/>(Faktur Resmi A4 & Struk Thermal POS)"]
```

---

## 4. Alur Proses Bisnis Servis AHASS

```mermaid
flowchart LR
    A["Motor Masuk ke AHASS"] --> B["Pendaftaran Service Advisor"]
    B --> C["Diagnosis & Terbitkan WO (Draft)"]
    C --> D["Mekanik Kerjakan di Pit (In Progress)"]
    D --> E["Inspeksi Akhir & Safety Check (Completed)"]
    E --> F["Kasir POS & Cetak Struk (Paid)"]
    F --> G["Penyerahan Unit & Garansi 500 KM"]
```

> Diagram alur lengkap, Sequence Diagram interaksi, dan Entity-Relationship Diagram (ERD) tersedia di [docs/system-flowchart.md](docs/system-flowchart.md).

---

## 5. Akun Pengguna Demo (Akses Cepat 1-Click)

Tersedia tombol login demo otomatis pada halaman `/login`:

| Peran (Role) | Username | Password | Tanggung Jawab Operasional | Halaman Utama (Default Home) |
| :--- | :--- | :--- | :--- | :--- |
| **Kepala Bengkel (Admin)** | `admin` | `password123` | Akses penuh dashboard eksekutif, data master, dan laporan omzet | `/dashboard` |
| **Service Advisor** | `advisor` | `password123` | Pendaftaran motor masuk, diagnosa keluhan, dan penerbitan SPK | `/work-orders` |
| **Kasir POS AHASS** | `kasir` | `password123` | Penyelesaian billing servis, kasir retail OTC, cetak struk thermal | `/pos-retail` |
| **Mekanik Senior** | `mekanik` | `password123` | Pelaksanaan pengerjaan servis di Pit dan pemantauan Kanban Board | `/work-orders` |

### Matriks Akses Halaman Berdasarkan Peran (Role-Based Access Control / RBAC)

| Halaman / Fitur Bengkel | Rute URL | Admin (Kepala Bengkel) | Service Advisor | Kasir POS | Mekanik Pit |
| :--- | :--- | :---: | :---: | :---: | :---: |
| **Executive Dashboard** | `/dashboard` | ✅ | ❌ | ❌ | ❌ |
| **Work Orders & Kanban Pit** | `/work-orders` | ✅ | ✅ (Home) | ✅ | ✅ (Home) |
| **Detail WO & Kasir Billing** | `/work-orders/:id` | ✅ | ✅ | ✅ | ✅ |
| **Pendaftaran Servis Baru** | `/work-orders/create` | ✅ | ✅ | ❌ | ❌ |
| **Kasir POS Retail OTC** | `/pos-retail` | ✅ | ✅ | ✅ (Home) | ❌ |
| **Showroom & NSC Finance** | `/showroom` | ✅ | ✅ | ✅ | ❌ |
| **Katalog Suku Cadang HGP** | `/inventory` | ✅ | ✅ | ✅ | ❌ |
| **Data Pelanggan & Motor** | `/customers` | ✅ | ✅ | ❌ | ❌ |
| **Mekanik & Pit Servis** | `/mechanics` | ✅ | ✅ | ❌ | ✅ |
| **TV Display Ruang Tunggu** | `/waiting-display` | ✅ | ✅ | ✅ | ✅ |

> **Catatan Keamanan & Tampilan Adaptif**: Navigasi menu (*Sidebar desktop, drawer mobile, dan bottom navigation bar*) **hanya menampilkan menu yang relevan bagi peran tersebut**. Jika pengguna mencoba mengetikkan URL terlarang secara manual di peramban, *Navigation Guard* Vue Router akan memblokir akses, memunculkan notifikasi peringatan (*Akses Terbatas*), dan mengarahkan pengguna kembali ke halaman utama mereka secara otomatis.

---

## 6. Showcase Fitur Utama

### A. Executive Workshop Dashboard
- **Omzet Hari Ini & MTD**: Agregasi *real-time* omset kotor, PPN 11%, dan jumlah unit motor yang diservis.
- **Top 5 Honda Genuine Parts**: Peringkat suku cadang terlaris berdasarkan kontribusi omset (misal: AHM Oil MPX2, Kampas Rem Cakram, V-Belt Kit).
- **Produktivitas Mekanik**: Pelacakan jumlah motor yang selesai diservis per Pit mekanik dan total pendapatan jasa mekanik.
- **Peringatan Stok Menipis**: Peringatan dini otomatis ketika stok fisik suku cadang berada di bawah batas minimum (`stock <= min_stock`).

### B. Kasir POS Billing & Engine Cetak Ganda
- **Pencegahan Floating-Point Drift**: Nilai moneter dikalkulasi menggunakan `shopspring/decimal` tanpa galat pembulatan IEEE-754.
- **Cetak Struk POS Thermal (58mm / 80mm)**: Format struk ESC/POS lengkap dengan header bengkel resmi AHASS Nusantara Sakti, data no. polisi plat motor, rincian part/jasa, PPN 11%, dan syarat garansi servis 500 KM.
- **Faktur Resmi Servis A4 (GoFPDF)**: Berkas PDF resmi berstandar A4 yang dihasilkan langsung dari backend Go, memuat rincian pekerjaan, tanda tangan ganda (*Pelanggan* & *Service Advisor*), serta stempel verifikasi digital AHASS.

```text
========================================
       AHASS NUSANTARA SAKTI MOTOR
     Jl. Raya Ciledug No. 88, Tangerang
       Telp: (021) 7345-6789 / 0812-3456
========================================
No. WO   : WO-AHASS-20261027-0001
Tgl/Jam  : 27/10/2026 14:30 WIB
Kasir    : kasir (Siti Rahma)
Mekanik  : Agus Setiawan (Pit 1)
----------------------------------------
Pelanggan: Budi Santoso
No. Plat : B 3456 KRF
Motor    : Honda Vario 160 ABS (12.500 KM)
----------------------------------------
RINCIAN SUKU CADANG / JASA    QTY  SUBTOTAL
----------------------------------------
AHM Oil MPX2 0.8L (Matik)      1  Rp 56.000
Busi NGK CPR9EA-9              1  Rp 25.000
Paket Servis Lengkap Matik     1  Rp 90.000
----------------------------------------
Subtotal Part & Jasa           : Rp 171.000
Diskon                         :       Rp 0
PPN 11%                        :  Rp 18.810
========================================
TOTAL TAGIHAN                  : Rp 189.810
Metode Bayar                   :      TUNAI
Diterima                       : Rp 200.000
Kembalian                      :  Rp 10.190
========================================
  * GARANSI SERVIS RESMI AHASS 500 KM *
  Gunakan selalu Honda Genuine Parts &
     Oli Standar Pabrikan AHM Oil!
 Terima kasih atas kepercayaan Anda di
       PT NUSANTARA SAKTI GROUP
========================================
```

---

## 7. Panduan Menjalankan Sistem (Local Setup)

### Prasyarat
- Go 1.24+ terpasang di sistem.
- Node.js v20+ dan npm terpasang.
- (Opsional) PostgreSQL 14+ atau Docker. *Backend secara cerdas mendukung SQLite otomatis untuk pengujian cepat jika PostgreSQL lokal belum aktif.*

### Langkah 1: Menjalankan Backend (Golang)
```bash
cd motoflow-backend

# Jalankan unit tests verifikasi finansial
go test -v ./...

# Jalankan linter vet
go vet ./...

# Jalankan server API (Port 8080)
go run main.go
```
*Saat pertama kali dijalankan, backend otomatis melakukan Auto-Migration dan men-seed data realistis pelanggan, motor Honda, suku cadang HGP, dan mekanik.*

### Langkah 2: Menjalankan Frontend (Vue 3)
```bash
cd motoflow-frontend

# Pasang dependensi
npm install

# Jalankan server pengembangan Vite (Port 5173)
npm run dev

# Untuk pengujian build produksi
npm run build
```

### Langkah 3: Menjalankan Seluruh Stack via Docker Compose (Produksi)
Untuk menjalankan seluruh ekosistem (*PostgreSQL 16, Redis 7, Golang Backend, dan Nginx Frontend*) hanya dengan 1 perintah:
```bash
docker compose up --build -d
```
Aplikasi web akan tersedia di port `http://localhost` (Port 80) dan API Backend di `http://localhost:8080`.

### Langkah 4: Menjalankan Database Migrations (Opsional untuk Produksi)
```bash
# Menjalankan migrasi DDL formal
migrate -path migrations/ -database "postgres://postgres:password123@localhost:5432/motoflow_ahass?sslmode=disable" up
```

---

## 8. Ringkasan Endpoint REST API

| Metode | Endpoint | Deskripsi | Akses Role |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | Otentikasi pengguna & pembuatan JWT | Publik |
| `GET` | `/api/v1/dashboard/summary` | Metrik omzet, unit servis, & low-stock alert | Terautentikasi |
| `GET` | `/api/v1/work-orders` | Daftar Work Order / Faktur Servis | Terautentikasi |
| `POST` | `/api/v1/work-orders` | Terbitkan Work Order baru & kurangi stok HGP | Service Advisor / Admin |
| `PATCH` | `/api/v1/work-orders/:id/status` | Update status: `in_progress`, `completed` | Mekanik / SA / Admin |
| `POST` | `/api/v1/work-orders/:id/pay` | Penyelesaian pembayaran kasir POS | Kasir / Admin |
| `GET` | `/api/v1/work-orders/:id/receipt` | Format teks struk kasir thermal (58mm/80mm) | Kasir / All |
| `GET` | `/api/v1/work-orders/:id/invoice-pdf` | Unduh dokumen faktur resmi A4 via GoFPDF | Kasir / All |
| `GET` | `/api/v1/spareparts` | Katalog Honda Genuine Parts & status stok | Terautentikasi |
| `PATCH` | `/api/v1/spareparts/:id/stock` | Penyesuaian stok masuk / keluar gudang | Admin / Staff |
| `GET` | `/api/v1/customers` | Master data pelanggan & nomor plat motor | Terautentikasi |
| `GET` | `/api/v1/mechanics` | Kinerja mekanik dan status Pit AHASS | Terautentikasi |

---

## 9. Showcase Query Pelaporan SQL

Seluruh query SQL agregasi tersedia di berkas [docs/reporting-queries.sql](docs/reporting-queries.sql):
- **Query 1**: Laporan Omset Harian & Bulanan Bengkel (`SUM`, `COUNT`, `AVG`, `GROUP BY DATE`).
- **Query 2**: Top 5 Honda Genuine Parts Terlaris (`JOIN`, `GROUP BY`, `ORDER BY omset DESC`).
- **Query 3**: Kinerja & Produktivitas Mekanik AHASS (`LEFT JOIN`, `AVG waktu pengerjaan`).
- **Query 4**: Notifikasi Peringatan Suku Cadang Menipis (`WHERE stock <= min_stock`).
- **Query 5**: Analisa Komparasi Omset Suku Cadang vs Jasa Servis (Revenue Mix).

---

## 10. Implementasi Fase 4: Production Maturity, Caching, Containerization & CI/CD

Sebagai bukti kesiapan ke tingkat *production-ready enterprise*, sistem telah dilengkapi dengan arsitektur Fase 4:

1. **High-Performance In-Memory & Redis Caching Layer (`services/cache_service.go`)**:
   - Caching thread-safe dengan TTL pada ringkasan dashboard eksekutif dan katalog suku cadang fast-moving.
   - *Cache Invalidation otomatis*: Cache metrik langsung di-reset begitu transaksi Work Order baru masuk, status diperbarui, atau kasir menyelesaikan pembayaran.
2. **Database Migration Management Resmi (`migrations/`)**:
   - `000001_create_initial_schema.up.sql` & `.down.sql`: DDL relasional PostgreSQL berstandar enterprise.
   - `000002_add_indexes_and_constraints.up.sql` & `.down.sql`: Composite covering indexes untuk performa tinggi pencarian plat motor dan nomor part.
3. **Containerization Multi-Stage Docker**:
   - **Backend Dockerfile**: Menggunakan `golang:1.24-alpine` sebagai builder dan `alpine:3.19` minimalis non-root user sebagai runtime dengan zona waktu Asia/Jakarta.
   - **Frontend Dockerfile**: Multi-stage build Node 20 dengan server reverse-proxy Nginx Alpine teroptimasi gzip.
   - **Docker Compose Orkestrasi Penuh (`docker-compose.yml`)**: Menggabungkan PostgreSQL 16, Redis 7, Golang API Backend, dan Vue 3 SPA Frontend.
4. **Automated CI/CD Pipeline (`.github/workflows/ci.yml`)**:
   - Otomatis menjalankan `go mod verify`, `go vet ./...`, `go test -v -race ./...`, kompilasi Go binary, dan `npm run build` pada setiap *push* atau *pull request* ke branch `main`.
5. **Live Public Tunneling Deployment**:
   - Terintegrasi dengan secure tunneling Cloudflare sehingga dapat diakses secara publik dan aman via protokol HTTPS.

---

## 11. Mobile-First Ergonomics & Inovasi UI/UX Bengkel AHASS

Untuk menghadirkan pengalaman pengguna kelas industri yang siap diimplementasikan langsung di lantai bengkel (*workshop floor*) maupun ruang tunggu pelanggan resmi AHASS, sistem ini dirancang dengan standar **Mobile-First App Experience**:

### A. Konfigurasi Viewport & Safe-Area Inset Standar PWA
- **Pencegahan Auto-Zoom iOS Safari**: Pengaturan `viewport-fit=cover, maximum-scale=1.0, user-scalable=no` memastikan kolom input teks tidak memicu pembesaran layar (*viewport shift*) saat disentuh teknisi di perangkat iPhone.
- **Dukungan Bezel-less & Gesture Bar**: Penggunaan `env(safe-area-inset-bottom)` pada navigasi bawah dan kontainer utama memastikan tombol aksi tidak terpotong oleh garis *Home Indicator* iOS maupun bilah navigasi gestur Android.
- **PWA Meta Tags**: Dilengkapi `<meta name="theme-color" content="#CC0000" />` dan `<meta name="apple-mobile-web-app-capable" content="yes" />` dengan status bar beraksen transparan resmi Honda Racing Red.

### B. Sticky Floating Action Bar (Standar Moka POS, Grab & Tokopedia)
- **POS Retail OTC (`/pos-retail`)**:
  - **Mobile Segmented Switcher**: Membagi tampilan smartphone menjadi tab *Katalog Part* vs *Keranjang Belanja* untuk mencegah antrean barang menenggelamkan kolom kasir.
  - **Sticky Floating Cart Bar**: Bilah keranjang melayang di atas navigasi bawah saat berada di katalog: `🛒 [X Item • Rp Total] [Lihat Keranjang & Bayar ➔]`, sehingga kasir cukup 1 ketukan jempol untuk langsung menyelesaikan transaksi.
  - **Touch Stepper**: Tombol penambah dan pengurang kuantiti (*+ / -*) berukuran minimal 36px ramah jempol dengan umpan balik klik instan.
- **Work Order Detail & Kasir Billing (`/work-orders/:id`)**:
  - **Mobile Dual Tab**: Tab *Rincian Servis & Part* vs *Kasir POS & Unit* membuat kasir bengkel tidak perlu menggulir puluhan rincian baut dan oli untuk mengakses terminal pembayaran.
  - **Sticky Floating Pay Bar**: Bilah aksi melayang `💰 [Total Rp X] [Bayar Kasir ➔]` aktif otomatis saat status nota masih belum lunas.

### C. Kanban Pit Board dengan Mobile Column Switcher
- **Segmented Column Tabs**: Pada layar smartphone (`< md`), 4 kolom Kanban yang biasanya menumpuk vertikal panjang (*infinite scrolling*) dikonversi menjadi bar navigasi kolom dinamis:
  `[1. Antrean SA (X)]` `[2. Pit Servis (Y)]` `[3. Kasir POS (Z)]` `[4. Selesai Lunas (W)]`.
- Hanya menampilkan kolom yang sedang dipantau secara fokus dan responsif, sementara pada layar tablet dan desktop (`>= md`) tetap menyajikan grid 4 kolom horizontal penuh.

### D. Transformasi Dialog Modal Menjadi Native Mobile Bottom Sheets
- Seluruh jendela modal kompleks pada sistem bertransisi secara mulus menjadi *Bottom Sheet* modern di layar ponsel (`rounded-t-3xl`, *drag handle bar*, pembatasan tinggi `max-h-[92vh]` dan *scrollable body*):
  - **Kalkulator Kredit NSC Finance** & **Formulir SPK Resmi Showroom**.
  - **Formulir Pendaftaran Pelanggan & Motor Honda Baru**.
  - **Penyesuaian Stok Gudang HGP**.
  - **Pratinjau Struk Kasir Thermal (58mm/80mm)** & **Faktur A4 Resmi**.
- Menghilangkan bug terpotongnya tombol konfirmasi akibat keyboard sentuh (*virtual keyboard*) pada layar HP.

### E. Responsive Card View & Form Pembuatan WO
- **Kartu WO Interaktif**: Tabel antrean bertransformasi menjadi kartu berdensitas informasi tinggi pada layar kecil (`< 768px`) dengan status alur servis beranimasi.
- **Formulir Item Servis Ramah Sentuh**: Formulir pendaftaran servis mengubah tabel 5 kolom menjadi kartu item berundak dengan pemilih suku cadang, pengatur kuantitas, dan kalkulator estimasi total biaya servis otomatis.
- **Tren Kinerja Dashboard**: Grafik 7-hari dashboard eksekutif beradaptasi dari 7 kolom di desktop menjadi `grid-cols-2` di mobile untuk keterbacaan metrik pendapatan harian.

### F. Kasir POS Keyboard Shortcuts (Hotkeys)
Mendukung alur transaksi kasir super cepat di jam sibuk tanpa perlu menyentuh mouse:
- <kbd>F4</kbd>: Rotasi instan metode pembayaran (*TUNAI* &rarr; *QRIS* &rarr; *KARTU DEBIT* &rarr; *TRANSFER*).
- <kbd>F8</kbd>: Isi otomatis uang pas (*Exact Cash Amount*).
- <kbd>F9</kbd>: Tampilkan pratinjau struk thermal 58mm/80mm siap cetak.
- <kbd>ESC</kbd>: Menutup jendela modal aktif secara instan.

### G. Dynamic QRIS Modal & Web Audio Synthesizer
- **QRIS Dinamis**: Modal QRIS otomatis dengan *countdown timer* 15 menit dan simulasi verifikasi transfer QRIS 1-klik.
- **Native Web Audio API Synthesizer**: Efek suara sintetis tanpa aset berkas eksternal mp3:
  - *Barcode Scanner Beep* (880Hz Sine Wave) saat memilih suku cadang HGP.
  - *Success Chime* (Arpeggio nada ganda 523Hz & 659Hz) saat transaksi kasir berhasil diselesaikan.
  - *Click Feedback* (1200Hz) untuk umpan balik respons tombol sentuh.

### H. Pinia Toast Notification System
Sistem notifikasi modern *non-blocking* menggantikan fungsi `alert()` bawaan browser, dilengkapi transisi halus, pewarnaan kontekstual (*Success, Error, Info*), dan *auto-dismiss* 4 detik.

### I. Customer Waiting Lounge TV Monitor (`/waiting-display`)
- **Fullscreen Kiosk Mode 16:9**: Tampilan layar lebar khusus TV monitor ruang tunggu pelanggan AHASS.
- **Live Status Polling**: Pembaruan status otomatis setiap 5 detik memantau antrean motor di Pit 1 s/d Pit 4.
- **Marquee Running Ticker**: Teks berjalan di bagian bawah mengedukasi pelanggan mengenai keaslian oli Honda AHM Oil, garansi servis 500 KM, dan tips keselamatan berkendara *#Cari_Aman*.

---

## 12. Modul Bisnis Retail Otomotif & Integrasi Pembiayaan NSC Multifinance

Untuk mencerminkan secara utuh model bisnis **PT Nusantara Sakti Group (NSS)** sebagai Dealer Resmi Honda, Bengkel AHASS, dan Pembiayaan Kredit Terbesar (**NSC Finance / PT Nusa Surya Ciptadana**), sistem ini mengintegrasikan proses bisnis retail modern:

### A. Direct Over-The-Counter (OTC) POS Cashier (`/pos-retail`)
- **Fast-Throughput Checkout**: Transaksi kasir langsung untuk pembelian suku cadang Honda Genuine Parts (HGP), oli standar pabrikan AHM Oil (MPX1, MPX2, SPX2), dan aksesoris resmi tanpa harus melalui tiket servis antrean mekanik.
- **Dukungan Barcode Scanner**: Bidang pencarian otomatis siap menerima input pemindai *barcode* fisik, lengkap dengan efek suara sintetis *Scanner Beep* (880Hz Sine Wave).
- **Promo Engine & Bundling**:
  - Kupon diskon otomatis (*misal: `NSSHEMAT10`, `TUNEUPMATIC`, `OLIASLI5K`*).
  - Peringatan batas minimal belanja dan kalkulasi diskon flat/persentase secara *real-time*.
- **Quick-Cash Presets & Struk Kasir Thermal**: Tombol uang pas, tombol pecahan cepat (*50k, 100k, 200k*), perhitungan kembalian instan, dan pratinjau struk kasir thermal 58mm/80mm ESC/POS.

### B. Showroom Unit Kendaraan & Alokasi Nomor Rangka Astra (`/showroom`)
- **Katalog Fisik Motor Showroom**: Menampilkan jajaran motor Honda terbaru (*Honda Vario 160 ABS, Beat Deluxe CBS-ISS, PCX 160, Scoopy Prestige, EM1 e: Motor Listrik, CBR150R*).
- **Pelacakan Fisik VIN & Nomor Mesin**: Setiap unit motor memiliki Nomor Rangka unik pabrikan Astra (`MH1...`) dan Nomor Mesin (`KF81E...`) untuk integritas pengurusan faktur dan STNK Samsat.
- **Status Alokasi Unit**: Status transisi unit kendaraan dari *Ready Showroom* &rarr; *Booked SPK* &rarr; *Delivered*.

### C. Kalkulator Simulasi Kredit Interaktif NSC Finance
- **Interactive DP Slider**: Calon debitur dapat menggeser *slider* Uang Muka (10% s/d 50% sesuai regulasi OJK) dengan konversi nominal Rupiah langsung.
- **Tenor Fleksibel**: Pilihan jangka waktu cicilan 11, 17, 23, 29, dan 35 bulan.
- **Engine Finansial Cicilan**:
  - Perhitungan Pokok Hutang (*Principal Loan*).
  - Bunga pembiayaan flat 18.00% per tahun (1.5% per bulan).
  - Estimasi cicilan per bulan, biaya fidusia/administrasi, dan asuransi kehilangan Total Loss Only (TLO).
  - Total Pembayaran Pertama (TDP: DP + Biaya Admin + Asuransi).

### D. Penerbitan Surat Pemesanan Kendaraan (SPK) Resmi NSS
- Penerbitan nomor SPK otomatis anti-tabrakan dengan *advisory locking*: `SPK-NSS-YYYYMMDD-XXXX`.
- Pilihan skema pembayaran: **Cash Keras** (dengan bukti Tanda Jadi / Booking Fee) atau **Kredit NSC Finance** (lengkap dengan data tenor dan angsuran bulanan).
- Cetak lembar konfirmasi SPK resmi berstandar dealer Astra Honda.

---

## 13. Profil Pengembang & Lisensi
- **Pengembang**: **Frans Alwan** ([@fransalwan](https://github.com/fransalwan))
- **Email**: fransalwan55@gmail.com
- **Posisi Tujuan**: **Programmer Junior** — PT Nusantara Sakti Group (NSS)
- **Lisensi**: Proyek ini dilisensikan di bawah lisensi [MIT](LICENSE). Dibuat dengan dedikasi tinggi sebagai portofolio teknis komprehensif.
