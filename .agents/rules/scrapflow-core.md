# MotoFlow & ScrapFlow Core Rules
### Aturan Dasar Arsitektur & Prinsip Rekayasa Perangkat Lunak Bengkel AHASS NSS

Dokumen ini merupakan aturan wajib (*engineering invariants*) yang harus dipatuhi oleh seluruh developer dan AI coding agent saat memodifikasi atau mengembangkan fitur pada proyek **MotoFlow Workshop System** (adaptasi domain AHASS dari basis arsitektur ScrapFlow).

---

## 1. Domain Business Invariants (Kaidah Inti Domain)

Sistem ini merepresentasikan operasional terpadu **PT Nusantara Sakti Group (NSS)**:
1. **Bengkel Resmi Honda AHASS**:
   * Alur Work Order wajib mengikuti 4 tahap linier: `Draft` (SA) &rarr; `In Progress` (Pit Mekanik) &rarr; `Completed` (QC Selesai) &rarr; `Paid` (Kasir POS).
   * Nomor WO wajib berformat: `WO-AHASS-YYYYMMDD-XXXX`.
2. **Retail Suku Cadang Over-The-Counter (HGP)**:
   * Penjualan suku cadang Honda Genuine Parts dan AHM Oil tanpa pendaftaran servis mekanik.
   * Nomor Faktur Retail: `INV-RET-YYYYMMDD-XXXX`.
3. **Showroom Motor & NSC Multifinance**:
   * Setiap unit motor fisik wajib memiliki Nomor Rangka unik pabrikan Astra (*VIN*) dan Nomor Mesin.
   * Simulasi kredit cicilan NSC Finance tunduk pada aturan OJK: Uang Muka (DP) minimal 10% s/d 50%, Bunga flat 18.00% p.a., biaya fidusia Rp 450.000, asuransi TLO Rp 350.000.
   * Nomor SPK resmi: `SPK-NSS-YYYYMMDD-XXXX`.

---

## 2. Aturan Finansial & Decimal Arithmetic (Zero Floating-Point Drift)

> [!CRITICAL]
> **DILARANG KERAS** menggunakan tipe data `float32` atau `float64` untuk seluruh perhitungan uang, harga satuan, diskon, PPN 11%, dan kembalian kasir!

* **Backend (Golang)**:
  * Wajib menggunakan `github.com/shopspring/decimal`.
  * Semua konstanta moneter didefinisikan sebagai `decimal.Decimal` (misal: `decimal.NewFromFloat(0.11)`).
  * Pembulatan pajak dan diskon menggunakan metode pembulatan standar akuntansi `Round(0)` atau `RoundBanker()`.
* **Frontend (TypeScript)**:
  * Seluruh representasi moneter di-format melalui fungsi terpusat `formatRupiah()` di `src/utils/format.ts`.
  * Input nominal kasir menggunakan `type="number"` dengan validasi batas bawah nol (`min="0"`).

---

## 3. Aturan Konkurensi & Integritas Penomoran (Anti Race-Condition)

* Penomoran otomatis `WO-AHASS-...`, `SPK-NSS-...`, dan `INV-RET-...` **WAJIB** berada di dalam transaksi database atomik (`tx := db.Begin()`).
* Menggunakan teknik *Advisory Lock* atau penguncian baris (`SELECT ... FOR UPDATE`) untuk mencegah tabrakan nomor faktur yang sama jika kasir melakukan transaksi secara bersamaan pada detik yang sama (*concurrency safe*).

---

## 4. Aturan Snapshot Data Transaksi (Historical Immutability)

* Pada tabel `work_order_items` dan `retail_sale_items`, dilarang hanya menyimpan relasi `sparepart_id`.
* Wajib menyimpan nilai historis beku (*immutable snapshot*):
  * `price_snapshot` (Harga saat transaksi terjadi)
  * `part_number_snapshot` (Nomor part saat transaksi)
  * `item_name_snapshot` (Nama item saat transaksi)
* Perubahan master harga di masa depan tidak boleh mengubah nilai historis nota yang telah terbit.

---

## 5. Aturan Keamanan & Role-Based Access Control (RBAC)

Sistem menerapkan prinsip pertahanan berlapis ganda (*Defense in Depth*):
1. **Lapisan Backend (Golang Middleware)**:
   * Setiap endpoint terlindungi wajib divalidasi oleh `middleware.AuthMiddleware()`.
   * Pengecekan peran spesifik dilakukan dengan `middleware.RequireRole("admin", "service_advisor")`.
2. **Lapisan Frontend (Vue Router Navigation Guard)**:
   * Guard `router.beforeEach()` memeriksa token dan array `meta: { roles: [...] }`.
   * Menu navigasi (*Sidebar, Bottom Bar, Navbar*) difilter secara dinamis; peran yang tidak berhak **tidak boleh melihat tombol menu yang dilarang**.
   * Jika URL diakses langsung secara manual, arahkan ke rute default peran pengguna bersangkutan (`getDefaultRouteForRole()`) disertai toast peringatan.

---

## 6. Aturan Mesin Cetak (Isolated Print Engine)

> [!WARNING]
> **DILARANG** memanggil `window.print()` langsung pada halaman global karena akan mencetak navbar, sidebar, dan latar belakang dialog modal!

* Cetak struk kasir thermal (58mm/80mm) **WAJIB** menggunakan `printThermalReceipt()` di `src/utils/printer.ts`.
* Cetak faktur formal A4 **WAJIB** menggunakan `printA4Document()` di `src/utils/printer.ts` atau endpoint backend GoFPDF `/api/v1/work-orders/:id/invoice-pdf`.
* Mesin cetak harus memanfaatkan *isolated iframe* dengan aturan CSS `@page` khusus agar bersih dari header URL dan penomoran halaman bawaan peramban.

---

## 7. Standar Ergonomi Mobile-First

* Seluruh halaman harus nyaman digunakan pada smartphone layar kecil (320px–480px).
* Input teks wajib menggunakan `viewport-fit=cover, maximum-scale=1.0, user-scalable=no` di `index.html` untuk mematikan auto-zoom iOS Safari.
* Navigasi bawah wajib menyertakan `pb-[max(0.5rem,env(safe-area-inset-bottom))]`.
* Target sentuh tombol minimal 40–44px untuk kemudahan penekanan jari teknisi di lantai bengkel.
* Seluruh jendela modal harus bertransformasi menjadi *native bottom sheet* (`rounded-t-3xl`) pada layar `< 640px`.
