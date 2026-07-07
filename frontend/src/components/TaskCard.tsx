'use client';

import { Task, TaskStatus } from '@/types';
import StatusBadge from './StatusBadge';
import { format } from 'date-fns';
import { id } from 'date-fns/locale';

interface TaskCardProps {
  task: Task;
  onEdit: (task: Task) => void;
  onDelete: (task: Task) => void;
  onStatusChange: (taskId: number, status: TaskStatus) => void;
}

export default function TaskCard({ task, onEdit, onDelete, onStatusChange }: TaskCardProps) {
  const deadline = new Date(task.deadline);
  const isOverdue = deadline < new Date() && task.status !== 'done';
  const isToday =
    deadline.toDateString() === new Date().toDateString();

  return (
    <div className="group relative bg-slate-800/60 backdrop-blur-sm border border-slate-700/50 rounded-2xl p-5 hover:border-sky-500/40 hover:bg-slate-800/80 transition-all duration-300 hover:shadow-lg hover:shadow-sky-500/5">
      {/* Bilah aksen gradien */}
      <div className="absolute top-0 left-0 right-0 h-0.5 rounded-t-2xl bg-gradient-to-r from-sky-500/0 via-sky-500/50 to-cyan-500/0 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />

      {/* Header */}
      <div className="flex items-start justify-between gap-3 mb-3">
        <h3 className="font-semibold text-white text-sm leading-snug line-clamp-2 flex-1">
          {task.title}
        </h3>
        <StatusBadge status={task.status} size="sm" />
      </div>

      {/* Deskripsi */}
      {task.description && (
        <p className="text-slate-400 text-xs leading-relaxed line-clamp-2 mb-4">
          {task.description}
        </p>
      )}

      {/* Tenggat waktu */}
      <div className="flex items-center gap-1.5 mb-3">
        <svg className="w-3.5 h-3.5 text-slate-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span
          className={`text-xs font-medium ${
            isOverdue
              ? 'text-red-400'
              : isToday
              ? 'text-amber-400'
              : 'text-slate-400'
          }`}
        >
          {isOverdue && '⚠️ '}
          {format(deadline, 'dd MMM yyyy, HH:mm', { locale: id })}
        </span>
      </div>

      {/* Assignee */}
      <div className="flex items-center gap-2 mb-4">
        <div className="h-6 w-6 rounded-full bg-gradient-to-br from-sky-500 to-cyan-500 flex items-center justify-center text-white text-xs font-bold flex-shrink-0">
          {task.assignee?.name?.charAt(0).toUpperCase() ?? '?'}
        </div>
        <span className="text-xs text-slate-400">{task.assignee?.name ?? 'Unknown'}</span>
      </div>

      {/* Perubahan status cepat */}
      <div className="mb-4">
        <select
          value={task.status}
          onChange={(e) => onStatusChange(task.id, e.target.value as TaskStatus)}
          className="w-full bg-slate-900/60 border border-slate-700/50 rounded-lg px-3 py-1.5 text-xs text-slate-300 focus:outline-none focus:border-sky-500/50 focus:ring-1 focus:ring-sky-500/30 transition-all cursor-pointer"
          onClick={(e) => e.stopPropagation()}
        >
          <option value="todo">📋 Todo</option>
          <option value="in_progress">⚡ In Progress</option>
          <option value="done">✅ Done</option>
        </select>
      </div>

      {/* Tombol aksi */}
      <div className="flex gap-2">
        <button
          onClick={() => onEdit(task)}
          className="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-xl bg-sky-500/10 hover:bg-sky-500/20 text-sky-400 text-xs font-medium transition-all duration-200 border border-sky-500/20 hover:border-sky-500/40"
        >
          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
          Edit
        </button>
        <button
          onClick={() => onDelete(task)}
          className="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 rounded-xl bg-red-500/10 hover:bg-red-500/20 text-red-400 text-xs font-medium transition-all duration-200 border border-red-500/20 hover:border-red-500/40"
        >
          <svg className="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          Hapus
        </button>
      </div>
    </div>
  );
}
