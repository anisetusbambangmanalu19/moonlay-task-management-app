# ERD — Moonlay Task Management App

## Diagram (dbdiagram.io)

Buat diagram ini di [dbdiagram.io](https://dbdiagram.io) dengan kode berikut:

```dbml
Table users {
  id bigint [pk, increment, note: "GENERATED ALWAYS AS IDENTITY"]
  name varchar(100) [not null]
  email varchar(150) [unique, not null]
  password_hash varchar(255) [not null, note: "bcrypt hash"]
  created_at timestamptz [not null, default: `now()`]
}

Table tasks {
  id bigint [pk, increment, note: "GENERATED ALWAYS AS IDENTITY"]
  title varchar(200) [not null]
  description text
  status task_status [not null, default: "todo", note: "ENUM: todo | in_progress | done"]
  deadline timestamptz [not null]
  assignee_id bigint [not null, ref: > users.id, note: "ON DELETE RESTRICT"]
  created_by bigint [ref: > users.id, note: "ON DELETE SET NULL — audit trail"]
  created_at timestamptz [not null, default: `now()`]
  updated_at timestamptz [not null, default: `now()`]
}

Ref: tasks.assignee_id > users.id
Ref: tasks.created_by > users.id
```

## Relasi

| Relasi | Tipe | Keterangan |
|---|---|---|
| `users` → `tasks` (assignee) | 1 ke banyak | Satu user bisa menjadi assignee banyak task |
| `users` → `tasks` (created_by) | 1 ke banyak | Audit trail siapa yang membuat task |

## Index

| Index | Tabel | Column | Tujuan |
|---|---|---|---|
| `idx_tasks_assignee` | tasks | assignee_id | Query task by assignee |
| `idx_tasks_status` | tasks | status | Filter by status |
| `idx_tasks_deadline` | tasks | deadline | Sort/filter by deadline |

## Catatan

- `password_hash` di tabel `users` selalu berupa bcrypt hash — tidak pernah plaintext
- `task_status` adalah PostgreSQL ENUM (`todo`, `in_progress`, `done`)
- `created_by` nullable (SET NULL jika user pembuat dihapus) — `assignee_id` NOT nullable (RESTRICT)
- `updated_at` diperbarui otomatis oleh database trigger `trg_tasks_updated_at`
