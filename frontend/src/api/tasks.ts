import { apiRequest } from './client'
import type { Task } from './types'

export type ListTasksParams = {
  status?: string
  priority?: string
  search?: string
}

export async function listTasks(params: ListTasksParams = {}): Promise<Task[]> {
  const res = await apiRequest<{ data: Task[] }>({
    path: '/tasks',
    query: params,
  })
  return res.data
}

export async function createTask(input: { title: string }): Promise<Task> {
  return await apiRequest<Task>({
    method: 'POST',
    path: '/tasks',
    body: input,
  })
}

