'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { Task, TaskStatus, User } from '@/types';
import api from '@/lib/api';
import { auth } from '@/lib/auth';
import TaskCard from '@/components/TaskCard';
import TaskForm from '@/components/TaskForm';
import ChatbotWidget from '@/components/ChatbotWidget';

type FilterStatus = 'all' | TaskStatus;

export default function TasksPage() {
  const router = useRouter();

  const [tasks, setTasks] = useState<Task[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editTask, setEditTask] = useState<Task | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<Task | null>(null);
  const [filterStatus, setFilterStatus] = useState<FilterStatus>('all');
  const [searchQuery, setSearchQuery] = useState('');
  const [deleting, setDeleting] = useState(false);
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null);

  useEffect(() => {
    if (!auth.isAuthenticated()) {
      router.push('/login');
      return;
    }
    const user = auth.getUser();
    setCurrentUser(user as User);
  }, [router]);

  const fetchTasks = useCallback(async () => {
    try {
      const res = await api.get('/tasks');
      setTasks(res.data.data ?? []);
    } catch {
      showToast('Gagal memuat task', 'error');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (auth.isAuthenticated()) fetchTasks();
  }, [fetchTasks]);

  const showToast = (message: string, type: 'success' | 'error') => {
    setToast({ message, type });
    setTimeout(() => setToast(null), 3000);
  };

  const handleStatusChange = async (taskId: number, status: TaskStatus) => {
    try {
      await api.patch(`/tasks/${taskId}/status`, { status });
      setTasks((prev) =>
        prev.map((t) => (t.id === taskId ? { ...t, status } : t))
      );
      showToast('Status berhasil diperbarui', 'success');
    } catch {
      showToast('Gagal mengubah status', 'error');
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    setDeleting(true);
    try {
      await api.delete(`/tasks/${deleteTarget.id}`);
      setTasks((prev) => prev.filter((t) => t.id !== deleteTarget.id));
      showToast('Task berhasil dihapus', 'success');
      setDeleteTarget(null);
    } catch {
      showToast('Gagal menghapus task', 'error');
    } finally {
      setDeleting(false);
    }
  };

  const handleLogout = () => {
    auth.logout();
    router.push('/login');
  };

  const filteredTasks = tasks
    .filter((t) => filterStatus === 'all' || t.status === filterStatus)
    .filter((t) =>
      searchQuery === '' ||
      t.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      t.assignee?.name?.toLowerCase().includes(searchQuery.toLowerCase())
    );

  const stats = {
    total: tasks.length,
    todo: tasks.filter((t) => t.status === 'todo').length,
    inProgress: tasks.filter((t) => t.status === 'in_progress').length,
    done: tasks.filter((t) => t.status === 'done').length,
  };

  return (
    <div className="min-h-screen gradient-bg">
      {/* Notifikasi toast */}
      {toast && (
        <div
          className={`fixed top-4 right-4 z-50 px-4 py-3 rounded-2xl text-sm font-medium shadow-lg animate-in transition-all ${toast.type === 'success'
              ? 'bg-emerald-500/20 border border-emerald-500/30 text-emerald-300'
              : 'bg-red-500/20 border border-red-500/30 text-red-300'
            }`}
        >
          {toast.message}
        </div>
      )}

      {/* Header */}
      <header className="sticky top-0 z-30 glass-card-light backdrop-blur-lg border-b border-slate-700/40">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="h-10 w-10 rounded-2xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center shadow-lg shadow-blue-500/20">
              <svg className="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
              </svg>
            </div>
            <div>
              <h1 className="text-lg font-bold text-white">Moonlay Tasks</h1>
              <p className="text-xs text-slate-500">PT Moonlay Technologies</p>
            </div>
          </div>

          <div className="flex items-center gap-3">
            {currentUser && (
              <div className="hidden sm:flex items-center gap-2.5 px-3.5 py-2 rounded-lg bg-slate-800/40 border border-slate-700/40">
                <div className="h-7 w-7 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center text-white text-xs font-bold">
                  {currentUser.name?.charAt(0).toUpperCase()}
                </div>
                <span className="text-sm text-slate-300 font-medium">{currentUser.name}</span>
              </div>
            )}
            <button
              onClick={handleLogout}
              className="flex items-center gap-1.5 px-4 py-2 rounded-lg bg-red-500/15 hover:bg-red-500/25 border border-red-500/30 text-red-400 hover:text-red-300 text-sm font-semibold transition-all"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              Keluar
            </button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Kartu statistik */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
          {[
            { label: 'Total Task', value: stats.total, color: 'text-slate-300', icon: '📋', bg: 'from-slate-600/20 to-slate-700/20' },
            { label: 'Todo', value: stats.todo, color: 'text-amber-400', icon: '⏳', bg: 'from-amber-500/15 to-amber-600/10' },
            { label: 'In Progress', value: stats.inProgress, color: 'text-blue-400', icon: '⚡', bg: 'from-blue-500/15 to-blue-600/10' },
            { label: 'Done', value: stats.done, color: 'text-emerald-400', icon: '✅', bg: 'from-emerald-500/15 to-emerald-600/10' },
          ].map((stat) => (
            <div
              key={stat.label}
              className={`glass-card-light bg-gradient-to-br ${stat.bg} rounded-2xl p-4 border-slate-700/40 hover:border-slate-700/60 transition-all`}
            >
              <div className="flex items-start justify-between">
                <div>
                  <p className="text-xs text-slate-500 font-semibold mb-2">{stat.label}</p>
                  <p className={`text-3xl font-bold ${stat.color}`}>{stat.value}</p>
                </div>
                <span className="text-2xl opacity-60">{stat.icon}</span>
              </div>
            </div>
          ))}
        </div>

        {/* Bilah alat */}
        <div className="flex flex-col sm:flex-row gap-3 mb-8">
          {/* Pencarian */}
          <div className="relative flex-1">
            <svg className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input
              type="text"
              placeholder="Cari task atau assignee..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-slate-800/40 border border-slate-700/40 rounded-xl pl-10 pr-4 py-2.5 text-sm text-white placeholder-slate-600 focus:outline-none focus:border-blue-500/60 focus:bg-slate-800/60 focus:ring-1 focus:ring-blue-500/20 transition-all"
            />
          </div>

          {/* Tombol filter */}
          <div className="flex gap-2 flex-wrap">
            {(['all', 'todo', 'in_progress', 'done'] as FilterStatus[]).map((status) => (
              <button
                key={status}
                onClick={() => setFilterStatus(status)}
                className={`px-4 py-2 rounded-lg text-xs font-semibold transition-all whitespace-nowrap ${filterStatus === status
                    ? 'bg-gradient-to-r from-blue-500 to-cyan-500 text-white shadow-lg shadow-blue-500/30'
                    : 'bg-slate-800/40 border border-slate-700/40 text-slate-400 hover:text-white hover:bg-slate-800/60 hover:border-slate-600/50'
                  }`}
              >
                {status === 'all' ? 'Semua' : status === 'todo' ? 'Todo' : status === 'in_progress' ? 'In Progress' : 'Done'}
              </button>
            ))}
          </div>

          {/* Tombol tambah */}
          <button
            onClick={() => { setEditTask(null); setShowForm(true); }}
            className="flex items-center justify-center sm:justify-start gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-cyan-500 hover:from-blue-600 hover:to-cyan-600 text-white text-sm font-semibold transition-all shadow-lg shadow-blue-500/25 hover:shadow-blue-500/40 active:scale-95"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 4v16m8-8H4" />
            </svg>
            <span className="hidden sm:inline">Tambah Task</span>
            <span className="sm:hidden">Tambah</span>
          </button>
        </div>

        {/* Grid task */}
        {loading ? (
          <div className="flex items-center justify-center py-24">
            <div className="flex flex-col items-center gap-4 animate-pulse">
              <svg className="w-12 h-12 text-blue-500/60" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <p className="text-slate-400 text-sm font-medium">Memuat task...</p>
            </div>
          </div>
        ) : filteredTasks.length === 0 ? (
          <div className="text-center py-24">
            <div className="inline-flex items-center justify-center h-16 w-16 rounded-2xl bg-slate-800/40 border border-slate-700/40 mb-4">
              <svg className="w-8 h-8 text-slate-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <p className="text-slate-300 font-semibold mb-2">
              {searchQuery || filterStatus !== 'all' ? 'Tidak ada task' : 'Belum ada task'}
            </p>
            <p className="text-slate-600 text-sm mb-6">
              {searchQuery || filterStatus !== 'all' ? 'Ubah filter atau kata kunci pencarian' : 'Buat task pertama Anda sekarang'}
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {filteredTasks.map((task) => (
              <TaskCard
                key={task.id}
                task={task}
                onEdit={(t) => { setEditTask(t); setShowForm(true); }}
                onDelete={(t) => setDeleteTarget(t)}
                onStatusChange={handleStatusChange}
              />
            ))}
          </div>
        )}
      </main>

      {/* Modal formulir task */}
      {showForm && (
        <TaskForm
          task={editTask}
          onClose={() => { setShowForm(false); setEditTask(null); }}
          onSuccess={() => {
            fetchTasks();
            showToast(editTask ? 'Task berhasil diperbarui' : 'Task berhasil dibuat', 'success');
          }}
        />
      )}

      {/* Konfirmasi hapus */}
      {deleteTarget && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div className="absolute inset-0 bg-black/70 backdrop-blur-sm" onClick={() => setDeleteTarget(null)} />
          <div className="relative w-full max-w-sm glass-card rounded-2xl p-6 shadow-2xl animate-in">
            <div className="h-0.5 -mx-6 -mt-6 mb-6 bg-gradient-to-r from-transparent via-red-500/50 to-transparent rounded-t-2xl" />

            <div className="flex items-center gap-3 mb-5">
              <div className="h-10 w-10 rounded-lg bg-red-500/20 border border-red-500/40 flex items-center justify-center flex-shrink-0">
                <svg className="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
              </div>
              <div>
                <h3 className="font-bold text-white">Hapus Task?</h3>
                <p className="text-xs text-slate-400">Tindakan ini tidak dapat dibatalkan</p>
              </div>
            </div>

            <div className="bg-slate-800/60 rounded-lg p-3 mb-5 border border-slate-700/40">
              <p className="text-sm text-slate-300 font-medium">&quot;{deleteTarget.title}&quot;</p>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => setDeleteTarget(null)}
                className="flex-1 py-2.5 rounded-lg border border-slate-700/40 text-slate-300 hover:text-white hover:bg-slate-800/60 text-sm font-semibold transition-all"
              >
                Batal
              </button>
              <button
                onClick={handleDelete}
                disabled={deleting}
                className="flex-1 py-2.5 rounded-lg bg-red-500 hover:bg-red-600 text-white text-sm font-semibold transition-all disabled:opacity-50 shadow-lg shadow-red-500/20"
              >
                {deleting ? 'Menghapus...' : 'Hapus'}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Chat mengambang */}
      <ChatbotWidget />
    </div>
  );
}