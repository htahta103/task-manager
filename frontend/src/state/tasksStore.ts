import { create } from 'zustand'
import type { Task } from '../api/types'
import { createTask, listTasks } from '../api/tasks'

type TasksState = {
  tasks: Task[]
  isLoading: boolean
  error: Error | null
  refresh: () => Promise<void>
  add: (title: string) => Promise<void>
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
  add: async (title: string) => {
    set({ error: null })
    try {
      const created = await createTask({ title })
      set((s) => ({ tasks: [created, ...s.tasks] }))
    } catch (e) {
      const err = e instanceof Error ? e : new Error('Failed to add task')
      set({ error: err })
    }
  },
}))

