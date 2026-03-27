export type TaskStatus = 'pending' | 'in_progress' | 'done'
export type TaskPriority = 'low' | 'medium' | 'high'

export type Task = {
  id: string
  title: string
  description?: string | null
  status: TaskStatus
  priority: TaskPriority
  due_date?: string | null
  created_at: string
  updated_at: string
}

export type ApiError = {
  message: string
  code?: string
  details?: unknown
}

