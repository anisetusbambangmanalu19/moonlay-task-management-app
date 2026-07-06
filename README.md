# Moonlay Task Management App

Aplikasi manajemen task sederhana dengan fitur CRUD task, autentikasi JWT, dan AI Chatbot berbasis RAG (Retrieval-Augmented Generation) — dibuat sebagai technical test untuk posisi **Fullstack Developer Intern @ PT Moonlay Technologies**.

---

## Tech Stack

| Layer | Teknologi |
|---|---|
| **Frontend** | Next.js 15 (App Router) + TypeScript + Tailwind CSS |
| **Backend** | Golang 1.23 + Gin Framework |
| **ORM/Query** | GORM (CRUD standar) + Raw SQL via `db.Raw()` (chatbot) |
| **Database** | PostgreSQL |
| **Auth** | JWT (access token, expiry 24 jam) + bcrypt |
| **Chatbot** | RAG: PostgreSQL → Prompt → Gemini 1.5 Flash |

---

## Struktur Folder

```
moonlay-task-management-app/
├── backend/
│   ├── cmd/api/main.go           # Entry point
│   ├── internal/
│   │   ├── config/               # DB connection, env loader
│   │   ├── models/               # Struct User, Task
│   │   ├── repository/           # Query layer (GORM + raw SQL)
│   │   ├── handlers/             # HTTP handler per resource
│   │   ├── middleware/           # JWT auth middleware
│   │   └── routes/               # Route registration
│   ├── migrations/
│   │   ├── 001_init.sql          # Schema migration
│   │   └── seed/main.go          # Seed script (bcrypt hash at runtime)
│   ├── .env.example
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── login/page.tsx
│   │   │   ├── tasks/page.tsx
│   │   │   └── layout.tsx
│   │   ├── components/
│   │   │   ├── TaskCard.tsx
│   │   │   ├── TaskForm.tsx
│   │   │   ├── ChatbotWidget.tsx
│   │   │   └── StatusBadge.tsx
│   │   ├── lib/
│   │   │   ├── api.ts            # Axios + interceptors
│   │   │   └── auth.ts           # Token helpers
│   │   └── types/index.ts
│   └── .env.local
└── docs/
    ├── Moonlay_TaskApp.postman_collection.json
    ├── Moonlay_TaskApp.postman_environment.json
    └── ERD.md
```

---

## Prasyarat

| Tool | Versi Minimum | Link |
|---|---|---|
| **Go** | 1.23+ | https://go.dev/dl |
| **Node.js** | 18+ (LTS) | https://nodejs.org |
| **PostgreSQL** | 14+ | https://postgresql.org/download |

---

## Setup Database

### 1. Buat database baru

Lewat Navicat atau psql:
```sql
CREATE DATABASE moonlay_task_db;
```

### 2. Jalankan migration schema

Buka `backend/migrations/001_init.sql` di Navicat Query Tool (pilih database `moonlay_task_db`) dan **Run**.

Schema yang dibuat:
- ENUM `task_status` (`todo`, `in_progress`, `done`)
- Tabel `users` dan `tasks`
- Index pada `assignee_id`, `status`, `deadline`
- Trigger auto-update `updated_at`

---

## Setup Backend

### 1. Konfigurasi environment

```bash
cd backend
copy .env.example .env
```

Edit `.env` dengan nilai yang sesuai:

```env
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_postgres_kamu
DB_NAME=moonlay_task_db
DB_SSLMODE=disable

JWT_SECRET=minimal_32_karakter_random_string_yang_aman
JWT_EXPIRY_HOURS=24

GEMINI_API_KEY=isi_dari_aistudio_google_com
```

