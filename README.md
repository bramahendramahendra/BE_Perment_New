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
│   ├── master_sumber/
│   ├── pencapaian_kpi/
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

## Redis

Aplikasi menggunakan Redis sebagai cache. Pastikan Redis sudah berjalan sebelum menjalankan server.

### Cek Status Redis

```bash
sudo systemctl status redis-server
```

### Start Redis

```bash
sudo systemctl start redis-server
```

### Verifikasi Redis Berjalan

```bash
redis-cli ping
# Output: PONG
```

---

### Troubleshooting Redis

#### Error: `stop-writes-on-bgsave-error`

```
failed to initialize redis: MISCONF Redis is configured to save RDB snapshots,
but it is currently not able to persist on disk.
```

Penyebab: Redis gagal menyimpan snapshot RDB ke disk karena **permission denied** pada folder `/var/lib/redis`.

**Solusi — perbaiki permission:**

```bash
sudo chown -R redis:redis /var/lib/redis
sudo chmod 770 /var/lib/redis
sudo systemctl restart redis-server
```

Verifikasi berhasil:

```bash
redis-cli ping
# Output: PONG
```

**Solusi alternatif (development only)** — nonaktifkan pengecekan RDB error tanpa restart:

```bash
redis-cli config set stop-writes-on-bgsave-error no
```

> Solusi alternatif ini hanya untuk environment development. Untuk production, gunakan solusi perbaikan permission di atas.

#### Cek Log Redis

Jika masalah masih terjadi, lihat detail error di log Redis:

```bash
sudo tail -50 /var/log/redis/redis-server.log
```

---

## Deployment

Repository: `https://bitbucket.bri.co.id/scm/pt/perment-api-v2.git`

### Branch Convention

Kedua environment (development dan production) menggunakan branch yang sama: **`main`**.

Perbedaan environment ditentukan oleh nilai `RELEASE_MODE` di file `.env` pada masing-masing server.

| Server      | Branch | `RELEASE_MODE` | File Config              |
|-------------|--------|----------------|--------------------------|
| Development | `main` | `dev`          | `config/config_dev.json` |
| Production  | `main` | `prod`         | `config/config_prod.json`|

---

### Langkah Push ke Bitbucket

1. **Pastikan berada di branch `main`**

   ```bash
   git checkout main
   git pull origin main
   ```

2. **Merge perubahan dari branch fitur**

   ```bash
   git merge feature/<nama-fitur>
   ```

3. **Push ke Bitbucket**

   ```bash
   git push origin main
   ```

---

### Deploy ke Server Development

#### Pertama kali (clone)

1. Clone repository ke server development:

   ```bash
   git clone https://bitbucket.bri.co.id/scm/pt/perment-api-v2.git
   cd perment-api-v2
   ```

2. Buat file `.env` di root project:

   ```env
   APP_NAME=permen_api
   APP_AUTHOR=BRI
   APP_VERSION=1.0.0
   APP_HOST=0.0.0.0
   APP_PORT=8006
   RELEASE_MODE=dev
   ```

3. Pastikan file `config/config_dev.json` sudah tersedia dan konfigurasinya sesuai server development (database, minio, dll).

4. Build dan jalankan:

   ```bash
   go build -o perment-api main.go
   ./perment-api
   ```

   Atau jika menggunakan process manager:
   ```bash
   systemctl restart perment-api
   ```

#### Update berikutnya

```bash
git pull origin main
go build -o perment-api main.go
systemctl restart perment-api
```

---

### Deploy ke Server Production

> **Perhatian:** Pastikan semua perubahan sudah diuji di server development sebelum deploy ke production.

#### Pertama kali (clone)

1. Clone repository ke server production:

   ```bash
   git clone https://bitbucket.bri.co.id/scm/pt/perment-api-v2.git
   cd perment-api-v2
   ```

2. Buat file `.env` di root project:

   ```env
   APP_NAME=permen_api
   APP_AUTHOR=BRI
   APP_VERSION=1.0.0
   APP_HOST=0.0.0.0
   APP_PORT=8006
   RELEASE_MODE=prod
   ```

3. Pastikan file `config/config_prod.json` sudah tersedia dan konfigurasinya sesuai server production (database, minio, dll).

4. Build dan jalankan:

   ```bash
   go build -o perment-api main.go
   ./perment-api
   ```

   Atau jika menggunakan process manager:
   ```bash
   systemctl restart perment-api
   ```

