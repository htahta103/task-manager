import { create } from 'zustand'
import type { Task } from '../api/types'
import { deleteTask, listTasks } from '../api/tasks'

type TasksState = {
  tasks: Task[]
  isLoading: boolean
  error: Error | null
  refresh: () => Promise<void>
  removeTask: (id: string) => Promise<void>
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
  removeTask: async (id) => {
    set({ error: null })
    try {
      await deleteTask(id)
      set((state) => ({
        tasks: state.tasks.filter((task) => task.id !== id),
      }))
    } catch (e) {
      const err =
        e instanceof Error ? e : new Error('Failed to delete task')
      set({ error: err })
      throw err
    }
  },
}))

