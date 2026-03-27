import { create } from 'zustand'
import type { Task } from '../api/types'
import { listTasks } from '../api/tasks'

type TasksState = {
  tasks: Task[]
  isLoading: boolean
  error: Error | null
  refresh: () => Promise<void>
}

export const useTasksStore = create<TasksState>((set) => ({
  tasks: [],
  isLoading: false,
  error: null,
  refresh: async () => {
    set({ isLoading: true, error: null })
    try {
      const tasks = await listTasks()
      set({ tasks, isLoading: false })
    } catch (e) {
      const err = e instanceof Error ? e : new Error('Failed to load tasks')
      set({ error: err, isLoading: false })
    }
  },
}))

