# Moonlay Task Management App

Aplikasi manajemen task sederhana dengan fitur CRUD task, autentikasi JWT, dan AI Chatbot berbasis RAG (Retrieval-Augmented Generation). Proyek ini dikembangkan sebagai technical test untuk posisi **Fullstack Developer Intern @ PT Moonlay Technologies**.

UI didesain secara profesional menyesuaikan branding Moonlay (Cyan/Sky Blue & Dark Mode).

---

## 🛠 Tech Stack

| Layer | Teknologi |
|---|---|
| **Frontend** | Next.js 15 (App Router) + TypeScript + Tailwind CSS v4 |
| **Backend** | Golang 1.23 + Gin Framework |
| **ORM/Query** | GORM (CRUD standar) + Raw SQL via `db.Raw()` (chatbot) |
| **Database** | PostgreSQL (Docker) |
| **Auth** | JWT (access token, expiry 24 jam) + bcrypt |
| **Chatbot** | RAG: PostgreSQL → Custom Prompt → Gemini 1.5 Flash (via REST HTTP) |

---

## 📂 Struktur Folder Utama

```
moonlay-task-management-app/
├── backend/
│   ├── cmd/api/main.go           # Entry point aplikasi backend
│   ├── internal/
│   │   ├── config/               # Koneksi DB dan loader environment
│   │   ├── models/               # Struct GORM untuk User dan Task
│   │   ├── repository/           # Query layer (GORM + Raw SQL)
│   │   ├── handlers/             # HTTP Handlers (Auth, Task, Chatbot)
│   │   ├── middleware/           # Middleware JWT untuk proteksi route
│   │   └── routes/               # Registrasi Endpoint API
│   ├── migrations/
│   │   ├── 001_init.sql          # Skema database (PostgreSQL)
│   │   └── seed/main.go          # Skrip seeder data awal
│   └── .env.example
└── frontend/
    ├── src/
    │   ├── app/                  # Next.js Pages (Login, Tasks)
    │   ├── components/           # UI Components (Card, Form, Chatbot)
    │   ├── lib/                  # Helpers (Axios interceptor, Token)
    │   └── types/                # TypeScript Interfaces
    └── .env.local
```

---

## 🚀 Panduan Setup & Instalasi

### 1. Prasyarat Sistem
Pastikan Anda sudah menginstal:
- **Go** (minimal v1.23)
- **Node.js** (minimal v20+)
- **PostgreSQL** (bisa via Docker atau lokal)

### 2. Setup Database (PostgreSQL)

Jika menggunakan Docker (opsional):
```bash
docker run --name postgres-moonlay -e POSTGRES_PASSWORD=postgrespassword -p 5434:5432 -d postgres
```

Buat database `moonlay_task_db`. Kemudian jalankan script SQL yang ada di `backend/migrations/001_init.sql` melalui tool seperti pgAdmin, DBeaver, atau Navicat.

### 3. Setup Backend (Golang)

1. Masuk ke folder backend dan copy environment:
```bash
cd backend
copy .env.example .env
```

2. Konfigurasikan `.env` sesuai dengan credentials database Anda:
```env
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=postgrespassword
DB_NAME=moonlay_task_db
DB_SSLMODE=disable

# Gunakan Gemini API Key format terbaru dari Google AI Studio (berawalan AQ. atau AIza)
GEMINI_API_KEY=AQ.Ab8RN6KL6Hc4ppmBXNIs...
```

3. Jalankan Seeder untuk mengisi data awal:
```bash
go run migrations/seed/main.go
```

4. Jalankan server Backend:
```bash
go run cmd/api/main.go
```
*Backend berjalan di port `8080`.*

### 4. Setup Frontend (Next.js)

1. Masuk ke folder frontend dan install dependencies:
```bash
cd frontend
npm install
```

2. Jalankan development server:
```bash
npm run dev
```
*Frontend dapat diakses di `http://localhost:3000`.*

---

## 🔑 Kredensial Uji Coba (Login)

Skrip seeder akan secara otomatis membuat 4 user berikut (Password: **nama123**):

