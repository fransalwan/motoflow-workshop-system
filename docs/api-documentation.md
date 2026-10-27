# Dokumentasi RESTful API MotoFlow AHASS
### PT Nusantara Sakti Group (NSS) - Sistem Informasi Manajemen Bengkel, Kasir POS & Showroom

Dokumen ini merupakan panduan teknis lengkap seluruh antarmuka pemrograman aplikasi (API) yang disediakan oleh backend Go (**Gin Framework**) pada port `http://localhost:8080` atau live tunnel gateway `/api/v1/*`.

---

## 1. Konvensi Standar API & Format Payload

Seluruh endpoint mengembalikan respons terstandarisasi dalam format JSON:

### Format Respons Sukses
```json
{
  "success": true,
  "message": "Deskripsi aksi berhasil",
  "data": { ... }
}
```

### Format Respons Galat (Error)
```json
{
  "success": false,
  "message": "Pesan deskripsi kesalahan ramah pengguna",
  "error": "KODE_ERROR_ATAU_DETAIL_VALIDASI"
}
```

### Header Autentikasi
Kecuali rute publik (`/api/v1/auth/login`), seluruh permintaan wajib menyertakan token JWT:
```http
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json
```

---

## 2. Modul Autentikasi & Akun Demo

### `POST /api/v1/auth/login`
Mengotentikasi pengguna dan mengembalikan token JWT bertanda tangan HMAC-SHA256 (masa aktif 24 jam).

* **Header**: `Content-Type: application/json`
* **Request Body**:
  ```json
  {
    "username": "kasir",
    "password": "password123"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Login berhasil",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsIn...",
      "user": {
        "id": 3,
        "username": "kasir",
        "full_name": "Siti Rahma (Kasir)",
        "role": "cashier"
      }
    }
  }
  ```

---

## 3. Modul Work Order & Servis Bengkel AHASS

### `GET /api/v1/work-orders`
Mengambil daftar antrean Work Order yang dapat difilter berdasarkan status.

* **Query Parameters**:
  * `status` (opsional): `draft`, `in_progress`, `completed`, `paid`
  * `search` (opsional): Nomor WO, nomor polisi plat motor, nama pemilik
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "data": [
      {
        "id": 1,
        "wo_number": "WO-AHASS-20261027-0001",
        "status": "in_progress",
        "payment_status": "unpaid",
        "entry_km": 12500,
        "grand_total": "189810.00",
        "customer": {
          "id": 1,
          "name": "Budi Santoso",
          "plate_number": "B 3456 KRF",
          "motor_model": "Honda Vario 160 ABS"
        },
        "mechanic": {
          "id": 1,
          "name": "Agus Setiawan",
          "pit_number": 1
        }
      }
    ]
  }
  ```

### `POST /api/v1/work-orders`
Menerbitkan Work Order baru oleh Service Advisor. Stok suku cadang HGP yang dipilih otomatis dipotong dari gudang.

* **Request Body**:
  ```json
  {
    "customer_id": 1,
    "mechanic_id": 1,
    "entry_km": 14500,
    "complaints": "Tarikan awal bergetar dan rem depan berdecit",
    "service_diagnosis": "Ganti kampas ganda CVT & kampas rem depan HGP",
    "items": [
      {
        "sparepart_id": 2,
        "quantity": 1,
        "discount_percent": 0
      },
      {
        "sparepart_id": 10,
        "quantity": 1,
        "discount_percent": 10
      }
    ]
  }
  ```

### `PATCH /api/v1/work-orders/:id/status`
Memperbarui status pengerjaan tiket servis.

* **Request Body**:
  ```json
  {
    "status": "completed"
  }
  ```

### `POST /api/v1/work-orders/:id/pay`
Penyelesaian kasir POS billing dengan kalkulasi kembalian berbasis tipe Decimal anti floating drift.

* **Request Body**:
  ```json
  {
    "payment_method": "cash",
    "paid_amount": 200000
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "message": "Pembayaran berhasil diproses",
    "data": {
      "wo_number": "WO-AHASS-20261027-0001",
      "payment_status": "paid",
      "grand_total": "189810.00",
      "paid_amount": "200000.00",
      "change_amount": "10190.00"
    }
  }
  ```

---

## 4. Modul Cetak Struk POS Thermal & Faktur A4 Resmi

### `GET /api/v1/work-orders/:id/receipt?width=58`
Menghasilkan teks struk kasir thermal berstandar ESC/POS untuk roll kertas 58mm atau 80mm.

* **Query Parameters**: `width=58` atau `width=80`
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "data": {
      "width_mm": 58,
      "raw_text": "================================\n  AHASS 01234 - NUSANTARA SAKTI\n================================\n...",
      "store_name": "AHASS 01234 - NSS",
      "warranty_notice": "Garansi servis 500 KM / 7 hari kerja."
    }
  }
  ```

### `GET /api/v1/work-orders/:id/invoice-pdf`
Menghasilkan dokumen faktur resmi formal A4 bertanda tangan ganda yang digenerate langsung oleh backend via library `jung-kurt/gofpdf`.

* **Content-Type**: `application/pdf`
* **Content-Disposition**: `inline; filename="Faktur-AHASS-WO-AHASS-XXXX.pdf"`

---

## 5. Modul Kasir Retail OTC (`/api/v1/retail/*`)

### `GET /api/v1/retail/promotions`
Mengambil daftar kupon promo aktif (*NSSHEMAT10*, *TUNEUPMATIC*, *OLIASLI5K*).

### `POST /api/v1/retail/checkout`
Checkout instan kasir meja Over-The-Counter tanpa alur mekanik pit.

* **Request Body**:
  ```json
  {
    "customer_name": "Walk-in Customer",
    "customer_phone": "08123456789",
    "payment_method": "cash",
    "paid_amount": 100000,
    "items": [
      {
        "sparepart_id": 1,
        "quantity": 1,
        "discount_percent": 0
      }
    ],
    "promo_code": "NSSHEMAT10"
  }
  ```

---

## 6. Modul Showroom Motor & Simulasi Kredit NSC Finance (`/api/v1/showroom/*`)

### `GET /api/v1/showroom/units`
Katalog motor showroom lengkap dengan Nomor Rangka Astra (*VIN*) dan Nomor Mesin.

### `POST /api/v1/showroom/simulate-credit`
Engine kalkulasi cicilan kredit NSC Finance:

* **Request Body**:
  ```json
  {
    "otr_price": 23500000,
    "dp_amount": 4700000,
    "tenor_months": 35
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "success": true,
    "data": {
      "otr_price": "23500000.00",
      "dp_amount": "4700000.00",
      "dp_percent": "20.00",
      "principal_loan": "18800000.00",
      "tenor_months": 35,
      "annual_interest_rate": "18.00",
      "total_interest": "9870000.00",
      "admin_fee": "450000.00",
      "insurance_tlo": "350000.00",
      "monthly_installment": "819143.00",
      "total_first_payment": "5500000.00"
    }
  }
  ```

### `POST /api/v1/showroom/spk`
Menerbitkan Surat Pemesanan Kendaraan (SPK) resmi PT Nusantara Sakti Group.
