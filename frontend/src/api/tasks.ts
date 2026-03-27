import { apiRequest } from './client'
import type { Task } from './types'

export type ListTasksParams = {
  status?: string
  priority?: string
  search?: string
}

type ListTasksResponse = {
  data: Task[]
  count: number
}

export async function listTasks(params: ListTasksParams = {}): Promise<Task[]> {
  const res = await apiRequest<Task[] | ListTasksResponse>({
    path: '/api/tasks',
    query: params,
  })
  if (Array.isArray(res)) return res
  return Array.isArray(res.data) ? res.data : []
}

