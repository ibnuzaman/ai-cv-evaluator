# AI CV Evaluator

Aplikasi untuk mengevaluasi CV menggunakan AI dengan backend Go, PostgreSQL, ChromaDB, dan Google Gemini API.

## Setup Database

### 1. Konfigurasi Environment Variables

Buat file `.env` di root directory dengan mengcopy dari `.env.example`:

```bash
cp .env.example .env
```

Edit file `.env` dengan konfigurasi database Anda:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name
DB_SSLMODE=disable

# Connection Pool Settings (optional)
DB_MAX_OPEN=25
DB_MAX_IDLE=25
DB_CONNECTION_LIFETIME=5m
DB_CONNECTION_IDLE=5m

# Application Configuration
APP_PORT=8080

# API Keys
GEMINI_API_KEY=your_gemini_api_key

# ChromaDB Configuration
CHROMADB_URL=http://localhost:8000
```

### 2. Setup PostgreSQL Database

Pastikan PostgreSQL sudah terinstall dan running, kemudian buat database:

```sql
CREATE DATABASE your_database_name;
CREATE USER your_username WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE your_database_name TO your_username;
```

### 3. Migrasi Database

Aplikasi akan otomatis menjalankan migrasi database saat pertama kali dijalankan. File migrasi terletak di `database/migrations/`.

## Menjalankan Aplikasi

### 1. Start Services dengan Docker Compose

```bash
# Start PostgreSQL dan ChromaDB
docker-compose up -d
```

### 2. Seed ChromaDB dengan Data Konteks (Opsional)

```bash
# Build dan jalankan seeder
go build ./cmd/seed
./seed
```

### 3. Jalankan Aplikasi

```bash
# Build aplikasi
go build ./cmd/server

# Jalankan aplikasi
./server
```

Atau langsung jalankan dengan:

```bash
go run ./cmd/server/main.go
```

## Fitur AI Pipeline

- **File Reading**: Mendukung parsing PDF dan file teks untuk ekstraksi konten
- **ChromaDB Integration**: Retrieval-Augmented Generation (RAG) untuk konteks evaluasi
- **Gemini AI Integration**: Multi-stage evaluation menggunakan Google Gemini API
- **Asynchronous Processing**: Evaluasi berjalan di background dengan status tracking
- **Structured Output**: Hasil evaluasi dalam format JSON terstruktur

## Fitur Konfigurasi Database

- **Auto Migration**: Aplikasi otomatis menjalankan migrasi database saat startup
- **Connection Pool**: Konfigurasi connection pool untuk performa optimal
- **Environment Variables**: Semua konfigurasi database menggunakan environment variables
- **SSL Support**: Mendukung konfigurasi SSL untuk koneksi database
- **Graceful Connection**: Koneksi database dengan error handling yang baik

## API Endpoints

- `POST /api/v1/evaluate` - Evaluate CV
- `GET /api/v1/result/:id` - Get evaluation result

## Struktur Database

Tabel `evaluations`:
- `id` (SERIAL PRIMARY KEY)
- `content` (TEXT) - Konten CV yang dievaluasi  
- `result` (TEXT) - Hasil evaluasi
- `score` (INTEGER) - Skor evaluasi
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)