> **Mendapatkan Gemini API Key**: Buka [aistudio.google.com](https://aistudio.google.com) → Get API Key → Create API Key (gratis)

### 2. Seed database

Jalankan dari folder `backend/`:
```bash
go run migrations/seed/main.go
```

Script ini akan insert 4 user dengan password yang di-hash bcrypt secara otomatis. Aman untuk dijalankan berulang kali (idempotent via `FirstOrCreate`).

### 3. Jalankan backend

```bash
go run cmd/api/main.go
```

Backend berjalan di `http://localhost:8080`

---

## Setup Frontend

### 1. Konfigurasi environment

File `.env.local` sudah ada dengan nilai default:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api
```

### 2. Install dependencies

```bash
cd frontend
npm install
```

### 3. Jalankan frontend

```bash
npm run dev
```

Frontend berjalan di `http://localhost:3000`

---

## Kredensial Login (untuk Testing)

| Nama | Email | Password |
|---|---|---|
| Admin | admin@moonlay.com | admin123 |
| Budi Santoso | budi@moonlay.com | budi123 |
| Siti Rahayu | siti@moonlay.com | siti123 |
| Rangga Pratama | rangga@moonlay.com | rangga123 |

> Password disimpan sebagai bcrypt hash di database — tidak ada plaintext.

---

## API Endpoints

| Method | Path | Auth | Deskripsi |
|---|---|---|---|
| `POST` | `/api/auth/login` | ❌ | Login, return JWT token |
| `GET` | `/api/users` | ✅ | List semua user |
| `GET` | `/api/tasks` | ✅ | List semua task + assignee |
| `GET` | `/api/tasks/:id` | ✅ | Detail satu task |
| `POST` | `/api/tasks` | ✅ | Buat task baru |
| `PUT` | `/api/tasks/:id` | ✅ | Update semua field task |
| `PATCH` | `/api/tasks/:id/status` | ✅ | Update status saja |
| `DELETE` | `/api/tasks/:id` | ✅ | Hapus task |
| `POST` | `/api/chatbot` | ✅ | Tanya AI tentang task |

Semua error response: `{"error": "pesan error"}` dengan HTTP status code yang sesuai.

---

## Dokumentasi API (Postman)

1. Import `docs/Moonlay_TaskApp.postman_collection.json` ke Postman
2. Import `docs/Moonlay_TaskApp.postman_environment.json` ke Postman
3. Pilih environment **"Moonlay Task API Env"**
4. Jalankan request **"Login"** — token JWT akan otomatis tersimpan ke variable `{{token}}`
5. Semua request lain langsung bisa dipakai

---

## Cara Kerja Chatbot (RAG)

Chatbot menggunakan pendekatan **Retrieval-Augmented Generation (RAG)** sederhana:

### Alur

```
User kirim pertanyaan
       ↓
Backend query PostgreSQL (raw SQL)
  → SELECT semua task + nama assignee
       ↓
Serialize hasil query ke JSON
       ↓
Susun prompt ke Gemini:
  [instruksi sistem]
  [data task sebagai context]
  [pertanyaan user]
       ↓
Gemini 1.5 Flash generate jawaban
  (temperature rendah = 0.3 untuk jawaban faktual)
       ↓
Return jawaban ke user
```

### Kenapa Raw SQL untuk Chatbot?

Query chatbot menggunakan `db.Raw()` (raw SQL) bukan GORM standar, sesuai spesifikasi teknis. Ini memungkinkan:
- JOIN yang fleksibel (task + user dalam satu query)
- Format output yang bisa disesuaikan untuk context LLM
- Tidak bergantung pada struct GORM yang bisa berubah

### Contoh Pertanyaan yang Bisa Dijawab

- "Berapa jumlah task yang sudah selesai?"
- "Tampilkan semua task yang statusnya belum dikerjakan"
- "Tugas apa saja yang deadlinenya paling dekat?"
- "Siapa assignee dari task [judul task]?"
- "Task apa saja yang dikerjakan oleh Budi?"

### Library / Model

- **Library**: `github.com/google/generative-ai-go` (SDK resmi Gemini dari Google)
- **Model**: `gemini-1.5-flash` (cepat, hemat token, cocok untuk Q&A berbasis data)
- **API Key**: Dapatkan gratis di [aistudio.google.com](https://aistudio.google.com)

---

## ERD

Lihat `docs/ERD.md` untuk diagram entitas dan relasi dalam format dbdiagram.io (DBML).

Relasi utama:
- 1 `user` → banyak `task` (sebagai `assignee_id`)
- 1 `user` → banyak `task` (sebagai `created_by`, opsional — audit trail)

---

## Asumsi yang Diambil

1. **Deadline menggunakan format datetime** (bukan hanya tanggal), supaya chatbot bisa menjawab pertanyaan seperti "task yang deadlinenya hari ini" dengan perbandingan waktu yang tepat.

2. **Tidak ada role/permission berbeda antar user** — semua user yang login punya akses yang sama (bisa CRUD semua task). Tidak ada konsep "hanya bisa edit task milik sendiri".

3. **Assignee wajib diisi saat membuat task**, tidak boleh kosong. Field `created_by` bersifat opsional dan diisi otomatis dari JWT token user yang sedang login.

4. **Token disimpan di localStorage** (bukan cookie), karena lebih sederhana untuk SPA dan tidak memerlukan konfigurasi server-side cookie handling.

5. **Seed dilakukan via Go script** (bukan SQL INSERT langsung) supaya bcrypt hash dibuat secara programatik di runtime — tidak perlu pre-compute hash secara manual.

6. **CORS dikonfigurasi** untuk hanya mengizinkan origin `http://localhost:3000` (frontend development server).

7. **Chatbot menggunakan model Gemini 1.5 Flash** — dipilih karena gratis, cepat, dan cukup untuk use case Q&A berbasis data terstruktur.

8. **Tidak ada fitur registrasi user** — semua user di-hardcode via seed script, sesuai spesifikasi.

---

## Commit History (Recommended)

Urutan commit yang disarankan untuk menunjukkan proses berpikir runtut:

```
feat: initial project setup (go.mod, next.js scaffold)
feat: database schema migration (001_init.sql)
feat: user model, repository, seed script
feat: jwt auth middleware and login endpoint
feat: task CRUD endpoints (repository + handlers + routes)
feat: chatbot endpoint with RAG flow (raw SQL + Gemini API)
feat: frontend login page with JWT storage
feat: frontend tasks page with CRUD modals
feat: chatbot widget component (floating UI)
docs: postman collection, ERD, README
```

---

*Dibuat oleh Anisetus Bambang Manalu — Technical Test PT Moonlay Technologies, Juli 2026*
