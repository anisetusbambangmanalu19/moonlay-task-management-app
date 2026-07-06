'use client';

import { TaskStatus } from '@/types';

interface StatusBadgeProps {
  status: TaskStatus;
  size?: 'sm' | 'md';
}

const statusConfig: Record<TaskStatus, { label: string; className: string }> = {
  todo: {
    label: 'Todo',
    className: 'bg-amber-500/20 text-amber-400 border border-amber-500/30',
  },
  in_progress: {
    label: 'In Progress',
    className: 'bg-blue-500/20 text-blue-400 border border-blue-500/30',
  },
  done: {
    label: 'Done',
    className: 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30',
  },
};

export default function StatusBadge({ status, size = 'md' }: StatusBadgeProps) {
  const config = statusConfig[status];
  const sizeClass = size === 'sm' ? 'px-2 py-0.5 text-xs' : 'px-3 py-1 text-xs';

  return (
    <span className={`inline-flex items-center gap-1.5 rounded-full font-medium ${sizeClass} ${config.className}`}>
      <span
        className={`h-1.5 w-1.5 rounded-full ${
          status === 'todo'
            ? 'bg-amber-400'
            : status === 'in_progress'
            ? 'bg-blue-400 animate-pulse'
            : 'bg-emerald-400'
        }`}
      />
      {config.label}
    </span>
  );
}