#### Update berikutnya

```bash
git pull origin main
go build -o perment-api main.go
systemctl restart perment-api
```

---

### Konfigurasi systemd (Process Manager)

Buat file service agar aplikasi otomatis berjalan saat server reboot dan auto-restart jika crash.

```bash
sudo nano /etc/systemd/system/perment-api.service
```

Isi file:

```ini
[Unit]
Description=Perment API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/perment-api-v2
ExecStart=/path/to/perment-api-v2/perment-api
Restart=on-failure
RestartSec=5s
EnvironmentFile=/path/to/perment-api-v2/.env

[Install]
WantedBy=multi-user.target
```

Aktifkan dan jalankan service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable perment-api
sudo systemctl start perment-api

# Cek status
sudo systemctl status perment-api
```

---

### Konfigurasi Nginx

Nginx berfungsi sebagai **reverse proxy** yang meneruskan request dari luar ke aplikasi Go yang berjalan di port `8006`.

#### Install Nginx (jika belum ada)

```bash
sudo apt install nginx        # Ubuntu/Debian
# atau
sudo yum install nginx        # CentOS/RHEL
```

#### Konfigurasi Server Development

Buat file konfigurasi Nginx:

```bash
sudo nano /etc/nginx/conf.d/perment-api-dev.conf
```

Isi file:

```nginx
server {
    listen 80;
    server_name dev-perment-api.bri.co.id;  # sesuaikan dengan domain/IP server dev

    # Ukuran maksimal upload (untuk upload file Excel)
    client_max_body_size 20M;

    # Timeout untuk request yang lama (generate PDF/Excel)
    proxy_read_timeout 120s;
    proxy_connect_timeout 10s;
    proxy_send_timeout 120s;

    location /api {
        proxy_pass         http://127.0.0.1:8006;
        proxy_http_version 1.1;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}
```

#### Konfigurasi Server Production

```bash
sudo nano /etc/nginx/conf.d/perment-api-prod.conf
```

Isi file:

```nginx
server {
    listen 80;
    server_name perment-api.bri.co.id;  # sesuaikan dengan domain production

    # Redirect HTTP ke HTTPS
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name perment-api.bri.co.id;  # sesuaikan dengan domain production

    # Sertifikat SSL (sesuaikan path)
    ssl_certificate     /etc/ssl/certs/perment-api.crt;
    ssl_certificate_key /etc/ssl/private/perment-api.key;
    ssl_protocols       TLSv1.2 TLSv1.3;
    ssl_ciphers         HIGH:!aNULL:!MD5;

    # Ukuran maksimal upload (untuk upload file Excel)
    client_max_body_size 20M;

    # Timeout untuk request yang lama (generate PDF/Excel)
    proxy_read_timeout 120s;
    proxy_connect_timeout 10s;
    proxy_send_timeout 120s;

    location /api {
        proxy_pass         http://127.0.0.1:8006;
        proxy_http_version 1.1;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
    }
}
```

#### Aktifkan dan Reload Nginx

```bash
# Test konfigurasi (pastikan tidak ada error syntax)
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

### Catatan Deployment

- File `.env` dan `config/*.json` **tidak** di-commit ke repository — pastikan sudah tersedia di masing-masing server sebelum menjalankan aplikasi.
- Sesuaikan `server_name` di konfigurasi Nginx dengan domain atau IP yang digunakan di masing-masing server.
- Sesuaikan path sertifikat SSL di konfigurasi production dengan sertifikat yang diterbitkan oleh tim infrastruktur BRI.
- Cek log aplikasi setelah deploy untuk memastikan tidak ada error startup:
  ```bash
  journalctl -u perment-api -f
  ```
- Cek log Nginx jika ada masalah koneksi:
  ```bash
  sudo tail -f /var/log/nginx/error.log
  ```

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

### Pencapaian KPI *(Protected)*

| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | `/pencapaian-kpi/get-all-pencapaian` | Daftar semua pencapaian KPI |
| POST | `/pencapaian-kpi/get-detail` | Detail pencapaian KPI |
| POST | `/pencapaian-kpi/get-excel` | Download pencapaian KPI format Excel |
| POST | `/pencapaian-kpi/get-pdf` | Download pencapaian KPI format PDF |

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
| POST | `/master-sumber/get-all` | Daftar sumber KPI | ✅
| POST | `/user/get-all` | Daftar user | ✅