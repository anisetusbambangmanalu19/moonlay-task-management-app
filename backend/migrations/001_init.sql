-- ============================================================
-- Migration 001: Initial Schema
-- Moonlay Task Management App
-- Run this file in psql or Navicat Query Tool on moonlay_task_db
-- ============================================================

-- Create custom ENUM type for task status
CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');

-- ============================================================
-- Table: users
-- Stores all application users (seeded manually, no registration)
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name          VARCHAR(100)  NOT NULL,
    email         VARCHAR(150)  UNIQUE NOT NULL,
    password_hash VARCHAR(255)  NOT NULL,  -- bcrypt hash, NEVER plaintext
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- ============================================================
-- Table: tasks
-- Core task management table with status, deadline, and assignee
-- ============================================================
CREATE TABLE IF NOT EXISTS tasks (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title       VARCHAR(200)  NOT NULL,
    description TEXT,
    status      task_status   NOT NULL DEFAULT 'todo',
    deadline    TIMESTAMPTZ   NOT NULL,
    assignee_id BIGINT        NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_by  BIGINT        REFERENCES users(id) ON DELETE SET NULL,  -- audit trail
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- ============================================================
-- Indexes: optimize common query patterns
-- ============================================================
CREATE INDEX IF NOT EXISTS idx_tasks_assignee ON tasks(assignee_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status   ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_deadline ON tasks(deadline);

-- ============================================================
-- Trigger: auto-update updated_at on every task UPDATE
-- ============================================================
CREATE OR REPLACE FUNCTION fn_update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_tasks_updated_at ON tasks;
CREATE TRIGGER trg_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION fn_update_updated_at();
