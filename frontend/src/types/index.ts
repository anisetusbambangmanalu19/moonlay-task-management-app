export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}

export type TaskStatus = 'todo' | 'in_progress' | 'done';

export interface Task {
  id: number;
  title: string;
  description: string;
  status: TaskStatus;
  deadline: string;
  assignee_id: number;
  created_by: number | null;
  created_at: string;
  updated_at: string;
  assignee?: User;
  creator?: User;
}

export interface LoginResponse {
  token: string;
  user: Pick<User, 'id' | 'name' | 'email'>;
}

export interface ApiResponse<T> {
  data?: T;
  message?: string;
  error?: string;
}

export interface ChatMessage {
  role: 'user' | 'bot';
  content: string;
  timestamp: Date;
}
