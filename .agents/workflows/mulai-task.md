# Workflow: Standar Operasional Memulai Task (SOP Development)
### Panduan Langkah Demi Langkah Eksekusi Fitur Baru, Modifikasi, atau Bugfix

Workflow ini adalah prosedur standar yang **WAJIB** diikuti developer atau AI agent setiap kali menerima permintaan penugasan task baru pada sistem **MotoFlow**.

---

## Tahap 1: Klarifikasi Kebutuhan & Analisis Dampak Domain

Sebelum menulis satu baris kode pun, lakukan analisis:
1. **Analisis Peran Pengguna (RBAC Impact)**:
   * Siapa yang akan menggunakan fitur ini? (*Admin, Service Advisor, Kasir, atau Mekanik*).
   * Pastikan hak akses dibatasi pada level backend (*middleware*) dan frontend (*router meta & sidebar filtering*).
2. **Analisis Integritas Finansial**:
   * Apakah task ini menyentuh nominal uang, diskon, PPN, cicilan, atau kasir?
   * Jika ya, **WAJIB** gunakan `decimal.Decimal` di Go dan `formatRupiah()` di Vue. Dilarang menggunakan `float`.
3. **Analisis Skema Database**:
   * Apakah ada penambahan entitas atau kolom baru?
   * Pastikan kolom relasi memiliki *Foreign Key* yang tepat dan mendukung SQLite *in-memory* serta PostgreSQL produksi.

---

## Tahap 2: Eksekusi Backend (Golang Clean Architecture)

Ikuti alur berlapis (*layered flow*) backend:

```
1. models/       -> Definisikan struct entitas GORM dengan tag JSON dan gorm
2. dto/          -> Buat struct Request & Response dengan tag binding:"required"
3. services/     -> Tulis logika bisnis & kalkulasi moneter via decimal.Decimal
4. controllers/  -> Tangani parsing JSON, error handling, dan response JSON terstandar
5. routes/       -> Daftarkan endpoint ke grup /api/v1/* dengan middleware RBAC yang tepat
```

*Aturan Tambahan*:
* Jika membuat penomoran dokumen baru, terapkan *transactional locking* agar anti tabrakan.
* Jika mutasi stok terjadi, jalankan dalam blok transaksi database atomik.

---

## Tahap 3: Eksekusi Frontend (Vue 3 + TypeScript)

Ikuti alur terstruktur frontend:

```
1. types/        -> Definisikan interface TypeScript data payload
2. stores/       -> Buat / perbarui modul Pinia store dengan action pemanggilan API
3. components/   -> Buat komponen modular (Modal, Button, Card, Badge)
4. pages/        -> Rakit tampilan halaman reaktif dengan <script setup lang="ts">
5. router/       -> Daftarkan rute dengan meta: { requiresAuth: true, roles: [...] }
```

*Standar Wajib Ergonomi Mobile-First*:
* Pastikan target sentuh tombol minimal 40–44px (`py-2.5 px-4`).
* Gunakan *native bottom sheet* untuk jendela dialog modal di layar smartphone (`rounded-t-3xl sm:rounded-2xl`).
* Jika ada form tabel panjang di desktop, sediakan tampilan *Item Cards* bertingkat untuk layar ponsel (`< 640px`).
* Jika menyangkut kasir retail atau nota servis, gunakan *sticky floating action bar* di bagian bawah.
* Cetak struk wajib memanggil `printThermalReceipt()` atau `printA4Document()` dari `src/utils/printer.ts`.

---

## Tahap 4: Quality Gate & Verifikasi Otomatis

Sebelum menyelesaikan penugasan, **WAJIB** menjalankan gerbang pengujian mutu berikut:

```bash
# 1. Jalankan unit test dan E2E test backend Go (Wajib 100% PASS)
cd motoflow-backend
go test -v ./...

# 2. Jalankan type checking dan compilation frontend Vue 3 (Wajib Exit Code 0)
cd motoflow-frontend
npm run build
```

> [!CAUTION]
> Dilarang melakukan commit jika masih ada tes yang gagal (*FAIL*) atau galat type error TypeScript!

---

## Tahap 5: Dokumentasi & Protokol Git Commit

1. **Perbarui Dokumentasi**:
   * Jika ada endpoint baru, tambahkan ke [docs/api-documentation.md](file:///c:/Users/Admin/Desktop/portfolio-opensource/motoflow-workshop-system/docs/api-documentation.md).
   * Jika ada formula bisnis baru, catat di [docs/business-architecture.md](file:///c:/Users/Admin/Desktop/portfolio-opensource/motoflow-workshop-system/docs/business-architecture.md).
   * Perbarui ringkasan fitur di [README.md](file:///c:/Users/Admin/Desktop/portfolio-opensource/motoflow-workshop-system/README.md).
2. **Commit Sesuai Standar Conventional Commits**:
   * `feat(...)`: Penambahan fitur baru.
   * `fix(...)`: Perbaikan bug.
   * `docs(...)`: Pembaruan dokumentasi teknis.
   * `refactor(...)`: Perapian struktur kode tanpa mengubah fungsionalitas.
3. **Push ke GitHub**:
   * Jalankan `git push origin main` dan pastikan pipeline CI/CD GitHub Actions berjalan sukses (*Green*).
