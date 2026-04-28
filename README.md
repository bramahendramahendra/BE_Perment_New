# BRIPerment Backend API

Backend service untuk aplikasi **Penyusunan KPI** BRI, dibangun menggunakan Go + Gin framework.

---

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** SQL (via `database` package)
- **Storage:** MinIO
- **Cache:** Redis *(opsional, nonaktif by default)*
- **Config:** Viper (`.env` + JSON config per environment)
- **Logger:** Zap

---

## Struktur Direktori

```
.
├── config/             # Konfigurasi aplikasi (env, database, redis, minio, dll)
├── domain/             # Domain logic (handler, service, repo, model, dto)
│   ├── audit_trail/
│   ├── auth/
│   ├── edm/
│   ├── master_context/
│   ├── master_divisi/
│   ├── master_kpi/
│   ├── master_perspektif/
│   ├── master_process/
│   ├── master_status/
│   ├── master_tahun/
│   ├── master_triwulan/
│   ├── penyusunan_kpi/
│   ├── realisasi_kpi/
│   ├── sample/
│   ├── template/
│   ├── user/
│   └── validasi_kpi/
├── dto/                # Shared DTO (response, error, filter, log)
├── errors/             # Custom error types
├── helper/             # Utility functions (ID generator, status code, dll)
├── middleware/         # Middleware (auth, CORS, audit trail, logging, error handler)
├── model/              # Shared model
├── pkg/                # Package eksternal (database, logger, minio, redis, transport)
├── repository/         # Shared repository (log request)
├── routes/             # Definisi routing
│   └── segment/        # Route per domain
├── server/             # Bootstrap & inisialisasi server
├── validation/         # Custom validator
└── main.go
```

---

## Konfigurasi

File konfigurasi dipilih berdasarkan nilai `RELEASE_MODE` di `.env`:

| RELEASE_MODE | File Config                  |
|-------------|-------------------------------|
| `local`     | `config/config_local.json`    |
| `dev`       | `config/config_dev.json`      |
| `uat`       | `config/config_uat.json`      |
| `prod`      | `config/config_prod.json`     |
| `bors`      | `config/config_bors_kost.json`|

Contoh `.env`:
```env
APP_NAME=permen_api
APP_AUTHOR=BRI
APP_VERSION=1.0.0
APP_HOST=localhost
APP_PORT=8006
RELEASE_MODE=local
```

---

## Menjalankan Server

```bash
go run main.go
```

Server berjalan di port yang dikonfigurasi (default: `8006`).

---

## API Endpoints

Base URL: `http://localhost:8006/api`

### Auth (Public)

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST   | `/auth`  | Mendapatkan Bearer Token |

---

### Penyusunan KPI *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/penyusunan-kpi/validate` | Validasi file Excel KPI | ✅
| POST | `/penyusunan-kpi/create` | Simpan data KPI | ✅
| POST | `/penyusunan-kpi/revision` | Revisi KPI (upload ulang) | ✅
| POST | `/penyusunan-kpi/approve` | Approval disetujui KPI | ✅
| POST | `/penyusunan-kpi/reject` | Approval ditolak KPI | ✅
| POST | `/penyusunan-kpi/get-all-approval` | Daftar KPI menunggu approval | ✅
| POST | `/penyusunan-kpi/get-all-tolakan` | Daftar KPI yang ditolak | ✅
| POST | `/penyusunan-kpi/get-all-daftar-penyusunan` | Daftar penyusunan KPI | ✅
| POST | `/penyusunan-kpi/get-all-daftar-approval` | Daftar approval KPI | ✅
| POST | `/penyusunan-kpi/get-detail` | Detail KPI | ✅
| POST | `/penyusunan-kpi/get-excel` | Download KPI format Excel | ✅
| POST | `/penyusunan-kpi/get-pdf` | Download KPI format PDF | ✅

---

### Realisasi KPI *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/realisasi-kpi/validate` | Validasi file Excel realisasi KPI |
| POST | `/realisasi-kpi/create` | Simpan data realisasi KPI |
| POST | `/realisasi-kpi/revision` | Revisi realisasi KPI (upload ulang) |
| POST | `/realisasi-kpi/approve` | Approval disetujui realisasi KPI |
| POST | `/realisasi-kpi/reject` | Approval ditolak realisasi KPI |
| POST | `/realisasi-kpi/get-all` | Daftar semua realisasi KPI |
| POST | `/realisasi-kpi/get-all-approval` | Daftar realisasi KPI menunggu approval |
| POST | `/realisasi-kpi/get-all-tolakan` | Daftar realisasi KPI yang ditolak |
| POST | `/realisasi-kpi/get-all-daftar-realisasi` | Daftar realisasi KPI |
| POST | `/realisasi-kpi/get-all-daftar-approval` | Daftar approval realisasi KPI |
| POST | `/realisasi-kpi/get-detail` | Detail realisasi KPI |

---

### Validasi KPI *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/validasi-kpi/input` | Input validasi KPI |
| POST | `/validasi-kpi/approval` | Kirim validasi ke approval |
| POST | `/validasi-kpi/approve` | Approve validasi KPI |
| POST | `/validasi-kpi/reject` | Reject validasi KPI |
| POST | `/validasi-kpi/batal` | Batalkan validasi KPI |
| POST | `/validasi-kpi/get-all-approval` | Daftar validasi menunggu approval |
| POST | `/validasi-kpi/get-all-tolakan` | Daftar validasi yang ditolak |
| POST | `/validasi-kpi/get-all-daftar-penyusunan` | Daftar penyusunan validasi |
| POST | `/validasi-kpi/get-all-daftar-approval` | Daftar approval validasi |
| POST | `/validasi-kpi/get-all-validasi` | Daftar semua validasi KPI |

---

### Template *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/template/format-penyusunan-kpi` | Download template penyusunan KPI | ✅
| POST | `/template/revision-penyusunan-kpi` | Download template revision KPI | ✅
| POST | `/template/format-realisasi-kpi` | Download template realisasi KPI | ✅

---

### EDM *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/edm/realisasi` | Ambil data realisasi dari EDM | 

---

### Master Data *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/master-triwulan/get-all` | Daftar triwulan | ✅
| POST | `/master-perspektif/get-all` | Daftar perspektif | ✅
| POST | `/master-tahun/get-all` | Daftar tahun | ✅
| POST | `/master-divisi/get-all` | Daftar divisi | ✅
| POST | `/master-kpi/get-all` | Daftar master KPI | ✅
| POST | `/master-status/get-all` | Semua status | ✅
| POST | `/master-process/get-all` | Daftar process | ✅
| POST | `/master-context/get-all` | Daftar context | ✅
| POST | `/user/get-all` | Daftar user | ✅