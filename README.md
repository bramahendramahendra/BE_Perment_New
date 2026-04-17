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
│   ├── master_challenge/
│   ├── master_divisi/
│   ├── master_kpi/
│   ├── master_method/
│   ├── master_perspektif/
│   ├── master_status/
│   ├── master_tahun/
│   ├── master_triwulan/
│   ├── penyusunan_kpi/
│   ├── sample/
│   └── template/
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
| POST | `/penyusunan-kpi/revision` | Revisi KPI (upload ulang) |
| POST | `/penyusunan-kpi/create` | Simpan data KPI | ✅
| POST | `/penyusunan-kpi/approval` | Approval KPI |
| POST | `/penyusunan-kpi/batal` | Batalkan KPI |
| POST | `/penyusunan-kpi/get-all-approval` | Daftar KPI menunggu approval | ✅
| POST | `/penyusunan-kpi/get-all-tolakan` | Daftar KPI yang ditolak | ✅
| POST | `/penyusunan-kpi/get-all-daftar-penyusunan` | Daftar penyusunan KPI | ✅
| POST | `/penyusunan-kpi/get-all-daftar-approval` | Daftar approval KPI | ✅
| POST | `/penyusunan-kpi/get-detail` | Detail KPI |
| POST | `/penyusunan-kpi/get-excel` | Download KPI format Excel |
| POST | `/penyusunan-kpi/get-pdf` | Download KPI format PDF |

---

### Template *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/template/format-penyusunan-kpi` | Download template penyusunan KPI |
| POST | `/template/tolakan-penyusunan-kpi` | Download template tolakan KPI |
| POST | `/template/format-realisasi-kpi` | Download template realisasi KPI |

---

### Master Data *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/master-triwulan/get-all` | Daftar triwulan |
| POST | `/master-perspektif/get-all` | Daftar perspektif |
| POST | `/master-tahun/get-all` | Daftar tahun |
| POST | `/master-divisi/get-all` | Daftar divisi |
| POST | `/master-kpi/get-all` | Daftar master KPI |
| POST | `/master-status/get-all` | Semua status |
| POST | `/master-status/get-draft` | Status draft |
| POST | `/master-challenge/get-all` | Daftar challenge |
| POST | `/master-method/get-all` | Daftar method |
| POST | `/user/get-all` | Daftar user |
