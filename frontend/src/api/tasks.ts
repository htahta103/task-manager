import { apiRequest } from './client'
import type { Task } from './types'

export type ListTasksParams = {
  status?: string
  priority?: string
  search?: string
}

export async function listTasks(params: ListTasksParams = {}): Promise<Task[]> {
  return await apiRequest<Task[]>({
    path: '/tasks',
    query: params,
  })
}

export async function deleteTask(id: string): Promise<void> {
  await apiRequest({
    method: 'DELETE',
    path: `/tasks/${id}`,
  })
}