| Nama | Email | Password |
|---|---|---|
| Admin | admin@moonlay.com | admin123 |
| Anisetus | anisetus@moonlay.com | anisetus123 |
| Bambang | bambang@moonlay.com | bambang123 |
| Manalu | manalu@moonlay.com | manalu123 |

> *Catatan: Seluruh password disimpan dalam bentuk **bcrypt hash** di database.*

---

## 🤖 Integrasi AI Chatbot (RAG)

Aplikasi ini memiliki fitur Chatbot untuk memonitor tugas. Fitur ini dibuat tanpa external framework yang berat, murni menggunakan pendekatan **Retrieval-Augmented Generation (RAG)** sederhana di backend:

1. **Pengumpulan Konteks (Raw SQL)**: 
   Sesuai dengan instruksi teknis, endpoint chatbot menjalankan kueri **Raw SQL** ke database PostgreSQL. Raw SQL digunakan untuk me-retrieve semua data task beserta nama assignee-nya secara langsung.
2. **Injeksi Prompt**: 
   Data JSON dari database diinjeksikan sebagai *context* ke dalam prompt pengguna.
3. **Panggilan HTTP REST API ke Gemini**: 
   Dikarenakan format API key Google terbaru (`AQ.`) saat ini memiliki isu kompatibilitas dengan Google SDK Go yang lama, panggilan AI dilakukan secara manual menggunakan **net/http** ke URL `https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent`. Key dikirim dengan header `x-goog-api-key`.
4. **Model Gemini 1.5 Flash**: 
   Digunakan karena cepat, hemat resource, dan memiliki free-tier quota yang tinggi (15 RPM).

**Contoh Pertanyaan Chatbot:**
- *"Berapa jumlah task yang sudah selesai?"*
- *"Tampilkan semua task yang statusnya belum dikerjakan"*
- *"Tugas apa saja yang deadlinenya paling dekat?"*
- *"Siapa assignee dari task [judul task]?"*

---

## 📡 API Endpoints

| Method | Path | Auth | Deskripsi |
|---|---|---|---|
| `POST` | `/api/auth/login` | ❌ | Login, mengembalikan JWT Token |
| `GET` | `/api/users` | ✅ | Mengambil daftar semua user (untuk dropdown Assignee) |
| `GET` | `/api/tasks` | ✅ | Mengambil daftar seluruh task |
| `POST` | `/api/tasks` | ✅ | Membuat task baru |
| `PUT` | `/api/tasks/:id` | ✅ | Memperbarui keseluruhan data task |
| `PATCH` | `/api/tasks/:id/status` | ✅ | Memperbarui status task saja |
| `DELETE` | `/api/tasks/:id` | ✅ | Menghapus task berdasarkan ID |
| `POST` | `/api/chatbot` | ✅ | Chat dengan AI (RAG Flow) |

---

## 🚀 Testing API dengan Postman

Anda dapat menguji seluruh endpoint API menggunakan Postman. Semua request, autentikasi otomatis, dan variabel sudah dikonfigurasi.

1. Buka aplikasi **Postman**.
2. Klik tombol **Import** (di kiri atas).
3. Import file berikut dari folder `docs/` di dalam project:
   - `Moonlay_TaskApp.postman_collection.json`
   - `Moonlay_TaskApp.postman_environment.json`
4. Di pojok kanan atas Postman, pastikan Anda mengubah environment dropdown dari *No Environment* menjadi **Moonlay Task API Env**.
5. Buka collection `Moonlay_TaskApp` dan jalankan request **Auth > Login**.
   > *Setelah login sukses, token JWT akan otomatis tersimpan ke variabel environment `{{token}}`.*
6. Anda sekarang bisa menjalankan request lain (seperti Get Tasks, Create Task, Chatbot) tanpa perlu memasukkan token secara manual!

---

## 🎨 UI/UX Design
UI frontend dirancang dengan **Tailwind CSS**, mengadopsi tema **Dark Mode** modern dengan skema warna aksen Cyan/Sky Blue (`#0095c8`) yang selaras dengan identitas **PT Moonlay Technologies**. Terdapat efek *glassmorphism*, gradient atraktif, dan animasi fluid untuk transisi.

---

*Dibuat oleh **Anisetus Bambang Manalu** — Technical Test PT Moonlay Technologies, Juli 2026*
